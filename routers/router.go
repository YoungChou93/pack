package routers

import (
	"github.com/YoungChou93/pack/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
    beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/localimage", &controllers.ListImageController{})
	beego.Router("/localimage/list", &controllers.ListImageController{},"post:List")
	beego.Router("/localimage/remove", &controllers.ListImageController{},"post:Remove")
	beego.Router("/registry", &controllers.RegistryController{})
	beego.Router("/registry/imagepush", &controllers.RegistryController{},"post:Push")
	beego.Router("/registry/imagepull", &controllers.RegistryController{},"post:Pull")
	beego.Router("/setting", &controllers.SetController{})
	beego.Router("/setting/registry", &controllers.SetController{},"post:SetRegistry")
	beego.Router("/setting/k8sui", &controllers.SetController{},"post:SetK8sui")
}
