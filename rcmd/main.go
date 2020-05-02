package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lcl101/rcmd/core"
)

var (
	Version = "0.2.0"
	Build   = "20181130"
)

var (
	version = flag.Bool("v", false, "版本信息")
	help    = flag.Bool("help", false, "帮助")
	host    = flag.String("h", "www.southlocal.cn", "主机地址")
	port    = flag.Int("P", 22, "端口")
	user    = flag.String("u", "", "用户名")
	passwd  = flag.String("p", "", "密码")
	key     = flag.String("k", "", "ssh免密的路径key")
	op      = flag.String("o", "cmd", "操作类型:执行命令:cmd,拷贝:cp")
	cmds    = flag.String("c", "cd /home, ls", "命令行,类似:cd /home/licl,ls")
	spath   = flag.String("s", "/home/aoki/work/wf/java/wowo/wowo-code/wowo-bweb/target/wowo-bweb-0.0.1-SNAPSHOT.tar.gz", "scp 源文件")
	dpath   = flag.String("d", "/home/software/wowo/wowo-b/test/wowo-bweb-0.0.1-SNAPSHOT.tar.gz", "scp 目标文件")
	en      = flag.String("en", "", "加密密码")
	// de     = flag.String("de", "", "解码")
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
	// 解码函数
	// if *de != "" {
	// 	s, err := core.Decrypt(*de)
	// 	if err != nil {
	// 		fmt.Println("en error: ", err)
	// 	} else {
	// 		fmt.Println(s)
	// 	}
	// 	os.Exit(0)
	// }

	pw, err := core.Decrypt(*passwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn := &core.SSHConn{User: *user, Password: pw, Port: *port, Host: *host}
	err = conn.Conn()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *op == "cmd" {
		session, err := conn.Session()
		defer session.Close()
		if err != nil {
			fmt.Println(err)
		}
		cmdList := core.SplitString(*cmds)
		c := core.NewCmd(cmdList, session)
		c.Run()
		str := fmt.Sprintf("%d---%s", c.RtnCode, c.RtnMsg)
		fmt.Println(str)
	} else if *op == "cp" {
		//scp
		// err := core.Scp(conn.Client(), *spath, *dpath)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		session, err := conn.Client().NewSession()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer session.Close()
		sc, err := core.NewSessionClient(session)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		rp, rf := filepath.Split(*spath)
		lp, lf := filepath.Split(*dpath)
		err = sc.Send(rp, rf, lp, lf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("不支持这个命令,-o仅支持cmd或者cp")
		flag.PrintDefaults()
	}
}

// 版本信息
func v() {
	fmt.Println("rcmd version: " + Version + ", Build " + Build + "。")
	fmt.Println("本程序源码：https://github.com/lcl101/rcmd。")
}
