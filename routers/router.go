package routers

import (
	"ppt2html/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/upload", &controllers.UploadController{})
}
