package controllers

import (
	"github.com/YoungChou93/pack/database"
	"github.com/YoungChou93/pack/entity"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
	"strconv"
	"strings"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) UserView() {

	c.TplName = "user/usermanage.html"
}

func (c *UserController) ListUser() {

	var users []*database.User
	_, err := database.Dao.QueryTable("user").All(&users)
	if err != nil {
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &users
	c.ServeJSON()
}

func (c *UserController) AddUser() {
	username := c.GetString("username")
	result := entity.Result{Success: true}

	var user database.User
	user.Username = username
	user.Password = "123456"
	user.Status = 1

	_, err := database.Dao.Insert(&user)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())

	} else {
		userpath := beego.AppConfig.String("nfspath") + "/default/tool"
		os.MkdirAll(userpath, os.ModePerm)
	}

	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *UserController) EnableUser() {
	id, _ := c.GetInt("id")
	result := entity.Result{Success: true}
	user := database.User{Id: id}
	err := database.Dao.Read(&user)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}
	user.Status = 1
	_, err = database.Dao.Update(&user)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *UserController) DisableUser() {
	id, _ := c.GetInt("id")
	result := entity.Result{Success: true}
	user := database.User{Id: id}
	err := database.Dao.Read(&user)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}
	user.Status = 0
	_, err = database.Dao.Update(&user)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &result
	c.ServeJSON()

}

func (c *UserController) GetUserRight() {
	id, _ := c.GetInt("id")

	var rights []*database.Right
	_, err := database.Dao.QueryTable("right").All(&rights)
	if err != nil {
		//错误日志
		beego.Error(err.Error())
	}

	treenodes := make([]entity.TreeNode, len(rights))
	for index, right := range rights {
		var rumap database.Rightusermap
		err = database.Dao.QueryTable("rightusermap").Filter("uid", id).Filter("rid", right.Id).One(&rumap)
		state := entity.State{}
		if err == orm.ErrNoRows {
			state.Selected = false
		} else {
			state.Selected = true
		}
		node := entity.TreeNode{right.Id, right.Rightname, state}
		treenodes[index] = node
	}

	c.Data["json"] = &treenodes
	c.ServeJSON()

}

func (c *UserController) SetUserRight() {
	id, _ := c.GetInt("id")
	ids := c.GetString("ids")
	result := entity.Result{Success: true}

	_, err := database.Dao.QueryTable("rightusermap").Filter("uid", id).Delete()

	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	idarray := strings.Split(ids, ",")
	rumaps := make([]database.Rightusermap, len(idarray))

	for index, rid := range idarray {
		ridi, _ := strconv.Atoi(rid)
		rumap := database.Rightusermap{Rid: ridi, Uid: id}
		rumaps[index] = rumap
	}

	_, err = database.Dao.InsertMulti(len(idarray), rumaps)
	if err != nil {
		result.Success = false
		result.Reason = err.Error()
		//错误日志
		beego.Error(err.Error())
	}

	c.Data["json"] = &result
	c.ServeJSON()
}

func (c *UserController) LoginView() {

	c.TplName = "user/login.html"
}

func (c *UserController) Login() {
	username := c.GetString("username")
	password := c.GetString("password")

	var user database.User

	err := database.Dao.QueryTable("user").Filter("username", username).One(&user)

	if err == orm.ErrNoRows {
		c.TplName = "user/login.html"
		c.Data["errormsg"] = "用户名或密码错误"
	} else {
		if user.Password != password {
			c.TplName = "user/login.html"
			c.Data["errormsg"] = "用户名或密码错误"
		} else if user.Status==0{
			c.TplName = "user/login.html"
			c.Data["errormsg"] = "该用户已禁用"

		}else{
			c.SetSession("user", user)
			c.Redirect("/", 302)
		}

	}
}

func (c *UserController) Logout() {
	c.DelSession("user")
	c.Redirect("/login", 302)
}

func (c *UserController) ModifyPasswordView() {
	c.TplName = "user/password.html"
}

func (c *UserController) ModifyPassword() {

	user:=c.GetSession("user").(database.User)
	result := entity.Result{Success: true}
	oldpassword := c.GetString("oldpassword")
	password := c.GetString("password")
	confirmpassword := c.GetString("confirmpassword")

	if password!=confirmpassword{
		result.Success = false
		result.Reason = "两次密码输入不一致"
	}else if user.Password !=oldpassword{
		result.Success = false
		result.Reason = "原密码错误"
	}else{
		user.Password=password
		_,err:=database.Dao.Update(&user)
		if err!=nil{
			result.Success = false
			result.Reason = "原密码错误"
		}else{
			c.SetSession("user",user)
		}

	}

	c.Data["json"] = &result
	c.ServeJSON()
}