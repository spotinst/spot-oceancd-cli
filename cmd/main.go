package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"

	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type Applyer interface {
	Apply(obj interface{}) error
}

func apply() {
	//api

}

func process(dir string, files []fs.FileInfo) {
	for _, f := range files {
		fmt.Println(f.Name())
		if !f.IsDir() {
			filePath := dir + "/" + f.Name()
			ext := filepath.Ext(filePath)
			fmt.Printf("type is %v\n", ext)
			if ext == ".json" {
				bytes, _ := ioutil.ReadFile(filePath)
				pod := core.Pod{}
				err := json.Unmarshal(bytes, &pod)
				if err != nil {
					continue
				}

				fmt.Printf("deployment %v\n", pod.Name)

			}
			if ext == ".yaml" {
				bytes, _ := ioutil.ReadFile(filePath)
				d := v1.Deployment{}
				err := yaml.Unmarshal(bytes, &d)
				if err != nil {
					continue
				}

				fmt.Printf("deployment %v\n", d.Name)

			}
		}
	}
}

func exec() {
	dir := "./yamls"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	process(dir, files)
}

func main() {
	err := NewRootCmd().Execute()
	if err != nil {
		panic(err)
	}
}
