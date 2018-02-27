package util

import (
	"io/ioutil"
)

type FileNode struct {
	Text  string     `json:"text"`
	Icon  string     `json:"icon"`
	Nodes []FileNode `json:"nodes"`
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
			node = FileNode{fi.Name(),"glyphicon glyphicon-folder-close", nodes}
			ListDir(dirPth+"/"+fi.Name(), &node)
		} else {
			node = FileNode{fi.Name(), "glyphicon glyphicon-file",nil}
		}
		filenode.Nodes = append(filenode.Nodes, node)
	}
	return nil
}
