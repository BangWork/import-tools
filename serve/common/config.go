package common

var (
	CachePath        string
	encryptKey       string
	cookieExpireTime = 86400 * 30 // second
	MaxScanTokenSize int
	installArea      = InstallAreaAsia
)

func GetCachePath() string {
	return CachePath
}

func SetCachePath(path string) {
	CachePath = path
}
func GetInstallArea() string {
	return installArea
}

func SetInstallArea(input string) {
	installArea = input
}

func GetEncryptKey() string {
	return encryptKey
}

func SetEncryptKey(input string) {
	encryptKey = input
}

func GetCookieExpireTime() int {
	return cookieExpireTime
}

func SetCookieExpireTime(input int) {
	cookieExpireTime = input
}

func GetMaxScanTokenSize() int {
	return MaxScanTokenSize
}

func SetMaxScanTokenSize(size int) {
	MaxScanTokenSize = size
}
