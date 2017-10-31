package routers

import (
	"github.com/YoungChou93/pack/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/listimages", &controllers.ListImageController{})
}
