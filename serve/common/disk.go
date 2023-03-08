package common

func CheckDiskAvailable(attachmentSize, ONESFileCapacitySize int64) bool {
	attachmentSize = attachmentSize / 1024 / 1024 / 1024
	return attachmentSize*2 < ONESFileCapacitySize
}
