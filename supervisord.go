package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ochinchina/go-ini"
	"os"
	"path/filepath"
)

var (
	filePath = []string{
		"./supervisord",
		"/usr/bin/supervisord",
		"/usr/local/bin/supervisord",
		"/bin/supervisord",
	}

	CfgPath = "/root/supervisor/supervisor.conf"
)

func SupervisorCfg() *ini.Ini {
	cfg := ini.NewIni()
	if pathExists(CfgPath) {
		cfg.LoadFile(CfgPath)
	} else {
		s := cfg.NewSection("inet_http_server")
		s.Add("port", "127.0.0.1:9001")
		cfg.AddSection(s)
		cfg.WriteToFile(CfgPath)
	}
	return cfg
}

func SupervisordAddProgram(id string, cmd string, directory string) error {
	c := SupervisorCfg()
	if c.HasSection(fmt.Sprintf("program:%s", id)) {
		return errors.New("已经存在")
	}
	s := c.NewSection(fmt.Sprintf("program:%s", id))
	s.Add("command", cmd)

	s.Add("directory", directory)
	s.Add("stdout_logfile", filepath.Join(directory, "stdout_logfile.log"))
	s.Add("stdout_logfile_maxbytes", "20m")
	s.Add("stderr_logfile", filepath.Join(directory, "stderr_logfile.log"))
	s.Add("stderr_logfile_maxbytes", "20m")
	c.AddSection(s)
	return c.WriteToFile(CfgPath)
}

func SupervisorSaveProgram(id string, cmd string, directory string, logOutPath string, logErrPath string) {
	SupervisorDeleteProgram(id)
	c := SupervisorCfg()
	section := fmt.Sprintf("program:%s", id)
	s := c.NewSection(section)
	s.Add("command", cmd)
	s.Add("directory", directory)
	s.Add("stdout_logfile", logOutPath)
	s.Add("stdout_logfile_maxbytes", "20m")
	s.Add("stderr_logfile", logErrPath)
	s.Add("stderr_logfile_maxbytes", "20m")
	c.AddSection(s)
	c.WriteToFile(CfgPath)
}

func SupervisorDeleteAllProgram() {
	os.Remove(CfgPath)
	c := SupervisorCfg()
	c.WriteToFile(CfgPath)
}

func SupervisorDeleteProgram(id string) {
	cfg := ini.NewIni()
	c := SupervisorCfg()
	section := fmt.Sprintf("program:%s", id)
	sl := c.Sections()
	for e := range sl {
		if sl[e].Name != section {
			cfg.AddSection(sl[e])
		}
	}
	cfg.WriteToFile(CfgPath)
}

type SupervisorProgram struct {
	Name          string `json:"name"`
	Group         string `json:"group"`
	Description   string `json:"description"`
	Start         int    `json:"start"`
	Stop          int    `json:"stop"`
	Now           int    `json:"now"`
	State         int    `json:"state"`
	Statename     string `json:"statename"`
	Spawnerr      string `json:"spawnerr"`
	Exitstatus    int    `json:"exitstatus"`
	Logfile       string `json:"logfile"`
	StdoutLogfile string `json:"stdout_logfile"`
	StderrLogfile string `json:"stderr_logfile"`
	Pid           int    `json:"pid"`
}

func SupervisorIsRun() bool {
	addr, _ := SupervisorCfg().GetValue("inet_http_server", "port")
	if addr != "" {
		_, err := HttpDo{
			Url: fmt.Sprintf("http://%s/program/list", addr),
		}.Get()
		if err != nil {
			return false
		}
		return true
	}

	return false
}

func SupervisorProgramList() []SupervisorProgram {
	var sp []SupervisorProgram
	addr, _ := SupervisorCfg().GetValue("inet_http_server", "port")
	if addr != "" {
		data, _ := HttpDo{
			Url: fmt.Sprintf("http://%s/program/list", addr),
		}.Get()
		json.Unmarshal(data, &sp)
	}

	return sp
}

func SupervisorProgramReload() {
	addr, _ := SupervisorCfg().GetValue("inet_http_server", "port")
	if addr != "" {
		_, _ = HttpDo{
			Url: fmt.Sprintf("http://%s/supervisor/reload", addr),
		}.Post()
	}
}

func SupervisorProgramStart(name string) string {
	addr, _ := SupervisorCfg().GetValue("inet_http_server", "port")
	if addr != "" {
		d, _ := HttpDo{
			Url: fmt.Sprintf("http://%s/program/start/%s", addr, name),
		}.Post()
		return string(d)
	}
	return ""
}

func SupervisorProgramStop(name string) string {
	addr, _ := SupervisorCfg().GetValue("inet_http_server", "port")
	if addr != "" {
		d, _ := HttpDo{
			Url: fmt.Sprintf("http://%s/program/stop/%s", addr, name),
		}.Post()
		return string(d)
	}
	return ""
}
