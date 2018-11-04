package main

import (
	"flag"
	"fmt"

	"github.com/lcl101/rcmd/cmd"
)

var (
	host   = flag.String("h", "www.southlocal.cn", "主机地址")
	port   = flag.Int("P", 22, "端口")
	user   = flag.String("u", "", "用户名")
	passwd = flag.String("p", "", "密码")
	op     = flag.String("o", "cmd", "操作类型:执行命令:cmd,拷贝:cp")
	cmds   = flag.String("c", "cd /home, ls", "命令行,类似:cd /home/licl,ls")
	spath  = flag.String("s", "/home/aoki/work/wf/java/wowo/wowo-code/wowo-bweb/target/wowo-bweb-0.0.1-SNAPSHOT.tar.gz", "scp 源文件")
	dpath  = flag.String("d", "/home/software/wowo/wowo-b/test/wowo-bweb-0.0.1-SNAPSHOT.tar.gz", "scp 目标文件")
)

func main() {
	flag.Parse()

	conn := &cmd.SSHConn{User: *user, Password: *passwd, Port: *port, Host: *host}
	err := conn.Conn()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	if *op == "cmd" {
		session, err := conn.Session()
		defer session.Close()
		if err != nil {
			fmt.Println(err)
		}
		cmdList := cmd.SplitString(*cmds)
		c := cmd.NewCmd(cmdList, session)
		c.Run()
		str := fmt.Sprintf("%d---%s", c.RtnCode, c.RtnMsg)
		fmt.Println(str)
	} else if *op == "cp" {
		//scp
		err := cmd.Scp(conn.Client(), *spath, *dpath)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("不支持这个命令,-o仅支持cmd或者cp")
		flag.PrintDefaults()
	}
}
