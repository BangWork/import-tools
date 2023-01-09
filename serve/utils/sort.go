package utils

import (
	"errors"
	"math"
	"sort"
)

const (
	sortPositionLimit = math.MaxInt64 - 1
)

type KVInt struct {
	Key   string
	Value int64
}

func NewKVInt(key string, value int64) *KVInt {
	kv := new(KVInt)
	kv.Key = key
	kv.Value = value
	return kv
}

// 将 key 加入到切片首部, value > fromValue 并且安value顺序递增, 返回需要变化的kv值
func ShuffleKV(key string, fromValue int64, kvs []*KVInt) (changeKVs []*KVInt, err error) {
	changeKVs = make([]*KVInt, 0)
	var steps int64 = 2
	currentKey := key
	currentValue := fromValue
	found := false
	for _, kv := range kvs {
		// 值太大时候重新洗牌
		if currentValue > sortPositionLimit {
			err = errors.New("out of max int64")
			return
		}
		if kv.Value-currentValue >= steps {
			changeKVs = append(changeKVs, NewKVInt(currentKey, currentValue+(kv.Value-currentValue)/2))
			found = true
		} else {
			changeKVs = append(changeKVs, NewKVInt(currentKey, currentValue+steps))
		}
		currentKey = kv.Key
		currentValue = currentValue + steps
		steps = steps * 2
		if found {
			break
		}
	}
	if !found {
		changeKVs = append(changeKVs, NewKVInt(currentKey, currentValue+steps))
	}
	return
}

func MarkPosition(elements []string) []*KVInt {
	kvs := make([]*KVInt, len(elements))
	var currentPosition int64 = 1000000
	for pos, element := range elements {
		kvs[pos] = NewKVInt(element, currentPosition)
		currentPosition *= 2
	}
	return kvs
}

func isdigit(i byte) bool {
	return '0' <= i && i <= '9'
}

func islower(i byte) bool {
	return 'a' <= i && i <= 'z'
}

func isupper(i byte) bool {
	return 'A' <= i && i <= 'Z'
}

func isalnum(i byte) bool {
	return isdigit(i) || islower(i) || isupper(i)
}

func compareAlnum(first byte, second byte) bool {
	b1 := first
	b2 := second
	if isupper(b1) {
		b1 = b1 + 'a' - 'A'
	}
	if isupper(b2) {
		b2 = b2 + 'a' - 'A'
	}
	return b1 < b2
}

func isEqual(first byte, second byte) bool {
	if isupper(first) {
		first = first + 'a' - 'A'
	}
	if isupper(second) {
		second = second + 'a' - 'A'
	}
	return first == second
}

func compareByte(first byte, second byte) bool {
	if isalnum(first) && isalnum(second) {
		return compareAlnum(first, second)
	} else if isalnum(first) {
		return true
	} else if isalnum(second) {
		return false
	} else {
		return first < second
	}
}

func CompareString(first string, second string) bool {
	len1 := len(first)
	len2 := len(second)
	i := 0
	for {
		if i == len1 || i == len2 {
			break
		}
		if !isEqual(first[i], second[i]) {
			return compareByte(first[i], second[i])
		}
		i++
	}
	if len1 == len2 {
		return false
	}
	return i == len1
}

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Int64Slice) Sort()              { sort.Sort(p) }
