package controllers

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Errorfail() {
	c.Data["content"] = "操作失败"
	c.TplName = "result.html"
}