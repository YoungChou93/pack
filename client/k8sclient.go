package client

import (
	"flag"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"fmt"
	"github.com/YoungChou93/pack/entity"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var K8sclient *kubernetes.Clientset

func ClientSetting() {

	kubeconfig := flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags(entity.Newk8sui.GetIpPort(), *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	K8sclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}



type K8sApp struct {
	Name          string
	UserName      string
	Image         string
	InstanceCount int32
	Port          int32
	TargetPort    int32
	NodePort      int32
	NameENV       string
	ValueENV       string
}

func CreateService(app entity.TaskMember)( *v1.Service, error){

	service := new(v1.Service)

	svTypemeta := metav1.TypeMeta{Kind: "Service", APIVersion: "v1"}
	service.TypeMeta = svTypemeta

	svObjectMeta := metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace, Labels: map[string]string{"name": app.Name}}
	service.ObjectMeta = svObjectMeta

	svServiceSpec := v1.ServiceSpec{
		Ports: []v1.ServicePort{
			v1.ServicePort{
				Name:       app.Name,
				Port:       app.Port,
				TargetPort: intstr.FromInt(int(app.TargetPort)),
				Protocol:   "TCP",
				//NodePort:   app.NodePort,
			},
		},
		Selector: map[string]string{"name": app.Name},
		Type:     v1.ServiceTypeNodePort,
	}
	service.Spec = svServiceSpec

	result, err := K8sclient.CoreV1().Services(app.Namespace).Create(service)
	if err !=nil{
		fmt.Println(err.Error())
	}

	return result,err

}


func CreateReplicationController(app entity.TaskMember,pod *v1.Pod) (*v1.ReplicationController,error){
	rc := new(v1.ReplicationController)

	rcTypeMeta := metav1.TypeMeta{Kind: "ReplicationController", APIVersion: "v1"}
	rc.TypeMeta = rcTypeMeta

	rcObjectMeta := metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace, Labels: map[string]string{"name": app.Name}}
	rc.ObjectMeta = rcObjectMeta

	rcSpec := v1.ReplicationControllerSpec{
		Replicas: &app.InstanceCount,
		Selector: map[string]string{
			"name": app.Name,
		},
		Template: &v1.PodTemplateSpec{
			ObjectMeta:metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,

				Labels: map[string]string{
					"name": app.Name,
				},
			},
			Spec:pod.Spec,
		},
	}
	rc.Spec = rcSpec
	result, err :=K8sclient.CoreV1().ReplicationControllers(app.Namespace).Create(rc)
	if err !=nil{
		fmt.Println(err.Error())
	}

	return result,err
}

func CreatePod(app entity.TaskMember)(*v1.Pod,error){
	container:=v1.Container{
		Name:  app.Name,
		Image: app.Image,
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				ContainerPort: app.TargetPort,
				Protocol:      v1.ProtocolTCP,
			},
		},

	}
	if len(app.Env)>0{
		container.Env=app.Env
	}
	pod:=new(v1.Pod)
	pod.TypeMeta=metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}
	pod.ObjectMeta=metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace, Labels: map[string]string{"name": app.Name}}
	pod.Spec=v1.PodSpec{
		RestartPolicy: v1.RestartPolicyAlways,
		Containers: []v1.Container{
			container,
		},
	}
	result, err := K8sclient.CoreV1().Pods(app.Namespace).Create(pod)
	if err !=nil{
		fmt.Println(err.Error())
	}

	return result,err
}


func RemoveReplicationController(namespace string,name string)error{
	err :=K8sclient.CoreV1().ReplicationControllers(namespace).Delete(name,&metav1.DeleteOptions{})
	if err !=nil{
		fmt.Println(err.Error())
	}
	return err
}

func RemovePod(namespace string,name string)error{
	err :=K8sclient.CoreV1().Pods(namespace).Delete(name,&metav1.DeleteOptions{})
	if err !=nil{
		fmt.Println(err.Error())
	}
	return err
}

func RemoveService(namespace string,name string)error{
	_,err:=K8sclient.CoreV1().Services(namespace).Get(name,metav1.GetOptions{})
	if err==nil {
		err = K8sclient.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			fmt.Println(err.Error())
		}
		return err
	}
	return nil

}

func GetPod(namespace string,name string)*v1.Pod{
	pod,_:=K8sclient.CoreV1().Pods(namespace).Get(name,metav1.GetOptions{})
	return pod

}


func ListReplicationController()[] entity.SimulationApp{
	namespace := "default"
	rcList, err:= K8sclient.CoreV1().ReplicationControllers(namespace).List(metav1.ListOptions{})
	if err !=nil{
		fmt.Println(err.Error())
	}

	apps:=make([]entity.SimulationApp,0)
	for _,rc :=range rcList.Items{
		app:=entity.SimulationApp{Name:rc.Name}
		apps=append(apps, app)
	}
	return apps
}

func ListPod()[] entity.SimulationApp{
	namespace := "default"
	podList, err:= K8sclient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err !=nil{
		fmt.Println(err.Error())
	}

	apps:=make([]entity.SimulationApp,0)
	for _,pod :=range podList.Items{
		app:=entity.SimulationApp{Name:pod.Name,Status:string(pod.Status.Phase)}
		apps=append(apps, app)
	}
	return apps
}


func ShowLogs(name string) string{
	namespace := "default"
	request:=K8sclient.CoreV1().Pods(namespace).GetLogs(name,&v1.PodLogOptions{})
	result:=request.Do()
	body,err:=result.Raw()
	if err !=nil{
		fmt.Println(err.Error())
	}
	return string(body)
}