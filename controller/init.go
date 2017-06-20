package controller

import (
	"github.com/astaxie/beego"
)

func Init() {
	InitManager()
	for _, handler := range allUrls{
		beego.Router(handler.url, handler.controller)
	}
}