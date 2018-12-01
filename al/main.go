package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lcl101/rcmd/core"
)

var (
	//Version 版本信息
	Version = "0.1.0"
	//Build 编译时间
	Build = "20181130"
)

var (
	version = flag.Bool("v", false, "版本信息")
	help    = flag.Bool("help", false, "帮助")
	config  = flag.String("c", "", "配置文件，默认al.conf")
	en      = flag.String("e", "", "加密密码")
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

	if *en != "" {
		s, err := core.Encrypt(*en)
		if err != nil {
			fmt.Println("en error: ", err)
		} else {
			fmt.Println(s)
		}
		os.Exit(0)
	}

	serverName := ""
	if len(os.Args) > 1 {
		serverName = os.Args[1]
	}

	core.Log.Category("main").Info("serverName: ", serverName)

	// core.Log.Category("main").Info("key=", core.StrKey)

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
	app.Init(serverName)
}

// 版本信息
func v() {
	fmt.Println("al version: " + Version + ", Build " + Build + "。")
	fmt.Println("本程序源码：https://github.com/lcl101/rcmd。")
}
