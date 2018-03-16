package controllers

import (
	"github.com/astaxie/beego"
	"github.com/YoungChou93/pack/entity"
	"github.com/YoungChou93/pack/client"
	"github.com/YoungChou93/pack/database"
)


type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.TplName = "index.html"
	user:=c.GetSession("user").(database.User)

	var rights []*database.Right

	database.Dao.Raw("select * from `right` where id  in (select rid from rightusermap where uid = ?)",user.Id).QueryRows(&rights)

	c.Data["rights"]=rights
	c.Data["user"] = "欢迎您！" + user.Username

}


type SetController struct {
	beego.Controller
}

func (c *SetController) Get() {
	c.TplName = "setting.html"
	c.Data["registry"]=&client.MajorRegistry
	c.Data["k8sui"]=&entity.Newk8sui
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