package controllers

import (
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/entity"
	"time"
	"k8s.io/api/core/v1"
	"strconv"
	"os"
	"fmt"
	"io"
	"strings"
	"github.com/YoungChou93/pack/client"
	"github.com/YoungChou93/pack/database"
)


type SimulationController struct {
	beego.Controller
}

func (c *SimulationController) TasksView() {
	c.TplName = "tasks.html"
}

func (c *SimulationController) ToolView() {
	c.TplName = "develop.html"
}

func (c *SimulationController) TaskView() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username

	name := c.GetString("name")
	task,_:=entity.App.GetTask(namespace,name)
	c.TplName = "task.html"
	c.Data["task"]=task
}

func (c *SimulationController) OneTask() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username

	name := c.GetString("name")

	task,_:=entity.App.GetTask(namespace,name)

	c.Data["json"]=task
	c.ServeJSON()
}


func (c *SimulationController) ListTask() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username
	tasks,_:=entity.App.GetTasks(namespace)
	c.Data["json"] = &tasks
	c.ServeJSON()
}

func (c *SimulationController) AddTask() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username


	name := c.GetString("name")
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	result := entity.Result{Success:true}
	if name=="tool"{
		result.Success=false
		result.Reason="tool 为系统保留名称"
	}else {
		err := entity.App.AddTask(namespace, entity.NewTask(name, namespace, timeNow))
		if err != nil {
			result.Success = false
			result.Reason = err.Error()
			//错误日志
			beego.Error(err.Error())
		}
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *SimulationController) RemoveTask() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username

	name := c.GetString("name")
	result := entity.Result{Success:true}
	err:=entity.App.RemoveTask(namespace,name)
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
		//错误日志
		beego.Error(err.Error())
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

//得到环境变量
func GetEnv(c *SimulationController)[]v1.EnvVar{
	envs:=make([]v1.EnvVar,0)
	for i:=0;;i++ {
		envname := "env[" + strconv.Itoa(i) + "][name]"
		envvalue := "env[" + strconv.Itoa(i) + "][value]"
		nameenv := c.GetString(envname)
		valueenv := c.GetString(envvalue)
		if len(nameenv)<=0{
			break
		}
		envVar:=v1.EnvVar{Name:nameenv,Value:valueenv}
		envs=append(envs,envVar)
	}
	return envs
}

//运行仿真成员，也可以叫addTaskMember
func (c *SimulationController) Run() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username

	taskname:=c.GetString("taskname")

	task,_:=entity.App.GetTask(namespace,taskname)

	name := c.GetString("name")
	image := c.GetString("image")
	cmd := c.GetString("cmd")
	types, _:= c.GetInt("type")
	result := entity.Result{Success:true}

	var cmds []string
	if len(cmd)>0{
		cmds=strings.Split(cmd," -")
		for i,_ :=range cmds{
			if i>0 {
				cmds[i]="-"+cmds[i]
			}
		}
	}

	if !strings.Contains(image,client.MajorRegistry.GetIpPort()){
		image=client.MajorRegistry.GetIpPort()+"/"+image
	}

	member:= entity.TaskMember{Name: name, Namespace: task.Namespace, Image: image, InstanceCount: 1, Cmd:cmds,Types:types}
	switch types {
	case 1:
		port, _ := c.GetInt32("port")
		member.Port=port
		member.TargetPort=60400

	case 2:
		port, _ := c.GetInt32("port")
		member.Port=port
		member.TargetPort=4200
		envs:=GetEnv(c)
		envVar1:=v1.EnvVar{Name:"SIAB_PASSWORD",Value:"123456"}
		envVar2:=v1.EnvVar{Name:"SIAB_SUDO",Value:"true"}
		envs=append(envs,envVar1)
		envs=append(envs,envVar2)
		member.Env=envs


	case 3:
		member.TargetPort=22
		member.Env=GetEnv(c)

	case client.TYPE_VNC:
		port, _ := c.GetInt32("port")
		member.Port=port
		member.TargetPort=6080
		member.Env=GetEnv(c)
	}
	err:=entity.App.AddTaskMember(task, member)
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
		//错误日志
		beego.Error(err.Error())
	}
	c.Data["json"] = &result
	c.ServeJSON()
}


func (c *SimulationController) RemoveMember() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username

	name := c.GetString("name")
	membername := c.GetString("membername")


	err:=entity.App.RemoveTaskMember(namespace,name,membername)

	result := entity.Result{Success:true}
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
		//错误日志
		beego.Error(err.Error())
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *SimulationController) GetLog() {
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username

	name := c.GetString("name")

	logs,_:=entity.App.GetLog(namespace,name)

	result := entity.Result{Success:true,Reason:logs}
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *SimulationController) UploadFile(){
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username


	file, header, _ := c.GetFile("fedfile")
	name:= c.GetString("name")

	nfspath := beego.AppConfig.String("nfspath")
	result := entity.Result{Success:true}
	task,_:=entity.App.GetTask(namespace,name)
	if len(task.Name)>0 {
		//保存待封装软件
		dirpath := nfspath + "/"+task.Namespace+"/"+task.Name
		os.MkdirAll(dirpath, os.ModePerm)
		f, err := os.OpenFile(dirpath+"/"+header.Filename, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			fmt.Println("文件打开失败")
			result.Reason=err.Error()
			//错误日志
			beego.Error(err.Error())

		}
		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Println("文件保存失败" + err.Error())
			result.Reason=err.Error()
			//错误日志
			beego.Error(err.Error())

		}

		f.Close()
	}else{
		result.Reason="任务不存在"
	}
	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *SimulationController) CreateTool(){
	//namespace := c.GetString("namespace")
	user:=c.GetSession("user").(database.User)
	namespace :=user.Username
	dirpath:=beego.AppConfig.String("nfspath") +user.Username+"/tool"
	os.MkdirAll(dirpath, os.ModePerm)


	name := c.GetString("name")
	image := c.GetString("image")
	task,err:=entity.App.GetTask(namespace,"tool")
	if err!=nil{
		tNow := time.Now()
		timeNow := tNow.Format("2006-01-02 15:04:05")
		entity.App.AddTask(namespace,entity.NewTask("tool",namespace,timeNow))
		task,_=entity.App.GetTask(namespace,"tool")
		//错误日志
		beego.Error(err.Error())
	}
	result := entity.Result{Success:true}

	if !strings.Contains(image,client.MajorRegistry.GetIpPort()){
		image=client.MajorRegistry.GetIpPort()+"/"+image
	}
	member:= entity.TaskMember{Name: name, Namespace: namespace, Image: image,Port:8080, InstanceCount: 1, Types:client.TYPE_TOOL,TargetPort:80}

	err=entity.App.AddTaskMember(task, member)
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
		//错误日志
		beego.Error(err.Error())
	}else {

	}
	c.Data["json"] = &result
	c.ServeJSON()
}