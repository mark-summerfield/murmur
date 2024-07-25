module github.com/mark-summerfield/murmur

go 1.22.5

replace github.com/mark-summerfield/clip => /home/mark/app/golib/clip

replace github.com/mark-summerfield/ufile => /home/mark/app/golib/ufile

replace github.com/mark-summerfield/uterm => /home/mark/app/golib/uterm

replace github.com/mark-summerfield/utext => /home/mark/app/golib/utext

require (
	github.com/mark-summerfield/clip v1.5.0
	github.com/mark-summerfield/ufile v0.0.0-00010101000000-000000000000
)

require (
	github.com/kopoli/go-terminal-size v0.0.0-20170219200355-5c97524c8b54 // indirect
	github.com/mark-summerfield/uterm v1.0.0 // indirect
	github.com/mark-summerfield/utext v1.0.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.22.0 // indirect
)
