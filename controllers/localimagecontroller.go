package controllers

import (
	"github.com/YoungChou93/pack/entity"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type ListImageController struct {
	beego.Controller
}

func (c *ListImageController) Get() {
	c.TplName = "localimage.html"
}

func (c *ListImageController) List() {
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
	result := entity.Result{}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	response, err := cli.ImageRemove(context.Background(), imageid, types.ImageRemoveOptions{Force:true,PruneChildren:true})

	println(response)

	if err != nil {
		result.Success = false
		result.Reason = err.Error()
	} else {

		result.Success = true
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

