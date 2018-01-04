package main

import (
	_ "github.com/YoungChou93/pack/routers"
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/controllers"
	"github.com/YoungChou93/pack/entity"
	"github.com/YoungChou93/pack/client"
)

func main() {
	entity.Setting()
	client.ClientSetting()
	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}

