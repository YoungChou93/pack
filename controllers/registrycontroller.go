package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
	"github.com/YoungChou93/pack/entity"
	"encoding/json"
	"github.com/docker/docker/client"
	"io"
	"os"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"strings"
	"encoding/base64"
)

type RegistryController struct {
	beego.Controller
}

func (c *RegistryController) Get(){

	c.TplName = "registry.html"
	c.Data["url"]=entity.Newregistry.GetUrl()

}


func (c *RegistryController) Post() {
        var resp *http.Response
        var err error
        var body []byte

		resp, err = http.Get(entity.Newregistry.GetUrl()+"/_catalog")
		if err != nil {
			c.Data["error"]=err.Error()
			return
		}

		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Data["error"]=err.Error()
			return
		}

		var repositoriesInfo entity.RepositoriesInfo
		json.Unmarshal(body,&repositoriesInfo)



	    imageList:=make([] entity.Image,0)

		registryimage:=entity.ResgitryImage{Images:imageList}

		for _,name:=range repositoriesInfo.Repositories {
			resp, err = http.Get(entity.Newregistry.GetUrl()+"/"+name+"/tags/list")

			body, err = ioutil.ReadAll(resp.Body)

			var image entity.Image

			json.Unmarshal(body,&image)

			registryimage.Images=append(registryimage.Images,image)
		}

	    c.Data["json"] = &registryimage
	    c.ServeJSON()
}



func (c *RegistryController) Push(){

	imagename := c.GetString("imagename")

	result:=entity.Result{}

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	var newTag string
	if strings.Contains(imagename,entity.Newregistry.GetIpPort()) {
		newTag=imagename
	}else{
		newTag = entity.Newregistry.GetIpPort() + "/" + imagename
	}

	cli.ImageTag(ctx,imagename,newTag)

	auth := types.AuthConfig{
		Username: "docker",
		Password: "docker",
	}
	authBytes, _ := json.Marshal(auth)
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)

	out, err := cli.ImagePush(ctx,newTag,types.ImagePushOptions{RegistryAuth:authBase64})

	if err != nil {
		result.Success=false
		result.Reason=err.Error()
		c.Data["error"]=&result
		c.ServeJSON()
		return
	}

	io.Copy(os.Stdout, out)

	result.Success=true
	c.Data["json"]=&result
	c.ServeJSON()

}