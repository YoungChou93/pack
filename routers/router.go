package routers

import (
	"github.com/YoungChou93/pack/controllers"
	"github.com/astaxie/beego"
)

func init() {

	//封装页面
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/encapsulationview", &controllers.UploadController{},"get:EncapsulationView")
	beego.Router("/encapsulation", &controllers.UploadController{},"Post:Encapsulation")

	//本地镜像
	beego.Router("/localimage", &controllers.ListImageController{})
	beego.Router("/localimage/list", &controllers.ListImageController{},"post:List")
	beego.Router("/localimage/remove", &controllers.ListImageController{},"post:Remove")
	beego.Router("/localimage/run", &controllers.ListImageController{},"post:Run")

	//镜像仓库
	beego.Router("/registry", &controllers.RegistryController{},)
	beego.Router("/registry/list", &controllers.RegistryController{},"post:List")
	beego.Router("/registry/imagepush", &controllers.RegistryController{},"post:Push")
	beego.Router("/registry/imagepull", &controllers.RegistryController{},"post:Pull")

	//镜像仓库管理
	beego.Router("/registrymanage", &controllers.RegistryController{},"get:RegistryView")
	beego.Router("/registry/addregistry", &controllers.RegistryController{},"post:AddRegistry")
	beego.Router("/registry/listregistry", &controllers.RegistryController{},"post:ListRegistry")
	beego.Router("/registry/majorregistry", &controllers.RegistryController{},"post:MajorRegistry")
	beego.Router("/registry/deleteregistry", &controllers.RegistryController{},"post:DeleteRegistry")

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

	//仿真开发
	beego.Router("/simulation/toolview", &controllers.SimulationController{},"get:ToolView")
	beego.Router("/simulation/createtool", &controllers.SimulationController{},"post:CreateTool")

	//文件管理
	beego.Router("/filemanage", &controllers.FileController{},"get:FileView")
	beego.Router("/getpath", &controllers.FileController{},"post:GetPath")
	beego.Router("/file/uploadfile", &controllers.FileController{},"post:UploadFile")
	beego.Router("/file/deletefile", &controllers.FileController{},"post:DeleteFile")
	beego.Router("/file/download", &controllers.FileController{},"get:Download")

	//用户管理
	beego.Router("/usermanage", &controllers.UserController{},"get:UserView")
	beego.Router("/user/adduser", &controllers.UserController{},"post:AddUser")
	beego.Router("/user/listuser", &controllers.UserController{},"post:ListUser")
	beego.Router("/user/disableuser", &controllers.UserController{},"post:DisableUser")
	beego.Router("/user/enableuser", &controllers.UserController{},"post:EnableUser")
	beego.Router("/user/getuserright", &controllers.UserController{},"post:GetUserRight")
	beego.Router("/user/setuserright", &controllers.UserController{},"post:SetUserRight")
	beego.Router("/login", &controllers.UserController{},"get:LoginView")
	beego.Router("/login", &controllers.UserController{},"post:Login")
	beego.Router("/logout", &controllers.UserController{},"get:Logout")
	beego.Router("/modifypassword", &controllers.UserController{},"get:ModifyPasswordView")
	beego.Router("/modifypassword", &controllers.UserController{},"post:ModifyPassword")

	//设置
	beego.Router("/setting", &controllers.SetController{})
	beego.Router("/setting/k8sui", &controllers.SetController{},"post:SetK8sui")

}
