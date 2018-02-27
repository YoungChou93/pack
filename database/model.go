package database

import "github.com/astaxie/beego/orm"

type User struct {
	Id int
	Username string
	Password string
	Status int
}

type Registry struct {
	Id int
	Name string
	Ip string
	Port int
	Version string
	Major int
}

type Right struct {
	Id int
	Rightname string
	Righturl string
	Icon string
}

type Rightusermap struct{
	Id int
	Rid int
	Uid int
}

func init(){
	orm.RegisterModel(new(User),new(Registry),new(Right),new(Rightusermap))
}