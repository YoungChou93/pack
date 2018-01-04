package entity

import (
	"container/list"
	"errors"
	"github.com/astaxie/beego"
	"k8s.io/api/core/v1"
)

//镜像仓库地址
type Registry struct {
	Ipaddr  string
	Port    string
	Version string
}

func (this *Registry) GetUrl() string {
	return "http://" + this.Ipaddr + ":" + this.Port + "/" + this.Version
}

func (this *Registry) GetIpPort() string {
	return this.Ipaddr + ":" + this.Port
}

//k8sui地址
type K8sui struct {
	Ipaddr string
	Port   string
	Route  string
}

func (this *K8sui) GetUrl() string {
	return "http://" + this.Ipaddr + ":" + this.Port + "/" + this.Route
}

func (this *K8sui) GetIpPort() string {
	return this.Ipaddr + ":" + this.Port
}

//构造镜像仓库镜像格式
type RepositoriesInfo struct {
	Repositories []string
}

type Image struct {
	Name string
	Tags []string
}

type ResgitryImage struct {
	Images []Image
}

type Result struct {
	Success bool
	Reason  string
}

type SimulationApp struct {
	Name   string
	Status string
}

type SimulationApps struct {
	Apps []SimulationApp
}

type Task struct {
	Name      string
	Namespace string
	Time      string
	Members   []TaskMember
}

type TaskMember struct {
	Name          string
	Namespace     string
	Types         int
	Image         string
	InstanceCount int32
	Port          int32
	TargetPort    int32
	NodePort      int32
	Env           []v1.EnvVar
	Service       *v1.Service
	Rc            *v1.ReplicationController
	Pod           *v1.Pod
}

func (this *Task) AddTaskMember(member TaskMember) {
	this.Members = append(this.Members, member)
}

func NewTask(name string, namespace string, time string) Task {
	members := make([]TaskMember, 0)
	return Task{Name: name, Namespace: namespace, Members: members, Time: time}
}

type Application struct {
	taskMap map[string]*list.List
	findMap map[string]*list.Element
}

func NewApplication() Application {
	taskmap := make(map[string]*list.List)
	findmap := make(map[string]*list.Element)
	return Application{taskMap: taskmap, findMap: findmap}
}

func (this *Application) AddTask(namespace string, task Task) error {
	key := namespace + "#" + task.Name
	if _, ok := this.findMap[key]; ok {
		return errors.New("The task has already existed !")
	} else {
		var e *list.Element
		if tasks, ok := this.taskMap[namespace]; ok {
			e = tasks.PushBack(task)
		} else {
			tasks := list.New()
			e = tasks.PushBack(task)
			this.taskMap[namespace] = tasks
		}
		this.findMap[key] = e
	}
	return nil
}

func (this *Application) RemoveTask(namespace string, name string)( Task,error) {
	key := namespace + "#" + name
	if e, ok := this.findMap[key]; ok {
		delete(this.findMap, key)
		if tasks, ok := this.taskMap[namespace]; ok {
			task := e.Value.(Task)
			tasks.Remove(e)
			return task,nil
		} else {
			return Task{},errors.New("error namespace")
		}
	} else {
		return Task{},errors.New("error name")
	}

}

func (this *Application) GetTasks(namespace string) ([]Task, error) {
	if tasks, ok := this.taskMap[namespace]; ok {
		return this.ListToTasks(tasks), nil
	} else {
		return nil, errors.New("Namespace does not exist")
	}
}

func (this *Application) ListToTasks(list *list.List) []Task {
	i := 0
	tasks := make([]Task, list.Len())
	for e := list.Front(); e != nil; e = e.Next() {
		tasks[i] = e.Value.(Task)
		i++
	}
	return tasks
}

func (this *Application) GetTask(namespace string, name string) Task {
	key := namespace + "#" + name
	if e, ok := this.findMap[key]; ok {
		task := e.Value.(Task)
		return task
	}
	return Task{}
}

func (this *Application) ModifyTask(task Task) error {
	key := task.Namespace + "#" + task.Name
	if e, ok := this.findMap[key]; ok {
		if tasks, ok := this.taskMap[task.Namespace]; ok {
			tasks.Remove(e)
			e = tasks.PushBack(task)
		}
		this.findMap[key] = e
	} else {
		return errors.New("task does not exist")
	}
	return nil
}

var App Application
var Newregistry Registry
var Newk8sui K8sui

func Setting() {
	Newregistry = Registry{beego.AppConfig.String("registryip"), beego.AppConfig.String("registryport"), beego.AppConfig.String("registryversion")}
	Newk8sui = K8sui{beego.AppConfig.String("k8sip"), beego.AppConfig.String("k8sport"), beego.AppConfig.String("k8sroute")}
	App = NewApplication()
}
