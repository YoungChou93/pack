package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) Get() {
	c.TplName = "packone.html"
}

func (c *UploadController) EncapsulationView() {
	c.TplName = "encapsulation.html"
}

func (c *UploadController) Encapsulation() {
	var err error

	imagename := c.GetString("imagename")
	version := c.GetString("version")
	software, header, _ := c.GetFile("software")
	commands := c.GetString("commands")
	baseimage := c.GetString("baseimage")

	//保存待封装软件
	packpath := beego.AppConfig.String("packpath")

	dirpath := packpath + imagename + version
	os.MkdirAll(dirpath, os.ModePerm)

	var delimiter string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		delimiter = "\\"
	} else {
		delimiter = "/"
	}

	if software != nil {
		f, err := os.OpenFile(dirpath+delimiter+header.Filename, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			fmt.Println("文件打开失败")
			c.Abort("fail")
		}

		_, err = io.Copy(f, software)
		if err != nil {
			fmt.Println("文件保存失败" + err.Error())
			c.Abort("fail")
		}

		f.Close()


		/*if strings.Contains(header.Filename, ".tar") {
			err := util.Untar(dirpath+delimiter+header.Filename, dirpath+delimiter)
			if err != nil {
				fmt.Println(err.Error())
				c.Abort("fail")
			}
		}*/

	}

	//生成Dockerfile
	dockerfile, err := os.OpenFile(dirpath+delimiter+"Dockerfile", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("文件打开失败")
		c.Abort("fail")
	}
	defer dockerfile.Close()

	io.WriteString(dockerfile, "FROM"+" "+baseimage)
	io.WriteString(dockerfile, "\n")
	io.WriteString(dockerfile, commands)

	//封装镜像
	tarpath := packpath + "tar/"
	os.MkdirAll(tarpath, os.ModePerm)

	tar := new(archivex.TarFile)
	tar.Create(packpath + "tar/" + imagename + version)
	tar.AddAll(dirpath, false)
	tar.Close()

	//删除临时文件
	defer func() {
		os.RemoveAll(dirpath)
		os.Remove(tarpath + imagename + version + ".tar")
	}()


	dockerBuildContext, err := os.Open(tarpath + imagename + version + ".tar")
	defer dockerBuildContext.Close()
	defaultHeaders := map[string]string{"Content-Type": "application/tar"}
	cli, _ := client.NewClient("unix:///var/run/docker.sock", "v1.27", nil, defaultHeaders)
	options := types.ImageBuildOptions{
		SuppressOutput: true,
		Remove:         true,
		ForceRemove:    true,
		Tags:           []string{imagename + ":" + version}}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	//defer buildResponse.Body.Close()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Printf("********* %s **********", buildResponse.OSType)
	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println(string(response))
	result:=string(response)
	if(strings.Contains(result,"error")){
		c.Data["content"] = "FAIL"
	}else{
		c.Data["content"] = "SUCCESS"
	}

	c.TplName = "result.html"

}

func (this *UploadController) Post() {
	var err error

	imagename := this.GetString("imagename")
	version := this.GetString("version")
	software, header, _ := this.GetFile("software")
	command := this.GetString("command")
	baseimage := this.GetString("baseimage")

	//保存待封装软件
	packpath := beego.AppConfig.String("packpath")

	dirpath := packpath + imagename + version
	os.MkdirAll(dirpath, os.ModePerm)

	var delimiter string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		delimiter = "\\"
	} else {
		delimiter = "/"
	}
	f, err := os.OpenFile(dirpath+delimiter+header.Filename, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer f.Close()

	_, err = io.Copy(f, software)
	if err != nil {
		fmt.Println("文件保存失败" + err.Error())
		return
	}

	//生成Dockerfile
	dockerfile, err := os.OpenFile(dirpath+delimiter+"Dockerfile", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer dockerfile.Close()

	io.WriteString(dockerfile, "FROM"+" "+baseimage)
	io.WriteString(dockerfile, "\n")
	io.WriteString(dockerfile, "ADD"+" "+header.Filename+" "+"/usr/local")
	io.WriteString(dockerfile, "\n")
	io.WriteString(dockerfile, "CMD"+" "+command)

	//封装镜像
	tarpath := packpath + "tar/"
	os.MkdirAll(tarpath, os.ModePerm)

	tar := new(archivex.TarFile)
	tar.Create(packpath + "tar/" + imagename + version)
	tar.AddAll(dirpath, false)
	tar.Close()
	dockerBuildContext, err := os.Open(tarpath + imagename + version + ".tar")
	defer dockerBuildContext.Close()
	defaultHeaders := map[string]string{"Content-Type": "application/tar"}
	cli, _ := client.NewClient("unix:///var/run/docker.sock", "v1.27", nil, defaultHeaders)
	options := types.ImageBuildOptions{
		SuppressOutput: true,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{imagename + ":" + version}}
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

	this.Ctx.Redirect(302, "/upload")

}
