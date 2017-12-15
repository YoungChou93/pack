package entity

import (
	"flag"
	"github.com/astaxie/beego"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"fmt"
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

type DockerApp struct {
	Name          string
	UserName      string
	Image         string
	InstanceCount int32
}

var Newregistry Registry
var Newk8sui K8sui

var K8sclient *kubernetes.Clientset

/*func CreatePod(){
pod:=new(v1.Pod)
pod.TypeMeta=metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}
pod.ObjectMeta=metav1.ObjectMeta{Name: app.Name, Namespace: app.UserName, Labels: map[string]string{"name": app.Name}}
pod.Spec=v1.PodSpec{
	RestartPolicy: v1.RestartPolicyAlways,
	Containers: []v1.Container{
		v1.Container{
			Name:  app.Name,
			Image: app.Image,
			Ports: []v1.ContainerPort{
				v1.ContainerPort{
					ContainerPort: 9080,
					Protocol:      v1.ProtocolTCP,
				},
			},
*/ /*Resources: v1.ResourceRequirements{
	Requests: v1.ResourceList{
		v1.ResourceCPU:    resource.MustParse(app.Cpu),
		v1.ResourceMemory: resource.MustParse(app.Memory),
	},
},*/ /*
			},
		},
	}

}*/

func CreateReplicationController(app DockerApp) {
	rc := new(v1.ReplicationController)

	rcTypeMeta := metav1.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}
	rc.TypeMeta = rcTypeMeta

	rcObjectMeta := metav1.ObjectMeta{Name: app.Name, Namespace: app.UserName, Labels: map[string]string{"name": app.Name}}
	rc.ObjectMeta = rcObjectMeta

	rcSpec := v1.ReplicationControllerSpec{
		Replicas: &app.InstanceCount,
		Selector: map[string]string{
			"name": app.Name,
		},
		Template: &v1.PodTemplateSpec{
			metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.UserName,
				Labels: map[string]string{
					"name": app.Name,
				},
			},
			v1.PodSpec{
				RestartPolicy: v1.RestartPolicyAlways,
				Containers: []v1.Container{
					v1.Container{
						Name:  app.Name,
						Image: app.Image,
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								ContainerPort: 9080,
								Protocol:      v1.ProtocolTCP,
							},
						},
					},
				},
			},
		},
	}
	rc.Spec = rcSpec
	result, err :=K8sclient.CoreV1().ReplicationControllers(app.UserName).Create(rc)
	if err !=nil{
		fmt.Println(err.Error())
	}
	fmt.Println(result.Name)

}

func Setting() {
	Newregistry = Registry{beego.AppConfig.String("registryip"), beego.AppConfig.String("registryport"), beego.AppConfig.String("registryversion")}
	Newk8sui = K8sui{beego.AppConfig.String("k8sip"), beego.AppConfig.String("k8sport"), beego.AppConfig.String("k8sroute")}

	kubeconfig := flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags(Newk8sui.GetIpPort(), *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	K8sclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

}
