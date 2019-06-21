package main

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"ppt2html/common"
	_ "ppt2html/routers"
)

func main() {
	src_path := "/Users/liuche/GO_WorkSpace/src/ppt2html/ppt2html.scpt"
	target_path := common.GetScriptPath()+"/"+"ppt2html.scpt"
	if _, err := os.Stat(src_path); err == nil {
		if _, err := os.Stat(target_path); err != nil {
			input, err1 := ioutil.ReadFile(src_path)
			if err1 == nil {
				ioutil.WriteFile(target_path, input, 0777)
			}
		}
	}

	beego.AddAPPStartHook()
	beego.SetStaticPath("/static", common.GetStaticPath())
	beego.SetViewsPath (common.GetViewsPath())
	beego.Run()
}


