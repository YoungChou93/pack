package controllers

import (
	"github.com/astaxie/beego"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
)

type ListImageController struct {
	beego.Controller
}

func (c *ListImageController) Get(){
	c.TplName = "localimage.html"
}



func (c *ListImageController) Post() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	c.Data["json"] = &images
	c.ServeJSON()
}


func (c *ListImageController) Remove() {

	imageid := c.GetString("imageid")

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	response, err := cli.ImageRemove(context.Background(),imageid,types.ImageRemoveOptions{})
	if err != nil {
		panic(err)
	}

	c.Data["json"] = &response
	c.ServeJSON()
}