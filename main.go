package main

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"ppt2html/common"
	_ "ppt2html/routers"
)

func main() {
	moveScript("/Users/liuche/GO_WorkSpace/src/ppt2html/ppt2html.scpt",common.GetScriptPath() + "/" + "ppt2html.scpt")
	beego.AddAPPStartHook()
	beego.SetStaticPath("/static", common.GetStaticPath())
	beego.SetViewsPath(common.GetViewsPath())
	beego.Run()
}

func moveScript(src_path string,target_path string)  {
	if _, err := os.Stat(src_path); err == nil {
		if _, err := os.Stat(target_path); err != nil {
			input, err1 := ioutil.ReadFile(src_path)
			if err1 == nil {
				ioutil.WriteFile(target_path, input, 0777)
			}
		}
	}
}
