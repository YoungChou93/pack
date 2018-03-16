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
	"github.com/docker/docker/api/types/container"
)

type ListImageController struct {
	beego.Controller
}

func (c *ListImageController) Get() {
	c.TplName = "localimage.html"
}

func (c *ListImageController) List() {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		//错误日志
		beego.Error(err.Error())
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		//错误日志
		beego.Error(err.Error())
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
	defer cli.Close()
	if err != nil {
		//错误日志
		beego.Error(err.Error())
		panic(err)
	}

	_, err = cli.ImageRemove(context.Background(), imageid, types.ImageRemoveOptions{Force:true,PruneChildren:true})


	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	} else {

		result.Success = true
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *ListImageController) Run() {
	imagename := c.GetString("imagename")

	result := entity.Result{Success:true}

	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		//错误日志
		beego.Error(err.Error())
		panic(err)
	}

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: imagename,
	}, nil, nil, "")
	if err != nil {
		//错误日志
		beego.Error(err.Error())
		result.Success = false
		result.Reason = err.Error()
	}else {
		if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
			//错误日志
			beego.Error(err.Error())
			result.Success = false
			result.Reason = err.Error()
		}else {
			if err := cli.ContainerStop(context.Background(), resp.ID, nil); err != nil {
				//错误日志
				beego.Error(err.Error())
				result.Success = false
				result.Reason = err.Error()
			}else{
				if err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{}); err != nil {
					//错误日志
					beego.Error(err.Error())
					result.Success = false
					result.Reason = err.Error()
				}
			}
		}
	}
	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *ListImageController) ListContainer() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	c.Data["json"] = &containers
	c.ServeJSON()

}