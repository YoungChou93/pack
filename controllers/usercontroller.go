package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/YoungChou93/pack/util"
	"github.com/YoungChou93/pack/database"
	"github.com/YoungChou93/pack/entity"
	"github.com/astaxie/beego/orm"
	"strings"
	"strconv"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) FileView() {
	c.TplName = "user/filemanage.html"
}

func (c *UserController) GetPackPath() {
	nodes:=make([]util.FileNode,0)
	node:=util.FileNode{"PackPath","glyphicon glyphicon-folder-close",nodes}
	err:=util.ListDir(beego.AppConfig.String("packpath"),&node)
	if err!=nil{
		fmt.Println(err)
	}

	result:=make([]util.FileNode,1)
	result[0]=node

	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *UserController) UserView() {
	c.TplName = "user/usermanage.html"
}

func (c *UserController) ListUser() {

	var users []*database.User
	_, err := database.Dao.QueryTable("user").All(&users)
	if err!=nil {
		fmt.Println(err)
	}

	c.Data["json"] = &users
	c.ServeJSON()
}



func (c *UserController) AddUser() {
	username := c.GetString("username")
	result := entity.Result{Success:true}

	var user database.User
	user.Username=username
	user.Password="123456"
	user.Status=1

	_,err:=database.Dao.Insert(&user)
	if err!=nil{
		result.Success=false
		result.Reason=err.Error()
	}

	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *UserController) EnableUser() {
	id,_ := c.GetInt("id")
	result := entity.Result{Success:true}
	user := database.User{Id: id}
	err := database.Dao.Read(&user)
	if err!=nil{
		result.Success=false
		result.Reason=err.Error()
	}
	user.Status = 1
	_,err=database.Dao.Update(&user)
	if err!=nil{
		result.Success=false
		result.Reason=err.Error()
	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *UserController) DisableUser() {
	id,_ := c.GetInt("id")
	result := entity.Result{Success:true}
	user := database.User{Id: id}
	err := database.Dao.Read(&user)
	if err!=nil{
		result.Success=false
		result.Reason=err.Error()
	}
	user.Status = 0
	_,err=database.Dao.Update(&user)
	if err!=nil{
		result.Success=false
		result.Reason=err.Error()
	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *UserController) GetUserRight() {
	id,_ := c.GetInt("id")

	var rights []*database.Right
	_, err := database.Dao.QueryTable("right").All(&rights)
	if err!=nil{
		fmt.Println(err)
	}


	treenodes:=make([]entity.TreeNode,len(rights))
	for index,right:=range rights{
		var rumap database.Rightusermap
		err=database.Dao.QueryTable("rightusermap").Filter("uid", id).Filter("rid",right.Id).One(&rumap)
		state:=entity.State{}
		if err == orm.ErrNoRows {
			state.Selected=false
		}else{
			state.Selected=true
		}
		node:=entity.TreeNode{right.Id,right.Rightname,state}
		treenodes[index]=node
	}

	c.Data["json"] = &treenodes
	c.ServeJSON()

}

func (c *UserController) SetUserRight() {
	id,_ := c.GetInt("id")
	ids:= c.GetString("ids")
	result := entity.Result{Success:true}


	_,err:=database.Dao.QueryTable("rightusermap").Filter("uid",id).Delete()

	if err!=nil{
		result.Success=false
		result.Reason=err.Error()
	}

	idarray:=strings.Split(ids,",")
	rumaps:=make([]database.Rightusermap,len(idarray))

	for index,rid:=range idarray{
		ridi,_:=strconv.Atoi(rid)
		rumap:=database.Rightusermap{Rid:ridi,Uid:id}
		rumaps[index]=rumap
	}

	_,err=database.Dao.InsertMulti(len(idarray),rumaps)
	if err!=nil{
		result.Success=false
		result.Reason=err.Error()
	}

	c.Data["json"] = &result
	c.ServeJSON()
}