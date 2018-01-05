package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/YoungChou93/pack/entity"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"fmt"
)

type RegistryController struct {
	beego.Controller
}

func (c *RegistryController) Get() {

	c.TplName = "registry.html"
	c.Data["url"] = entity.Newregistry.GetUrl()

}

func (c *RegistryController) List() {
	var resp *http.Response
	var err error
	var body []byte

	resp, err = http.Get(entity.Newregistry.GetUrl() + "/_catalog")
	if err != nil {
		c.Data["error"] = err.Error()
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data["error"] = err.Error()
		return
	}

	var repositoriesInfo entity.RepositoriesInfo
	json.Unmarshal(body, &repositoriesInfo)

	imageList := make([]entity.Image, 0)

	registryimage := entity.RegistryImage{Images: imageList}

	for _, name := range repositoriesInfo.Repositories {
		resp, err = http.Get(entity.Newregistry.GetUrl() + "/" + name + "/tags/list")

		body, err = ioutil.ReadAll(resp.Body)

		var image entity.Image

		json.Unmarshal(body, &image)

		registryimage.Images = append(registryimage.Images, image)
	}

	c.Data["json"] = &registryimage
	c.ServeJSON()
}

func (c *RegistryController) Push() {

	imagename := c.GetString("imagename")

	result := entity.Result{}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	var newTag string
	if strings.Contains(imagename, entity.Newregistry.GetIpPort()) {
		newTag = imagename
	} else {
		newTag = entity.Newregistry.GetIpPort() + "/" + imagename
	}

	cli.ImageTag(ctx, imagename, newTag)

	defer func() {
		cli.ImageTag(ctx, newTag,imagename)
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
	} else {
		io.Copy(os.Stdout, out)
		result.Success = true
	}
	c.Data["json"] = &result
	c.ServeJSON()

}


func (c *RegistryController) Pull() {

	imagename := c.GetString("imagename")

	imagename = entity.Newregistry.GetIpPort() + "/" + imagename

	result := entity.Result{}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imagename, types.ImagePullOptions{})

	defer out.Close()

	io.Copy(os.Stdout, out)
	if err != nil {
		fmt.Println(err.Error())
		result.Success=false
		result.Reason=err.Error()
	}else{
		result.Success=true

	}

	c.Data["json"] = &result
	c.ServeJSON()
}