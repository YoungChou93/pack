package entity

import (
	"container/list"
	"github.com/YoungChou93/pack/client"
	"errors"
)

type Application struct {
	taskMap map[string]*list.List
	findMap map[string]*list.Element
	client  client.KubernetesClient
}

func NewApplication(client client.KubernetesClient) Application {
	taskmap := make(map[string]*list.List)
	findmap := make(map[string]*list.Element)
	return Application{taskMap: taskmap, findMap: findmap,client:client}
}

func (this *Application) AddTask(namespace string, task Task) error {
	key := namespace + "#" + task.Name
	if _, ok := this.findMap[key]; ok {
		return errors.New("The task has already existed !")
	} else {
		var e *list.Element
		if tasks, ok := this.taskMap[namespace]; ok {
			e = tasks.PushBack(task)
		} else {
			tasks := list.New()
			e = tasks.PushBack(task)
			this.taskMap[namespace] = tasks
		}
		this.findMap[key] = e
	}
	return nil
}

func (this *Application) RemoveTask(namespace string, name string)(error) {
	key := namespace + "#" + name
	if e, ok := this.findMap[key]; ok {
		delete(this.findMap, key)
		if tasks, ok := this.taskMap[namespace]; ok {
			task := e.Value.(Task)
			tasks.Remove(e)
			for _, member := range task.Members {
				if member.Service != nil {
					err:=this.client.RemoveService(member.Namespace, member.Name)
					if err!=nil{
						return err
					}
				}
				err:=this.client.RemoveReplicationController(member.Namespace, member.Name)
				if err!=nil{
					return err
				}
				err=this.client.RemovePod(member.Namespace, member.Name)
				if err!=nil{
					return err
				}
			}
			return nil
		} else {
			return errors.New("error namespace")
		}
	} else {
		return errors.New("error name")
	}

}

func (this *Application) GetTasks(namespace string) ([]Task, error) {
	if tasks, ok := this.taskMap[namespace]; ok {
		return this.ListToTasks(tasks), nil
	} else {
		return nil, errors.New("Namespace does not exist")
	}
}

func (this *Application) ListToTasks(list *list.List) []Task {
	i := 0
	tasks := make([]Task, list.Len())
	for e := list.Front(); e != nil; e = e.Next() {
		tasks[i] = e.Value.(Task)
		i++
	}
	return tasks
}

func (this *Application) GetTask(namespace string, name string) Task {
	key := namespace + "#" + name
	if e, ok := this.findMap[key]; ok {
		task := e.Value.(Task)
		for index, member := range task.Members {
			//更新pod信息
			member.Pod ,_= this.client.GetPod(member.Namespace, member.Name)
			task.Members[index]=member
		}
		this.modifyTask(task)
		return task
	}
	return Task{}
}

func (this *Application) modifyTask(task Task) error {
	key := task.Namespace + "#" + task.Name
	if e, ok := this.findMap[key]; ok {
		if tasks, ok := this.taskMap[task.Namespace]; ok {
			tasks.Remove(e)
			e = tasks.PushBack(task)
		}
		this.findMap[key] = e
	} else {
		return errors.New("task does not exist")
	}
	return nil
}

func (this *Application) AddTaskMember(task Task,member TaskMember) error{
	pod,err:=this.client.CreatePod(member.GetK8sApp())
	if err!=nil{
		return err
	}
	rc,err:=this.client.CreateReplicationController(member.GetK8sApp(),pod)
	if err!=nil{
		return err
	}
	member.Pod=pod
	member.Rc=rc
	if member.Types!=3 {
		s, err := this.client.CreateService(member.GetK8sApp())
		if err!=nil{
			return err
		}
		member.Service = s
	}
	task.AddTaskMember(member)
	this.modifyTask(task)
	return nil
}

func (this *Application) RemoveTaskMember(namespace,name,membername string) error{

	task:=this.GetTask(namespace,name)
	member,err:=task.RemoveTaskMember(membername)
	if err !=nil{
		return err
	}
	if member.Service != nil {
		err=this.client.RemoveService(member.Namespace, member.Name)
		if err!=nil{
			return err
		}
	}
	err=this.client.RemoveReplicationController(member.Namespace, member.Name)
	if err!=nil{
		return err
	}
	err=this.client.RemovePod(member.Namespace, member.Name)
	if err!=nil{
		return err
	}
	this.modifyTask(task)
	return nil
}

func (this *Application) GetLog(namespace string,name string)(string,error){
	return this.client.ShowLogs(namespace,name)
}