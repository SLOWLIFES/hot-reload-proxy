package main

import (
	"fmt"
	"github.com/SLOWLIFES/hot-reload-proxy/hot_reload_proxy"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var c = config{}
var proxy = hot_reload_proxy.New()

func registered(writer http.ResponseWriter, request *http.Request) {
	addr := request.URL.Query().Get("addr")
	token := request.URL.Query().Get("token")
	fmd5 := request.URL.Query().Get("md5")
	log.Println("addr:", addr)
	log.Println("token:", token)
	log.Println("md5:", fmd5, c.MainMd5)
	writer.WriteHeader(200)
	if token != c.Token {
		fmt.Errorf("error1")
		_, _ = writer.Write([]byte("error1"))
		return
	}

	addrSplit := strings.Split(addr, ":")
	if len(addrSplit) != 2 {
		log.Println("error2")
		_, _ = writer.Write([]byte("error2"))
		return
	}
	if addrSplit[0] == "" || addrSplit[0] == "0.0.0.0" {
		addrSplit[0] = "127.0.0.1"
	}
	addr = strings.Join(addrSplit, ":")

	if fmd5 != c.MainMd5 {
		log.Println("error4")
		_, _ = writer.Write([]byte("error4"))
		return
	}

	isOk := false
	for i := 0; i < 50; i++ {
		log.Printf("端口检查第%d次", i)
		res, err := HttpDo{
			Url: fmt.Sprintf("http://%s", addr),
		}.Get()
		if err == nil && string(res) != "" {
			isOk = true
			break
		}
		time.Sleep(time.Second)
	}
	if !isOk {
		log.Println("error3")
		_, _ = writer.Write([]byte("error3"))
		return
	}

	err := proxy.SwitchHost(addr)
	c.RemoteAddr = addr
	saveConfig()
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
	} else {
		spl := SupervisorProgramList()
		for e := range spl {
			if spl[e].Name != c.MainMd5 {
				SupervisorDeleteProgram(spl[e].Name)
			}
		}
		SupervisorProgramReload()

		log.Println("ok")
		_, _ = writer.Write([]byte("ok"))
	}
}

func main() {
	loadConfig()
	go func() {
		recover()
		for {
			main, _ := ioutil.ReadFile(c.Main)
			time.Sleep(4 * time.Second)
			main2, _ := ioutil.ReadFile(c.Main)
			if len(main) > 100 && md5Str(main) == md5Str(main2) && string(main) != "" {
				md5 := md5Str(main)
				if c.MainMd5 != md5 {
					newMain := filepath.Join("upload-files", md5)
					if runtime.GOOS == "windows" {
						newMain = newMain + ".exe"
					}
					newMain, _ = filepath.Abs(newMain)
					_ = ioutil.WriteFile(newMain, main, 0777)
					c.MainMd5 = md5
					saveConfig()

					err := SupervisordAddProgram(c.MainMd5, newMain+" "+c.Arg, filepath.Dir(c.Main))
					SupervisorProgramReload()
					if err != nil {
						log.Println("SupervisordAddProgram(c.MainMd5, newMain+c.Arg):", err.Error())
					} else {
						log.Println("SupervisorProgramStart(c.MainMd5):", c.MainMd5)

						time.Sleep(1 * time.Second)
						log.Println("SupervisorProgramStart(c.MainMd5):", SupervisorProgramStart(c.MainMd5))
					}
				}
			}
		}
	}()

	_ = proxy.SwitchHost(c.RemoteAddr)

	http.HandleFunc("/hot-reload-proxy/registered", registered)

	http.Handle("/", proxy)
	err := http.ListenAndServe(c.LocalAddr, nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
