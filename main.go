package main

import (
	_ "github.com/YoungChou93/pack/routers"
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/entity"
)

func main() {
	entity.Setting()
	beego.Run()
}
