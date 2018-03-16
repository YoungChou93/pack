package main



import (
	_ "github.com/YoungChou93/pack/routers"
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/controllers"
	"github.com/YoungChou93/pack/entity"
	"github.com/YoungChou93/pack/database"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

func main() {
	database.Dbsetting()
	entity.Setting()
	beego.ErrorController(&controllers.ErrorController{})
	beego.BConfig.WebConfig.Session.SessionOn = true

	//验证用户是否已经登录
	var FilterUser = func(ctx *context.Context) {
		user, ok := ctx.Input.Session("user").(database.User)
		if !ok && ctx.Request.RequestURI != "/login" {
			ctx.Redirect(302, "/login")
		}else{
			if ctx.Request.RequestURI != "/login" {
				beego.Informational(user.Username + " call " + ctx.Request.RequestURI)
			}
		}
	}

	beego.InsertFilter("/*",beego.BeforeRouter,FilterUser)

	beego.SetLogger(logs.AdapterMultiFile,`{"filename":"logs/project.log","separate":["error", "info","debug"],"level":7,"daily":true,"maxdays":30}`)
	beego.SetLogFuncCall(true)

	beego.Run()
}


