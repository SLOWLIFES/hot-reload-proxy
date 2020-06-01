package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
)

func jsonByFile(file string, v interface{}) {
	data, _ := ioutil.ReadFile(file)
	json.Unmarshal(data, v)
}
func jsonToFile(file string, v interface{}) bool {
	data, err := json.Marshal(v)
	if err != nil {
		return false
	}
	return ioutil.WriteFile(file, data, 0644) == nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func createUUID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}

func loadConfig() {
	if !pathExists("upload-files") {
		os.Mkdir("upload-files", 0777)
	}
	if !pathExists("hot-reload-proxy.json") {
		c.Token = createUUID()
		saveConfig()
	}
	jsonByFile("hot-reload-proxy.json", &c)
}

func saveConfig() {
	jsonToFile("hot-reload-proxy.json", &c)
}

type config struct {
	Token      string `json:"token"`
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
	Main       string `json:"main"`
	MainMd5    string `json:"main_md5"`
	Arg        string `json:"arg"`
}

func md5Str(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
