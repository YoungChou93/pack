package util

import (
	"io/ioutil"
	"os"
)

type FileNode struct {
	Text  string     `json:"text"`
	Icon  string     `json:"icon"`
	Nodes []FileNode `json:"nodes"`
	Path  string     `json:"path"`
	Isfile bool    `json:"isfile"`
	Isroot bool   `json:"isroot"`

}


func ListDir(dirPth string, filenode *FileNode) error {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}
	for _, fi := range dir {
		var node FileNode
		if fi.IsDir() {
			nodes := make([]FileNode, 0)
			node = FileNode{fi.Name(),"glyphicon glyphicon-folder-close", nodes,dirPth+"/"+fi.Name(),false,false}
			ListDir(dirPth+"/"+fi.Name(), &node)
		} else {
			node = FileNode{fi.Name(), "glyphicon glyphicon-file",nil,dirPth+"/"+fi.Name(),true,false}
		}
		filenode.Nodes = append(filenode.Nodes, node)
	}
	return nil
}


func Chmod(dirPth string) error {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}
	err=os.Chmod(dirPth,0777)
	for _, fi := range dir {
		err=os.Chmod(dirPth+"/"+fi.Name(),0777)
		if fi.IsDir() {
			Chmod(dirPth+"/"+fi.Name())
		}
	}
	return err
}