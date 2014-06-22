package routers

import (
	"github.com/astaxie/beego"
	"iforgetgo/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
