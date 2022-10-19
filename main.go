package main

import (
	"encoding/base64"
	"fmt"
	"github.com/wenlaizhou/middleware"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const codeTpl = `package ${package}

import (
	"encoding/base64"
	"github.com/wenlaizhou/middleware"
)

func init() {
	${code}
}
`

const varTpl = "var ${name}, _ = base64.StdEncoding.DecodeString(`${value}`)\n" + `
	middleware.RegisterHandler("${path}", func(context middleware.Context) {
		context.OK("${contentType}", ${name})
	})`

var backg, _ = base64.StdEncoding.DecodeString(`IyBkaXN0IGJ1aWxkZXINCg0KPiDliY3nq69kaXN0IOaJk+WMheaIkG1pZGRsZXdhcmXmnI3liqEsIOWPr+ebtOaOpeWwhuW
Jjeerr+S7o+eggee8luivkeS4umdvbGFuZ+S7o+eggSwg5b2i5oiQ5Y2V5Y+v5omn6KGM5paH5Lu2`)

type DistVar struct {
	Name        string
	Value       string
	Path        string
	ContentType string
}

type PageCode struct {
	Package string
	Code    string
}

func init() {

}

func BuildFiles(pkg string, distPath string) string {
	page := PageCode{
		Package: pkg,
	}
	codeBuilder := strings.Builder{}

	if err := filepath.Walk(distPath, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		newPath := strings.ReplaceAll(path, string(os.PathSeparator), "/")
		subPath := strings.Replace(newPath, fmt.Sprintf("%v/", distPath), "", 1)
		if info.IsDir() {
			if indexContent, err := os.ReadFile(fmt.Sprintf("%v%vindex.html", path, string(os.PathSeparator))); err == nil {
				name := fmt.Sprintf("var_rand_%v", rand.Int()) + strings.ReplaceAll(strings.ReplaceAll(info.Name(), ".", "_"), "-", "_")
				codeBuilder.WriteString(middleware.StringFormatStructs(varTpl, DistVar{
					Name:        name,
					Value:       base64.StdEncoding.EncodeToString(indexContent),
					Path:        fmt.Sprintf("/%v", subPath),
					ContentType: middleware.Html,
				}))
				indexCode := fmt.Sprintf(`
	middleware.RegisterHandler("/%v" ,func(context middleware.Context) {
		context.AddCacheHeader(3600 * 24 * 30)
		context.OK(middleware.Html, %v)
	})
`, subPath, name)
				indexCode2 := fmt.Sprintf(`
	middleware.RegisterHandler("/%v/index.html" ,func(context middleware.Context) {
		context.AddCacheHeader(3600 * 24 * 30)
		context.OK(middleware.Html, %v)
	})
`, subPath, name)
				codeBuilder.WriteString(indexCode)
				codeBuilder.WriteString(indexCode2)
			} else {
				println(path)
				println("没有index.html")
				println(err.Error())
			}

			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			println(err.Error())
			return nil
		}

		name := fmt.Sprintf("var_rand_%v", rand.Int()) + strings.ReplaceAll(strings.ReplaceAll(info.Name(), ".", "_"), "-", "_")
		if subPath == "index.html" {
			name = "index_html"
		} else {
			if info.Name() == "index.html" {
			}
		}
		url := fmt.Sprintf("/%v", subPath)
		contentType := http.DetectContentType(content)
		switch middleware.Ext(path) {
		case ".js":
		case ".jsx":
			contentType = middleware.Js
			break
		case ".css":
			contentType = middleware.Css
			break
		case ".json":
			contentType = middleware.Json
			break
		case ".html":
			contentType = middleware.Html
			break
		}
		codeBuilder.WriteString(middleware.StringFormatStructs(varTpl, DistVar{
			Name:        name,
			Value:       base64.StdEncoding.EncodeToString(content),
			Path:        url,
			ContentType: contentType,
		}))

		codeBuilder.WriteString("\n")
		return nil
	}); err != nil {
		println(err.Error())
	}
	codeBuilder.WriteString(`
	middleware.RegisterIndex(func(context middleware.Context) {
		context.AddCacheHeader(3600 * 24 * 30)
		context.OK(middleware.Html, index_html)
	})
`)

	page.Code = codeBuilder.String()
	return middleware.StringFormatStructs(codeTpl, page)
}

func main() {

	if err := ioutil.WriteFile("code.go", []byte(BuildFiles("main", "dist")), os.ModePerm); err != nil {
		println(err.Error())
	}

	// middleware.RegisterDefaultIndex(middleware.DefaultIndexStruct{
	// 	Title:              "DistBuilder",
	// 	BackgroundUrl:      "/backg",
	// 	HeaderLinks:        nil,
	// 	CenterContentLines: []string{"<h2>hello world</h2>", "123"},
	// 	Buttons: []middleware.DefaultIndexButton{
	// 		{
	// 			Text: "进入系统",
	// 			Link: "/",
	// 		}, {
	// 			Text: "文档",
	// 			Link: "/",
	// 		}, {
	// 			Text: "code",
	// 			Link: "/",
	// 		},
	// 	},
	// 	PoweredBy:     "middleware",
	// 	FooterLinks:   nil,
	// 	ExtendedStyle: "",
	// 	EnableSwagger: false,
	// })

	middleware.StartServer("", 8080)
}
