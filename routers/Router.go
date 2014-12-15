package routers

import (
	"github.com/astaxie/beego"
	"notifier/controllers"
)

func init() {
	// роутеры запросов, направляют запросы на обработку в соответствующий контроллер
	beego.Router("/", &controllers.MainController{})
	beego.Router("/test", &controllers.TestController{})
}
