package main

import (
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-assets"
)

func main() {
	r := gin.New()
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}

	r.SetHTMLTemplate(t)
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "/html/bar.tmpl", gin.H{"Foo": "World"})
	})
	r.GET("/bar", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "/html/bar.tmpl", gin.H{"Bar": "World"})
	})

	r.Run()

}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}

		h, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Assets
var _Assetsbfa8d115ce0617d89507412d5393a462f8e9b003 = "<!doctype html>\n<body>\n  <p>Can you see this? → {{.Bar}}</p>\n</body>"
var _Assets3737a75b5254ed1f6d588b40a3449721f9ea86c2 = "<!doctype html>\n<body>\n  <p>Hello, {{.Foo}}</p>\n</body>"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": {"html"}, "/html": {"bar.tmpl", "index.tmpl"}}, map[string]*assets.File{
	"/html": {
		Path:     "/html",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1715698523, 1715698523114378388),
		Data:     nil,
	}, "/html/bar.tmpl": {
		Path:     "/html/bar.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1715698530, 1715698530473692611),
		Data:     []byte(_Assetsbfa8d115ce0617d89507412d5393a462f8e9b003),
	}, "/html/index.tmpl": {
		Path:     "/html/index.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1715698539, 1715698539918794004),
		Data:     []byte(_Assets3737a75b5254ed1f6d588b40a3449721f9ea86c2),
	}, "/": {
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1715698512, 1715698512707941562),
		Data:     nil,
	}}, "")
