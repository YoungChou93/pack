package client

import (
"github.com/astaxie/beego"
)

var Newregistry Registry

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

func RegistrySetting(){
	Newregistry = Registry{beego.AppConfig.String("registryip"), beego.AppConfig.String("registryport"), beego.AppConfig.String("registryversion")}
}