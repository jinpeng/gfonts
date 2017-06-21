# gfonts
Accessing googleapis including google fonts is blocked in China. But many open source projects are using google fonts, e.g. revealjs. This project is a simple standalone server to solve this problem, by retrieving and caching google fonts and serving. It is written in Go programming language.

Usage:
Modify the const in main.go:

const (
	PORT        = ":7011"
	CSS_HOST    = "fonts.sample.com"
	FONT_HOST   = "fonts.sample.com"
	FONT_FOLDER = "./static"
)

to your own server port, domain and file path for caching font files.

The server should support HTTPS. One could use Let's Encrypt to get free certificates.

$ go build
$ ./gfonts
or 
$ nohup ./gfonts &

