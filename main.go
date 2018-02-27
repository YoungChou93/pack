package main



import (
	_ "github.com/YoungChou93/pack/routers"
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/controllers"
	"github.com/YoungChou93/pack/entity"
	"github.com/YoungChou93/pack/database"
)

func main() {
	database.Dbsetting()
	entity.Setting()
	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}

