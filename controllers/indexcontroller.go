package controllers

import (
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/entity"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.TplName = "index.html"
	c.Data["k8s"]=entity.Newk8sui.GetUrl();
}


type SetController struct {
	beego.Controller
}

func (c *SetController) Get() {
	c.TplName = "setting.html"
	c.Data["registry"]=&entity.Newregistry
	c.Data["k8sui"]=&entity.Newk8sui
}