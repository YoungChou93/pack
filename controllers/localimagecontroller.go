package controllers

import (
	"github.com/YoungChou93/pack/entity"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	myclient "github.com/YoungChou93/pack/client"
	"strconv"
	"time"
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

	imageList:=make([]myclient.LocalImageForUI,0)
	for _,image:=range images{
		if len(image.RepoTags)<1{
			continue
		}
		size:=strconv.FormatInt(image.Size/1000000,10)+"MB"
		time:=time.Unix(image.Created, 0).Format("2006-01-02 15:04:05")
		imageList=append(imageList,myclient.LocalImageForUI{Name:image.RepoTags[0],Id:image.ID[7:19],Size:size,Created:time})
	}
	c.Data["json"] = &imageList
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

