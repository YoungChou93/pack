package entity

import (
	"errors"
	"flag"
	"github.com/YoungChou93/pack/client"
	"github.com/astaxie/beego"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type TreeNode struct {
	Id    int    `json:"id"`
	Text  string `json:"text"`
	State State  `json:"state"`
}

type State struct {
	Selected bool `json:"selected"`
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

type Result struct {
	Success bool
	Reason  string
}


//仿真成员 或者 工具 或者 模型
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
	Cmd           []string
	Service       *v1.Service
	Rc            *v1.ReplicationController
	Pod           *v1.Pod
}

func (this *TaskMember) GetK8sApp(taskname string) client.K8sApp {
	app := client.K8sApp{TaskName: taskname, Name: this.Name, Namespace: this.Namespace,
		Image: this.Image, InstanceCount: this.InstanceCount, Port: this.Port,
		TargetPort: this.TargetPort, NodePort: this.NodePort, Cmd: this.Cmd, Types: this.Types}
	if len(this.Env) > 0 {
		app.Env = this.Env
	}
	return app
}

//仿真任务
type Task struct {
	Name      string
	Namespace string
	Time      string
	Members   []TaskMember
}

func NewTask(name string, namespace string, time string) Task {
	members := make([]TaskMember, 0)
	return Task{Name: name, Namespace: namespace, Members: members, Time: time}
}

func (this *Task) AddTaskMember(member TaskMember) {
	this.Members = append(this.Members, member)
}

func (this *Task) RemoveTaskMember(name string) (TaskMember, error) {
	index := -1
	for i, member := range this.Members {
		if member.Name == name {
			index = i
		}
	}
	if index != -1 {
		taskMember := this.Members[index]
		if len(this.Members) == 1 {
			this.Members = make([]TaskMember, 0)
		} else {
			this.Members = append(this.Members[:index], this.Members[index+1:]...)
		}
		return taskMember, nil
	}
	return TaskMember{}, errors.New("error name")
}

var App Application

var Newk8sui K8sui

func Setting() {

	client.RegistrySetting()
	Newk8sui = K8sui{beego.AppConfig.String("k8sip"), beego.AppConfig.String("k8sport"), beego.AppConfig.String("k8sroute")}
	kubeconfig := flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags(Newk8sui.GetIpPort(), *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	K8sclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	kclient := client.KubernetesClient{K8sclient}

	App = NewApplication(kclient)
	App.Read()

}
