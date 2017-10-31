package controllers

import (
	"github.com/astaxie/beego"
	"os"
	"io"
	"fmt"
	"github.com/docker/docker/client"
	"io/ioutil"
	"github.com/jhoonb/archivex"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) Get() {
	c.TplName = "packone.html"
}

func (this *UploadController) Post() {
	var err error

	imagename := this.GetString("imagename")
	version := this.GetString("version")
	software,header,_ := this.GetFile("software")
	command := this.GetString("command")
	baseimage := this.GetString("baseimage")


	//保存待封装软件
	packpath:=beego.AppConfig.String("packpath")

	dirpath:=packpath+imagename+version
        os.MkdirAll(dirpath,os.ModePerm)


	var delimiter string
	if os.IsPathSeparator('\\') {  //前边的判断是否是系统的分隔符
		delimiter = "\\"
	} else {
		delimiter = "/"
	}
	f, err := os.OpenFile(dirpath+delimiter+header.Filename, os.O_CREATE | os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer f.Close()

	_,err=io.Copy(f, software)
	if err != nil {
		fmt.Println("文件保存失败"+err.Error())
		return
	}


	//生成Dockerfile
	dockerfile, err := os.OpenFile(dirpath+delimiter+"Dockerfile", os.O_CREATE | os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer dockerfile.Close()

	io.WriteString(dockerfile,"FROM"+" "+baseimage)
	io.WriteString(dockerfile,"\n")
	io.WriteString(dockerfile,"ADD"+" "+header.Filename+" "+"/usr/local")
	io.WriteString(dockerfile,"\n")
	io.WriteString(dockerfile,"CMD"+" "+command)

	//封装镜像
	tarpath :=packpath+"tar/"
	os.MkdirAll(tarpath,os.ModePerm)

	tar := new(archivex.TarFile)
	tar.Create(packpath+"tar/"+imagename+version)
	tar.AddAll(dirpath, false)
	tar.Close()
	dockerBuildContext, err := os.Open(tarpath+imagename+version+".tar")
	defer dockerBuildContext.Close()
	defaultHeaders := map[string]string{"Content-Type": "application/tar"}
	cli, _ := client.NewClient("unix:///var/run/docker.sock", "v1.27", nil, defaultHeaders)
	options := types.ImageBuildOptions{
		SuppressOutput: true,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{imagename+":"+version}}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	//defer buildResponse.Body.Close()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	//time.Sleep(5000 * time.Millisecond)
	fmt.Printf("********* %s **********", buildResponse.OSType)
	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println(string(response))

	this.Ctx.Redirect(302,"/upload")

}