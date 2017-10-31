package controllers

import (
	"github.com/astaxie/beego"
	"github.com/docker/docker/client"
	"fmt"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
)

type ListImageController struct {
	beego.Controller
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

	for _, image := range images {
		fmt.Println(image.RepoTags[0])
	}

	c.Data["json"] = &images
	c.ServeJSON()
}