package common

var MaxScanTokenSize int

func GetMaxScanTokenSize() int {
	return MaxScanTokenSize
}

func SetMaxScanTokenSize(size int) {
	MaxScanTokenSize = size
}
