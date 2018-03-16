package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	myclient "github.com/YoungChou93/pack/client"
	"github.com/YoungChou93/pack/database"
	"github.com/YoungChou93/pack/entity"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"os"
	"strings"
)

type RegistryController struct {
	beego.Controller
}

func (c *RegistryController) Get() {
	c.TplName = "registry.html"
}

func (c *RegistryController) RegistryView() {
	c.TplName = "registrymanage.html"
}

func (c *RegistryController) List() {
	id, err := c.GetInt("id")

	var newRegistry myclient.Registry

	if err !=nil{

		newRegistry=myclient.MajorRegistry
	}else{
		registry := database.Registry{Id: id}
		database.Dao.Read(&registry)
		newRegistry = myclient.GetRegistry(registry)
	}

	images, err := newRegistry.List()

	if err != nil {
		fmt.Println(err)
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &images
	c.ServeJSON()
}

func (c *RegistryController) Push() {

	imagename := c.GetString("imagename")
	id, _ := c.GetInt("id")

	registry := database.Registry{Id: id}
	database.Dao.Read(&registry)

	newRegistry := myclient.GetRegistry(registry)

	result := entity.Result{Success: true}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	var newTag string
	if strings.Contains(imagename, newRegistry.GetIpPort()) {
		newTag = imagename
	} else {
		if strings.Contains(imagename, "/") {
			strs := strings.Split(imagename, "/")
			imagename = strs[len(strs)-1]

		}
		newTag = newRegistry.GetIpPort() + "/" + imagename

		cli.ImageTag(ctx, imagename, newTag)
	}

	defer func() {
		cli.ImageTag(ctx, newTag, imagename)
	}()

	auth := types.AuthConfig{
		Username: "docker",
		Password: "docker",
	}
	authBytes, _ := json.Marshal(auth)
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)

	out, err := cli.ImagePush(ctx, newTag, types.ImagePushOptions{RegistryAuth: authBase64})

	if err != nil {
		fmt.Println(err.Error())
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	io.Copy(os.Stdout, out)

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *RegistryController) Pull() {

	imagename := c.GetString("imagename")

	id, _ := c.GetInt("id")

	registry := database.Registry{Id: id}
	database.Dao.Read(&registry)

	newRegistry := myclient.GetRegistry(registry)

	result := entity.Result{Success: true}

	err := newRegistry.Pull(imagename)

	if err != nil {
		fmt.Println(err.Error())
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *RegistryController) AddRegistry() {
	name := c.GetString("name")
	ip := c.GetString("ip")
	port, _ := c.GetInt("port")
	version := c.GetString("version")

	result := entity.Result{Success: true}

	var registry database.Registry
	registry.Name = name
	registry.Ip = ip
	registry.Port = port
	registry.Version = version

	_, err := database.Dao.Insert(&registry)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *RegistryController) ListRegistry() {

	var registries []*database.Registry
	_, err := database.Dao.QueryTable("registry").All(&registries)
	if err != nil {
		fmt.Println(err)
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &registries
	c.ServeJSON()

}

func (c *RegistryController) MajorRegistry() {

	id, _ := c.GetInt("id")
	result := entity.Result{Success: true}
	registry := database.Registry{Id: id}
	err := database.Dao.Read(&registry)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	} else {
		_, err := database.Dao.QueryTable("registry").Filter("major", 1).Update(orm.Params{"major": 0})
		registry.Major = 1
		_, err = database.Dao.Update(&registry)
		if err != nil {
			result.Success = false
			result.Reason = err.Error()
			//错误日志
			beego.Error(err.Error())
		} else {
			myclient.MajorRegistry = myclient.GetRegistry(registry)
		}

	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *RegistryController) DeleteRegistry() {
	id, _ := c.GetInt("id")
	result := entity.Result{Success: true}
	user := database.Registry{Id: id}
	_, err := database.Dao.Delete(&user)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &result
	c.ServeJSON()
}
