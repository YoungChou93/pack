package routers

import (
	"github.com/YoungChou93/pack/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//封装页面
	beego.Router("/", &controllers.IndexController{})
    beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/encapsulationview", &controllers.UploadController{},"get:EncapsulationView")
	beego.Router("/encapsulation", &controllers.UploadController{},"Post:Encapsulation")

	//本地镜像
	beego.Router("/localimage", &controllers.ListImageController{})
	beego.Router("/localimage/list", &controllers.ListImageController{},"post:List")
	beego.Router("/localimage/remove", &controllers.ListImageController{},"post:Remove")

	//镜像仓库
	beego.Router("/registry", &controllers.RegistryController{},)
	beego.Router("/registry/list", &controllers.RegistryController{},"post:List")
	beego.Router("/registry/imagepush", &controllers.RegistryController{},"post:Push")
	beego.Router("/registry/imagepull", &controllers.RegistryController{},"post:Pull")

    //仿真任务
	beego.Router("/simulation/tasksview", &controllers.SimulationController{},"get:TasksView")
	beego.Router("/simulation/addtask", &controllers.SimulationController{},"post:AddTask")
	beego.Router("/simulation/removetask", &controllers.SimulationController{},"post:RemoveTask")
	beego.Router("/simulation/listtask", &controllers.SimulationController{},"post:ListTask")
	beego.Router("/simulation/uploadfile", &controllers.SimulationController{},"post:UploadFile")

	//仿真成员
	beego.Router("/simulation/taskview", &controllers.SimulationController{},"post:TaskView")
	beego.Router("/simulation/onetask", &controllers.SimulationController{},"post:OneTask")
	beego.Router("/simulation/run", &controllers.SimulationController{},"post:Run")
	beego.Router("/simulation/log", &controllers.SimulationController{},"post:GetLog")
	beego.Router("/simulation/removemember", &controllers.SimulationController{},"post:RemoveMember")

    //设置
	beego.Router("/setting", &controllers.SetController{})
	beego.Router("/setting/registry", &controllers.SetController{},"post:SetRegistry")
	beego.Router("/setting/k8sui", &controllers.SetController{},"post:SetK8sui")
}
