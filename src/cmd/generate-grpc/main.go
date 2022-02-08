package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)

func main() {

	paths := handle("./")
	group := make(map[string][]string)
	for _, path := range paths {
		group[filepath.Dir(path)] = append(group[filepath.Dir(path)], path)
	}
	for _, file := range group {
		command := append([]string{
			//"--gofast_out=.", // gofast 导致重复生成多个文件
			"--go_out=.", "--go_opt=paths=source_relative",
			"--go-grpc_out=.", "--go-grpc_opt=paths=source_relative"},
			file...)
		cmd := exec.Command("protoc", command...)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			log.Println("generate grpc proto failed!")
			log.Println(out.String())
		} else {
			log.Println("generate grpc proto successfully!")
			for _, path := range file {
				log.Println("proto file:", path)
			}
		}
	}

}

func handle(dir string) (filePaths []string) {
	var fileInfo, err = ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, info := range fileInfo {
		path := filepath.Join(dir, info.Name())
		if info.IsDir() {
			filePaths = append(filePaths, handle(path)...)
		} else if filepath.Ext(path) == ".proto" {
			filePaths = append(filePaths, path)
		}
	}
	return filePaths
}
