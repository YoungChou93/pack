package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/YoungChou93/pack/entity"
	"os"
	"io"
	"github.com/YoungChou93/pack/util"
	"github.com/YoungChou93/pack/database"
	"os/exec"
)

type FileController struct {
	beego.Controller
}

func (c *FileController) FileView() {
	c.TplName = "user/filemanage.html"
}

func (c *FileController) GetPath() {
	user:=c.GetSession("user").(database.User)
	dirpath:=beego.AppConfig.String("nfspath") +user.Username+"/tool"
	os.MkdirAll(dirpath, os.ModePerm)

	nodes := make([]util.FileNode, 0)
	node := util.FileNode{"个人目录", "glyphicon glyphicon-folder-close", nodes, dirpath, false, true}
	err := util.ListDir(beego.AppConfig.String("nfspath")+"/"+"default/tool", &node)
	if err != nil {
		fmt.Println(err)
		//错误日志
		beego.Error(err.Error())
	}

	result := make([]util.FileNode, 1)
	result[0] = node

	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *FileController) UploadFile() {
	path := c.GetString("path")
	file, header, _ := c.GetFile("file")
	result := entity.Result{Success: true}
	if file != nil {
		f, err := os.OpenFile(path+"/"+header.Filename, os.O_CREATE|os.O_RDWR, 0777)
		defer f.Close()
		if err != nil {
			fmt.Println("文件打开失败"+ err.Error())
			result.Success = false
			result.Reason = err.Error()
			//错误日志
			beego.Error(err.Error())
		}else {
			_, err = io.Copy(f, file)
			if err != nil {
				fmt.Println("文件保存失败" + err.Error())
				result.Success = false
				result.Reason = err.Error()
				//错误日志
				beego.Error(err.Error())
			}
		}
	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *FileController) DeleteFile() {
	path := c.GetString("path")
	result := entity.Result{Success: true}
	cmd:=exec.Command("chmod","-R","0777",path)
	cmd.Run()
	err := os.RemoveAll(path)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *FileController) Download() {
	path := c.GetString("path")
	result := entity.Result{Success: true}
	file, err := os.Stat(path)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	var filename string

	if file.IsDir() {
		util.CreateTar(path)
		filename = path + ".tar"
		defer os.Remove(filename)
	} else {
		filename = path
	}

	c.Ctx.Output.Download(path + ".tar")

}