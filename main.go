package main

import (
	"encoding/base64"
	"github.com/wenlaizhou/middleware"
	"io/ioutil"
	"os"
)

var backg, _ = base64.StdEncoding.DecodeString(`IyBkaXN0IGJ1aWxkZXINCg0KPiDliY3nq69kaXN0IOaJk+WMheaIkG1pZGRsZXdhcmXmnI3liqEsIOWPr+ebtOaOpeWwhuW
Jjeerr+S7o+eggee8luivkeS4umdvbGFuZ+S7o+eggSwg5b2i5oiQ5Y2V5Y+v5omn6KGM5paH5Lu2`)

func main() {
	if code, err := middleware.DistFrontend2Code("main", "build", ""); err == nil {
		if err := ioutil.WriteFile("code.go", []byte(code), os.ModePerm); err != nil {
			println(err.Error())
		}
	} else {
		println(err.Error())
	}

	// middleware.RegisterHandler("/backg", func(context middleware.Context) {
	// 	context.ServeFile("bg.png")
	// })
	//
	// middleware.RegisterDefaultIndex(middleware.DefaultIndexStruct{
	// 	Title: "DistBuilder",
	// 	HeaderLinks: []middleware.DefaultIndexLink{
	// 		{
	// 			Text: "Welcome",
	// 			Link: "",
	// 		},
	// 	},
	// 	CenterContentLines: []string{"<h2>hello world</h2>", "123"},
	// 	Buttons: []middleware.DefaultIndexLink{
	// 		{
	// 			Text: "Enter",
	// 			Link: "/web",
	// 		}, {
	// 			Text: "Doc",
	// 			Link: "/doc",
	// 		},
	// 	},
	// 	PoweredBy: "middleware",
	// 	FooterLinks: []middleware.DefaultIndexLink{
	// 		{
	// 			Text: "A",
	// 			Link: "",
	// 		},
	// 	},
	// 	EnableSwagger: false,
	// })

	middleware.StartServer("", 8080)
}
