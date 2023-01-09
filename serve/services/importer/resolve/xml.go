package resolve

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/bangwork/import-tools/serve/services/cache"

	"golang.org/x/net/html"
)

type stack struct {
	data []interface{}
}

func (s *stack) empty() bool {
	return len(s.data) == 0
}

func (s *stack) push(value interface{}) {
	s.data = append(s.data, value)
}

func (s *stack) pop() interface{} {
	value := s.data[len(s.data)-1]
	s.data[len(s.data)-1] = nil
	s.data = s.data[:len(s.data)-1]
	return value
}

func (s *stack) peek() interface{} {
	if len(s.data) == 0 {
		return nil
	}
	return s.data[len(s.data)-1]
}

type Element struct {
	Tag     string
	AttrMap map[string]string
	Text    string
	Child   []*Element
	parent  *Element
}

func (e *Element) addChild(t *Element) {
	e.Child = append(e.Child, t)
}

func (e *Element) setText(text string) {
	e.Text = text
}

func newElement(tag string, parent *Element) *Element {
	e := &Element{
		Tag:     tag,
		AttrMap: make(map[string]string),
		Child:   make([]*Element, 0),
		parent:  parent,
	}
	if parent != nil {
		parent.addChild(e)
	}
	return e
}

func isWhitespace(s string) bool {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			return false
		}
	}
	return true
}

type XmlScanner struct {
	decoder *xml.Decoder
	rootTag string
	cache.ResolveResult
	HoursPerDay string
	DaysPerWeek string
}

func NewXmlScanner(reader io.Reader, rootTag string) *XmlScanner {
	return &XmlScanner{decoder: xml.NewDecoder(reader), rootTag: rootTag}
}

func (o *XmlScanner) NextElement() *Element {
	var stack stack
	for {
		t, err := o.decoder.RawToken()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Println(err)
			return nil
		}

		var top *Element
		switch e := stack.peek().(type) {
		case *Element:
			top = e
		case nil:
			top = nil
		}

		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == o.rootTag {
				continue
			}
			e := newElement(t.Name.Local, top)
			for _, a := range t.Attr {
				e.AttrMap[a.Name.Local] = a.Value
			}
			stack.push(e)
		case xml.EndElement:
			if t.Name.Local == o.rootTag {
				continue
			}
			e := stack.pop().(*Element)
			if e.parent == nil {
				return e
			}
		case xml.CharData:
			s := string(t)
			if top != nil && !isWhitespace(s) {
				text := fmt.Sprintf("<![CDATA[%s]]>", s)
				top.setText(text)
			}
		case xml.Comment:
			o.parseComment(string(t))
		}
	}
}

func (o *XmlScanner) NextElements(number int) []*Element {
	var stack stack
	var elements = make([]*Element, 0)
	var flag = true
	for flag {
		t, err := o.decoder.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}

		var top *Element
		switch e := stack.peek().(type) {
		case *Element:
			top = e
		case nil:
			top = nil
		}

		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == o.rootTag {
				continue
			}
			e := newElement(t.Name.Local, top)
			for _, a := range t.Attr {
				e.AttrMap[a.Name.Local] = a.Value
			}
			stack.push(e)
		case xml.EndElement:
			if t.Name.Local == o.rootTag {
				continue
			}

			e := stack.pop().(*Element)

			if e.parent == nil {
				elements = append(elements, e)
			}

			if len(elements) == number {
				flag = false
			}
		case xml.CharData:
			s := string(t)
			if top != nil && !isWhitespace(s) {
				text := fmt.Sprintf("<![CDATA[%s]]>", s)
				top.setText(text)
			}
		case xml.Comment:
			o.parseComment(string(t))
		}
	}
	return elements
}

func (o *XmlScanner) parseComment(comment string) {
	searchMap := map[string]bool{
		"_User_":                            true,
		"_FileAttachment_":                  true,
		"_Project_":                         true,
		"_Issue_":                           true,
		"_JIRA Build_":                      true,
		"_Server ID_":                       true,
		"_jira.timetracking.hours.per.day_": true,
		"_jira.timetracking.days.per.week_": true,
	}
	lines := strings.Split(comment, "\n")
	for _, s := range lines {
		strSlice := strings.Split(s, " : ")
		if len(strSlice) != 2 {
			continue
		}
		name := fmt.Sprintf("_%s_", strings.TrimSpace(strSlice[0]))
		value := strings.TrimSpace(strSlice[1])
		if !searchMap[name] {
			continue
		}
		switch name {
		case "_User_":
			o.ResolveResult.MemberCount, _ = strconv.ParseInt(value, 10, 64)
		case "_FileAttachment_":
			o.ResolveResult.AttachmentCount, _ = strconv.ParseInt(value, 10, 64)
		case "_Issue_":
			o.ResolveResult.IssueCount, _ = strconv.ParseInt(value, 10, 64)
		case "_Project_":
			o.ResolveResult.ProjectCount, _ = strconv.ParseInt(value, 10, 64)
		case "_JIRA Build_":
			strSlice = strings.Split(value, "#")
			if len(strSlice) != 2 {
				continue
			}
			o.ResolveResult.JiraVersion = strSlice[0]
		case "_Server ID_":
			o.ResolveResult.JiraServerID = value
		case "_jira.timetracking.hours.per.day_":
			o.HoursPerDay = value
		case "_jira.timetracking.days.per.week_":
			o.DaysPerWeek = value
		}
	}
	if o.ResolveResult.JiraVersion == "" {
		o.ResolveResult.JiraVersion = "Cloud"
	}
	return
}

func (o *Element) Encode() string {
	if o == nil {
		return ""
	}
	str := "<" + o.Tag
	for name, value := range o.AttrMap {
		str += fmt.Sprintf(` %s="%s"`, name, html.EscapeString(value))
	}
	str += ">"
	if len(o.Child) == 0 {
		str += fmt.Sprintf("%s</%s>", o.Text, o.Tag)
		return str
	}
	for _, child := range o.Child {
		str += child.Encode()
	}
	str += fmt.Sprintf("%s</%s>", o.Text, o.Tag)
	return str
}
