package controllers

import (
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/entity"
	"time"
	"k8s.io/api/core/v1"
	"strconv"
)

var NAMESPACE string

func init()  {
	NAMESPACE="default"
}


type SimulationController struct {
	beego.Controller
}

func (c *SimulationController) TasksView() {
	c.TplName = "tasks.html"
}

func (c *SimulationController) TaskView() {
	namespace := c.GetString("namespace")
	name := c.GetString("name")
	task:=entity.App.GetTask(namespace,name)
	c.TplName = "task.html"
	c.Data["task"]=task
}

func (c *SimulationController) OneTask() {
	namespace := c.GetString("namespace")
	name := c.GetString("name")

	task:=entity.App.GetTask(namespace,name)

	c.Data["json"]=task
	c.ServeJSON()
}


func (c *SimulationController) ListTask() {
	namespace := NAMESPACE
	tasks,_:=entity.App.GetTasks(namespace)
	c.Data["json"] = &tasks
	c.ServeJSON()
}

func (c *SimulationController) AddTask() {
	name := c.GetString("name")
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	result := entity.Result{Success:true}
	err:=entity.App.AddTask(NAMESPACE,entity.NewTask(name,NAMESPACE,timeNow))
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *SimulationController) RemoveTask() {
	namespace := c.GetString("namespace")
	name := c.GetString("name")
	result := entity.Result{Success:true}
	err:=entity.App.RemoveTask(namespace,name)
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
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

	taskname:=c.GetString("taskname")
	namespace:=c.GetString("namespace")

	task:=entity.App.GetTask(namespace,taskname)

	name := c.GetString("name")
	image := c.GetString("image")
	types, _:= c.GetInt("type")
	result := entity.Result{Success:true}

	member:= entity.TaskMember{Name: name, Namespace: task.Namespace, Image: image, InstanceCount: 1, Types:types}
	switch types {
	case 1:
		port, _ := c.GetInt32("port")
		member.Port=port
		member.NodePort=port
		member.TargetPort=60400

	case 2:
		port, _ := c.GetInt32("port")
		member.Port=port
		member.NodePort=port
		member.TargetPort=22
		member.Env=GetEnv(c)

	case 3:
		member.TargetPort=22
		member.Env=GetEnv(c)
	}
	err:=entity.App.AddTaskMember(task, member)
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
	}else {

	}
	c.Data["json"] = &result
	c.ServeJSON()
}


func (c *SimulationController) RemoveMember() {
	name := c.GetString("name")
	membername := c.GetString("membername")
	namespace := c.GetString("namespace")

	err:=entity.App.RemoveTaskMember(namespace,name,membername)

	result := entity.Result{Success:true}
	if err !=nil{
		result.Success=false
		result.Reason=err.Error()
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *SimulationController) GetLog() {
	name := c.GetString("name")
	namespace := c.GetString("namespace")
	logs,_:=entity.App.GetLog(namespace,name)

	result := entity.Result{Success:true,Reason:logs}
	c.Data["json"] = &result
	c.ServeJSON()
}