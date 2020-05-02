package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lcl101/rcmd/core"
)

var (
	//Version 版本信息
	Version = "0.2.1"
	//Build 编译时间
	Build = "20190707"
)

var (
	version = flag.Bool("v", false, "版本信息")
	help    = flag.Bool("help", false, "帮助")
	config  = flag.String("c", "", "配置文件，默认al.conf")
	en      = flag.String("e", "", "加密密码")
	de      = flag.String("d", "", "licl")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		v()
		os.Exit(0)
	}

	defer func() {
		if err := recover(); err != nil {
			core.Log.Category("main").Error("recover", err)
		}
	}()
	conf := ""
	if *config == "" {
		// tmp, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		tmp, _ := core.GetExecPath()
		conf = tmp + "al.conf"
	} else {
		conf, _ = core.ParsePath(*config)
	}
	//only test
	// conf = "/home/aoki/work/bin/al.conf"

	core.Log.Category("main").Info("config path=", conf)

	_, err := os.Stat(conf)
	if err != nil {
		if os.IsNotExist(err) {
			core.Printer.Errorln("config file", conf+" not exists")
			core.Log.Category("main").Error("config file not exists")
		} else {
			core.Printer.Errorln("unknown error", err)
			core.Log.Category("main").Error("unknown error", err)
		}
		os.Exit(0)
	}
	app := core.App{
		ConfigPath: conf,
	}
	// only test
	// path0 := "./licl.txt"
	// path1 := "licl1:/root/licl.txt"

	path0 := os.Args[1]
	path1 := os.Args[2]

	path0 = strings.Trim(path0, " ")
	path1 = strings.Trim(path1, " ")

	//download = 1, upload = 2
	method := 0
	serverName := ""
	i := strings.Index(path0, ":")
	if i > 0 {
		serverName = path0[0:i]
		path0 = path0[i+1:]
		method = method | 1
	}

	i = strings.Index(path1, ":")
	if i > 0 {
		serverName = path1[0:i]
		path1 = path1[i+1:]
		method = method | 2
	}
	if method == 3 {
		//说明path0，path1都含服务器标签，
		core.Log.Category("main").Error("cp地址异常，都含有服务器标签")
		os.Exit(1)
	}

	if method == 0 {
		//说明path0，path1都不含服务器标签，
		core.Log.Category("main").Error("cp地址异常，都不含有服务器标签")
		os.Exit(1)
	}

	server := app.GetServer(serverName)
	passwd, err := core.Decrypt(server.Password)
	if err != nil {
		core.Log.Category("main").Error("解析密码错误!", err)
		os.Exit(1)
		// fmt.Println("de error: ", err)
		//处理异常
	}
	var rp, rf, lp, lf string
	//处理本地相对地址 下载
	if method == 1 {
		// 判读本地路径是否带文件名
		isDir := strings.HasSuffix(path1, "/")
		if strings.Index(path1, "./") >= 0 {
			path1 = filepath.Join(core.GetCurrentDirectory(), path1)
		}
		if !strings.HasPrefix(path1, "/") {
			path1 = filepath.Join(core.GetCurrentDirectory(), path1)
		}
		rp, rf = filepath.Split(path0)
		if isDir {
			lp = path1
			lf = rf
		} else {
			lp, lf = filepath.Split(path1)
		}
	}

	// 处理上传
	if method == 2 {
		isDir := strings.HasSuffix(path1, "/")
		if strings.Index(path0, "./") >= 0 {
			path0 = filepath.Join(core.GetCurrentDirectory(), path0)
		}
		if !strings.HasPrefix(path0, "/") {
			path0 = filepath.Join(core.GetCurrentDirectory(), path0)
		}
		lp, lf = filepath.Split(path0)
		if isDir {
			rp = path1
			rf = lf
		} else {
			rp, rf = filepath.Split(path1)
		}
	}
	fmt.Println(path0)
	fmt.Println(path1)

	conn := &core.SSHConn{User: server.User, Password: passwd, Port: server.Port, Host: server.IP}
	err = conn.Conn()
	defer conn.Close()
	if err != nil {
		core.Log.Category("main").Error("conn error!", err)
		return
	}
	session, err := conn.Client().NewSession()
	defer session.Close()
	if err != nil {
		core.Log.Category("main").Error("down file error!", err)
		return
	}
	sc, err := core.NewSessionClient(session)
	if err != nil {
		core.Log.Category("main").Error("down file error!", err)
		return
	}
	if method == 1 {
		err = sc.Receive(rp, rf, lp, lf)
		if err != nil {
			core.Log.Category("main").Error("down file error!", err)
			return
		}
	} else if method == 2 {
		err = sc.Send(rp, rf, lp, lf)
		if err != nil {
			core.Log.Category("main").Error("send file error!", err)
			return
		}
	} else {
		core.Log.Category("main").Error("scp err!", err)
	}

}

// 版本信息
func v() {
	fmt.Println("al version: " + Version + ", Build " + Build + "。")
	fmt.Println("本程序源码：https://github.com/lcl101/rcmd。")
}
