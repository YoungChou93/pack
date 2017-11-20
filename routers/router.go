package routers

import (
	"github.com/YoungChou93/pack/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
    beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/localimage/list", &controllers.ListImageController{})
	beego.Router("/localimage/remove", &controllers.ListImageController{},"post:Remove")
	beego.Router("/setting", &controllers.SetController{})
	beego.Router("/registry", &controllers.RegistryController{})
	beego.Router("/registry/imagepush", &controllers.RegistryController{},"post:Push")
}
