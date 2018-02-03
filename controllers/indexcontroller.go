package controllers

import (
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/entity"
	"github.com/YoungChou93/pack/client"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.TplName = "index.html"
	c.Data["k8s"]=entity.Newk8sui.GetUrl()
}


type SetController struct {
	beego.Controller
}

func (c *SetController) Get() {
	c.TplName = "setting.html"
	c.Data["registry"]=&client.Newregistry
	c.Data["k8sui"]=&entity.Newk8sui
}

func (c *SetController) SetRegistry() {
	ipaddr := c.GetString("ipaddr")
	port := c.GetString("port")
	version := c.GetString("version")

	result:=entity.Result{}

	client.Newregistry.Ipaddr=ipaddr
	client.Newregistry.Port=port
	client.Newregistry.Version=version

	result.Success=true
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *SetController) SetK8sui() {
	ipaddr := c.GetString("ipaddr")
	port := c.GetString("port")
	route := c.GetString("route")

	result:=entity.Result{}

	entity.Newk8sui.Ipaddr=ipaddr
	entity.Newk8sui.Port=port
	entity.Newk8sui.Route=route

	result.Success=true
	c.Data["json"] = &result
	c.ServeJSON()
}