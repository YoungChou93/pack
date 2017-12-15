package main

import (
	_ "github.com/YoungChou93/pack/routers"
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/entity"
	"github.com/YoungChou93/pack/controllers"

)

func main() {
	entity.Setting()
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()


	/*err:=util.Untar("/home/zhouyang/pack/fedtest1.0/build_certi.tar","/home/zhouyang/pack/fedtest1.0/")

	if err!=nil{
		fmt.Println(err.Error())
	}*/

	//app:=entity.DockerApp{"tomcat2","default","tomcat",1}
	//entity.CreateReplicationController(app)
}

