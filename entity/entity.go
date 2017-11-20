package entity

//镜像仓库地址
type Registry struct{
	Ipaddr string
	Port string
	Version string
}


func (this *Registry)GetUrl()string{
	return "http://"+this.Ipaddr+":"+this.Port+"/"+this.Version
}

func (this *Registry)GetIpPort()string{
	return this.Ipaddr+":"+this.Port
}

//k8sui地址
type K8sui struct{
	Ipaddr string
	Port string
	Route string
}


func (this *K8sui)GetUrl()string{
	return "http://"+this.Ipaddr+":"+this.Port+"/"+this.Route
}

var Newregistry Registry
var Newk8sui K8sui

func Setting()  {
	Newregistry=Registry{"192.168.182.151","5000","v2"}
	Newk8sui=K8sui{"192.168.182.150","8080","ui"}
}

//构造镜像仓库镜像格式
type RepositoriesInfo struct{
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
	Reason string
}