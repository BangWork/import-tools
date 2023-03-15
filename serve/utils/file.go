package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	. "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	_ "golang.org/x/image/webp"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/rwcarlsen/goexif/exif"
)

type ExifValType struct {
	Val  string `json:"val"`
	Type int    `json:"type"`
}

type FileInfo struct {
	Name        string
	Size        int64
	Mime        string
	Hash        string
	ImageWidth  int
	ImageHeight int
	Exif        string
}

const (
	DetectSize int64 = 512
	BlockBits        = 22 // Indicate that the blocksize is 4M
	BlockSize        = 1 << BlockBits
)

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcPath := path.Join(src, fd.Name())
		dstPath := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func CopyFile(src, dst string) error {
	var err error
	var srcFile *os.File
	var dstFile *os.File
	var srcInfo os.FileInfo

	if srcFile, err = os.Open(src); err != nil {
		return err
	}
	defer srcFile.Close()

	if dstFile, err = os.Create(dst); err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

func GetDirSize(path string) (int64, error) {
	if !CheckPathExist(path) {
		return 0, nil
	}
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return 0, err
	}
	return size, nil
}

func GetFileSize(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func GetFileModTime(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	fi, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fi.ModTime().Unix(), nil
}

func CheckPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ListZipFile(path string) ([]string, error) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, common.Errors(common.NotFoundError, nil)
	}
	sortFiles := make(ByModTime, 0)
	for _, v := range fileInfoList {
		if strings.HasSuffix(v.Name(), ".zip") {
			sortFiles = append(sortFiles, v)
		}
	}
	sort.Sort(sortFiles)

	files := make([]string, 0, len(sortFiles))
	for _, v := range sortFiles {
		files = append(files, v.Name())
	}
	return files, nil
}

type ByModTime []os.FileInfo

func (fis ByModTime) Len() int {
	return len(fis)
}

func (fis ByModTime) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis ByModTime) Less(i, j int) bool {
	return fis[j].ModTime().Before(fis[i].ModTime())
}

type ByModTimeAsc []os.FileInfo

func (fis ByModTimeAsc) Len() int {
	return len(fis)
}

func (fis ByModTimeAsc) Swap(i, j int) {
	fis[i], fis[j] = fis[j], fis[i]
}

func (fis ByModTimeAsc) Less(i, j int) bool {
	return fis[j].ModTime().After(fis[i].ModTime())
}

func GetFileInfo(file *os.File) (fileInfo *FileInfo, err error) {
	seekZero := func(r *os.File) error {
		if _, err = r.Seek(0, 0); err != nil {
			return err
		}
		return nil
	}

	ext := filepath.Ext(file.Name())
	fMime, err := DetectFromReaderAndExt(file, ext)
	name := filepath.Base(file.Name())
	if err = seekZero(file); err != nil {
		return
	}
	hash, _ := QEtag(file)
	if err = seekZero(file); err != nil {
		return
	}
	fi, _ := file.Stat()
	size := fi.Size()

	var imageWidth, imageHeight int
	var imageExif string
	if strings.HasPrefix(fMime, "image/") {
		if err = seekZero(file); err != nil {
			return
		}
		imageWidth, imageHeight, _ = DecodeDimension(file)
		if err = seekZero(file); err != nil {
			return
		}
		imageExif, _ = DecodeExifToJSONString(file)
	}
	return &FileInfo{
		Name:        name,
		Size:        size,
		Mime:        fMime,
		Hash:        hash,
		ImageWidth:  imageWidth,
		ImageHeight: imageHeight,
		Exif:        imageExif,
	}, nil
}

func DetectFromBytes(b []byte) string {
	m := http.DetectContentType(b)
	i := strings.IndexByte(m, ';')
	if i > 0 {
		m = m[:i]
	}
	return m
}

func DetectFromReaderAndExt(r io.Reader, ext string) (string, error) {
	b := make([]byte, DetectSize)
	n, err := r.Read(b)
	if err != nil && err != io.EOF {
		return "", err
	}
	return DetectFromBytesAndExt(b[:n], ext), nil
}

func DetectFromBytesAndExt(b []byte, ext string) string {
	m := DetectFromBytes(b)
	if m == "application/octet-stream" ||
		m == "application/zip" ||
		m == "text/plain" {

		extm := DetectFromExt(ext)
		if extm != "" {
			m = extm
		}
	}
	return m
}

func DetectFromExt(ext string) string {
	m := mime.TypeByExtension(ext)
	i := strings.IndexByte(m, ';')
	if i > 0 {
		m = m[:i]
	}
	return m
}

func DecodeDimension(r io.Reader) (width int, height int, err error) {
	config, _, err := DecodeConfig(r)
	return config.Width, config.Height, err
}

func DecodeExifToJSONString(r io.Reader) (string, error) {
	x, err := exif.Decode(r)
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return "", err
	}
	t, _ := x.Get(exif.Orientation)
	if t == nil {
		return "", nil
	}
	m := map[exif.FieldName]*ExifValType{
		exif.Orientation: &ExifValType{
			Val:  orientationLabel(t.String()),
			Type: int(t.Type),
		},
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

func orientationLabel(val string) string {
	switch val {
	case "1":
		return "Top-left"
	case "2":
		return "Top-right"
	case "3":
		return "Bottom-right"
	case "4":
		return "Bottom-left"
	case "5":
		return "Left-top"
	case "6":
		return "Right-top"
	case "7":
		return "Right-bottom"
	case "8":
		return "Left-bottom"
	default:
		return "Top-left"
	}
}

func QEtag(f *os.File) (etag string, err error) {
	fi, err := f.Stat()
	if err != nil {
		return
	}

	size := fi.Size()
	blockCnt := blockCount(size)
	sha1Buf := make([]byte, 0, 21)

	if blockCnt <= 1 { // file size <= 4M
		sha1Buf = append(sha1Buf, 0x16)
		sha1Buf, err = calSha1(sha1Buf, f)
		if err != nil {
			return
		}
	} else { // file size > 4M
		sha1Buf = append(sha1Buf, 0x96)
		sha1BlockBuf := make([]byte, 0, blockCnt*20)
		for i := 0; i < blockCnt; i++ {
			body := io.LimitReader(f, BlockSize)
			sha1BlockBuf, err = calSha1(sha1BlockBuf, body)
			if err != nil {
				return
			}
		}
		sha1Buf, _ = calSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
	}
	etag = base64.URLEncoding.EncodeToString(sha1Buf)
	return
}

func blockCount(fsize int64) int {
	return int((fsize + (BlockSize - 1)) >> BlockBits)
}

func calSha1(b []byte, r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(b), nil
}
