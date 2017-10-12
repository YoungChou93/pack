package controllers

import (
	"github.com/astaxie/beego"
	"os"
	"io"
	"fmt"
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
	f, err := os.OpenFile(dirpath+delimiter+header.Filename, os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer f.Close()

	_,err=io.Copy(f, software)
	if err != nil {
		fmt.Println("文件保存失败")
		return
	}


	//生成Dockerfile
	dockerfile, err := os.OpenFile(dirpath+delimiter+"Dockerfile", os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("文件打开失败")
		return
	}
	defer dockerfile.Close()

	io.WriteString(dockerfile,"FROM"+" "+baseimage)
	io.WriteString(dockerfile,"\r\n")
	io.WriteString(dockerfile,"ADD"+" "+header.Filename+" "+"C:\\software")
	io.WriteString(dockerfile,"\r\n")
	io.WriteString(dockerfile,"CMD"+" "+"["+"C:\\software"+header.Filename+"]")
	io.WriteString(dockerfile,"\r\n")

	this.Ctx.Redirect(302,"/upload")

}