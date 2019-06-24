package routers

import (
	"github.com/astaxie/beego"
	"ppt2html/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/upload", &controllers.UploadController{})
}
