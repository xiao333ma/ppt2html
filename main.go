package main

import (
	"ppt2html/common"
	_ "ppt2html/routers"
	"github.com/astaxie/beego"
)

func main() {

	beego.AddAPPStartHook()
	beego.SetStaticPath("/static", common.GetStaticPath())
	beego.SetViewsPath (common.GetViewsPath())
	beego.Run()
}

