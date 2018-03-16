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
	"github.com/YoungChou93/pack/entity"
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
	result := entity.Result{Success:true}


	imagename := c.GetString("imagename")
	version := c.GetString("version")
	software, header, _ := c.GetFile("software")
	commands := c.GetString("commands")
	baseimage := c.GetString("baseimage")

	fmt.Println(imagename)

	//保存待封装软件
	packpath := beego.AppConfig.String("packpath")

	dirpath := packpath +"/"+ imagename + version
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
			//错误日志
			beego.Error(err.Error())
		}

		_, err = io.Copy(f, software)
		if err != nil {
			fmt.Println("文件保存失败" + err.Error())
			//错误日志
			beego.Error(err.Error())
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
	tarpath := packpath +"/"+ "tar/"
	os.MkdirAll(tarpath, os.ModePerm)

	tar := new(archivex.TarFile)
	tar.Create(tarpath + imagename + version)
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
	cli, _ := client.NewClient("unix:///var/run/docker.sock", beego.AppConfig.String("dockerversion"), nil, defaultHeaders)
	defer cli.Close()
	options := types.ImageBuildOptions{
		SuppressOutput: true,
		Remove:         true,
		ForceRemove:    true,
		Tags:           []string{imagename + ":" + version}}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	defer buildResponse.Body.Close()
	if err != nil {
		fmt.Printf("%s", err.Error())
		//错误日志
		beego.Error(err.Error())
	}
	fmt.Printf("********* %s **********", buildResponse.OSType)
	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s", err.Error())
		//错误日志
		beego.Error(err.Error())
	}
	fmt.Println(string(response))
	results:=string(response)
	if(strings.Contains(results,"error")){
		result.Success=false
		result.Reason=results
	}
	c.Data["json"] = &result
	c.ServeJSON()

}

