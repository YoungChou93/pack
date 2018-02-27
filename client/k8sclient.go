package client

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"github.com/astaxie/beego"
)

const (
	TYPE_RTIG  = 1
	TYPE_FIRST_FED = 2
	TYPE_OTHER_FED   = 3
	TYPE_VNC=4
)

type KubernetesClient struct {
	K8sclient *kubernetes.Clientset
}

type K8sApp struct {
	TaskName      string
	Name          string
	Namespace     string
	Image         string
	InstanceCount int32
	Port          int32
	TargetPort    int32
	NodePort      int32
	Env           []v1.EnvVar
	Cmd           []string
	Types         int
}

func (this * KubernetesClient)CreateService(app K8sApp)( *v1.Service, error){

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

	result, err := this.K8sclient.CoreV1().Services(app.Namespace).Create(service)
	if err !=nil{
		fmt.Println(err.Error())
	}

	return result,err

}


func (this * KubernetesClient)CreateReplicationController(app K8sApp,pod *v1.Pod) (*v1.ReplicationController,error){
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
	result, err :=this.K8sclient.CoreV1().ReplicationControllers(app.Namespace).Create(rc)
	if err !=nil{
		fmt.Println(err.Error())
	}

	return result,err
}

func (this * KubernetesClient)createcontainer(name,image string ,cmd []string,targetport int32)(v1.Container){
	container:=v1.Container{
		Name:  name,
		Image: image,
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				ContainerPort: targetport,
				Protocol:      v1.ProtocolTCP,
			},
		},
	}

	if len(cmd)>0 && len(cmd[0])>0 {
		container.Command=cmd
	}

	return container
}


func (this * KubernetesClient)CreatePod(app K8sApp)(*v1.Pod,error){

	pod:=new(v1.Pod)
	pod.TypeMeta=metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"}
	pod.ObjectMeta=metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace, Labels: map[string]string{"name": app.Name}}
	pod.Spec=v1.PodSpec{
		RestartPolicy: v1.RestartPolicyAlways,
	}

	var container v1.Container
	if app.Types==TYPE_VNC{
		container=this.createcontainer(app.Name,app.Image,app.Cmd,5900)
	}else{
		container=this.createcontainer(app.Name,app.Image,app.Cmd,app.TargetPort)
	}
	if len(app.Env)>0{
		container.Env=app.Env
	}

	if app.Types==TYPE_RTIG{
		nfsServer:=beego.AppConfig.String("nfsserver")
		nfsPath:=beego.AppConfig.String("nfspath")

		volumes:=make([]v1.Volume,1)
		volume:=v1.Volume{
			Name:"nfs-storage",
		}
		volume.NFS=&v1.NFSVolumeSource{nfsServer, nfsPath, false}
		volumes[0]=volume
		pod.Spec.Volumes=volumes

		volumeMounts:=make([]v1.VolumeMount,1)
		volumeMount:=v1.VolumeMount{
			Name:"nfs-storage",
			MountPath:"/root/certi/fom_files",
			SubPath:app.Namespace+"/"+app.TaskName,
		}
		volumeMounts[0]=volumeMount
		container.VolumeMounts=volumeMounts
	}

	pod.Spec.Containers= []v1.Container{container}

	if app.Types==TYPE_VNC{
		containernovnc:=this.createcontainer(app.Name+"-novnc",MajorRegistry.GetIpPort()+"/novnc",nil,app.TargetPort)
		pod.Spec.Containers=append(pod.Spec.Containers,containernovnc)
	}

	result, err := this.K8sclient.CoreV1().Pods(app.Namespace).Create(pod)
	if err !=nil{
		fmt.Println(err.Error())
	}

	return result,err
}


func (this * KubernetesClient)RemoveReplicationController(namespace string,name string)error{
	err :=this.K8sclient.CoreV1().ReplicationControllers(namespace).Delete(name,&metav1.DeleteOptions{})
	if err !=nil{
		fmt.Println(err.Error())
	}
	return err
}

func (this * KubernetesClient)RemovePod(namespace string,name string)error{
	err :=this.K8sclient.CoreV1().Pods(namespace).Delete(name,&metav1.DeleteOptions{})
	if err !=nil{
		fmt.Println(err.Error())
	}
	return err
}

func (this * KubernetesClient)RemoveService(namespace string,name string)error{
	_,err:=this.K8sclient.CoreV1().Services(namespace).Get(name,metav1.GetOptions{})
	if err==nil {
		err =this.K8sclient.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			fmt.Println(err.Error())
		}
		return err
	}
	return nil

}

func (this * KubernetesClient)GetPod(namespace string,name string)(*v1.Pod,error){
	pod,err:=this.K8sclient.CoreV1().Pods(namespace).Get(name,metav1.GetOptions{})
	return pod,err

}

func (this * KubernetesClient)ShowLogs(namespace string,name string) (string,error){
	request:=this.K8sclient.CoreV1().Pods(namespace).GetLogs(name,&v1.PodLogOptions{})
	result:=request.Do()
	body,err:=result.Raw()
	if err !=nil{
		fmt.Println(err.Error())
	}
	return string(body),err
}