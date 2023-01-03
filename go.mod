module github.com/bangwork/import-tools

go 1.16

replace (
	github.com/rwcarlsen/goexif/exif => ./serve/external/goexif/exif
	github.com/rwcarlsen/goexif/tiff => ./serve/external/goexif/tiff
)

require (
	github.com/beevik/etree v1.1.0
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-contrib/i18n v0.0.1
	github.com/gin-gonic/gin v1.8.1
	github.com/juju/errors v0.0.0-20220203013757-bd733f3c86b9
	github.com/juju/testing v1.0.2 // indirect
	github.com/mozillazg/go-pinyin v0.19.0
	github.com/pelletier/go-toml/v2 v2.0.1
	github.com/rwcarlsen/goexif/exif v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.2.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	golang.org/x/text v0.5.0
)
