package client

import (
    "github.com/astaxie/beego"
	"github.com/YoungChou93/pack/database"
	"strconv"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/docker/docker/client"
	"io"
	"os"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
)

var MajorRegistry Registry

//构造镜像仓库镜像格式
type RepositoriesInfo struct {
	Repositories []string
}

type Image struct {
	Name string
	Tags []string
}

type RegistryImage struct{
	Images []Image
}

type ImageForUI struct{
	Name string
	Tag string
	Registry string
}

type LocalImageForUI struct{
	Name string
	Size string
	Id string
	Created string
}


//镜像仓库地址
type Registry struct {
	Ipaddr  string
	Port    string
	Version string
}

func (this *Registry) GetUrl() string {
	return "https://" + this.Ipaddr + ":" + this.Port + "/" + this.Version
}

func (this *Registry) GetIpPort() string {
	return this.Ipaddr + ":" + this.Port
}

func (this *Registry) List() ([]ImageForUI,error){
	var resp *http.Response
	var err error
	var body []byte

	images:=make([]ImageForUI,0)

	resp, err = http.Get(this.GetUrl() + "/_catalog")
	if err != nil {
		fmt.Println(err)
		return images,err

	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return images,err
	}

	var repositoriesInfo RepositoriesInfo
	json.Unmarshal(body, &repositoriesInfo)

	imageList := make([]Image, 0)

	registryimage := RegistryImage{Images: imageList}

	for _, name := range repositoriesInfo.Repositories {
		resp, err = http.Get(this.GetUrl() + "/" + name + "/tags/list")

		if err != nil {
			fmt.Println(err)
			return images,err
		}

		body, err = ioutil.ReadAll(resp.Body)

		var image Image

		json.Unmarshal(body, &image)

		registryimage.Images = append(registryimage.Images, image)
	}


	for _,image:=range registryimage.Images{
		for _,tag:=range image.Tags{
			images=append(images,ImageForUI{image.Name,tag,this.GetIpPort()})
		}
	}

	return images,nil

}

func (this *Registry)Pull(imagename string)error{

	imagename = this.GetIpPort() + "/" + imagename

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	out, err := cli.ImagePull(ctx, imagename, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	return nil
}

func RegistrySetting(){
	MajorRegistry = Registry{beego.AppConfig.String("registryip"), beego.AppConfig.String("registryport"), beego.AppConfig.String("registryversion")}
}

func GetRegistry(registry database.Registry) Registry{

	return Registry{registry.Ip,strconv.Itoa(registry.Port),registry.Version}
}