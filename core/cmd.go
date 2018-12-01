package core

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

//Cmd 执行命令
type Cmd struct {
	session *ssh.Session
	cmdlist []string
	RtnCode int
	RtnMsg  string
}

//NewCmd 新生成一个cmd
func NewCmd(cmdlist []string, session *ssh.Session) *Cmd {
	return &Cmd{cmdlist: cmdlist, session: session}
}

//Run 执行命令
func (c *Cmd) Run() {
	cmdlist := append(c.cmdlist, "exit")
	newcmd := strings.Join(cmdlist, "&&")

	var outbt, errbt bytes.Buffer
	c.session.Stdout = &outbt

	c.session.Stderr = &errbt
	err := c.session.Run(newcmd)
	if err != nil {
		c.RtnCode = 10
		c.RtnMsg = fmt.Sprintf("%s", err.Error())
		return
	}

	if errbt.String() != "" {
		c.RtnCode = 11
		c.RtnMsg = errbt.String()
	}

	c.RtnCode = 0
	c.RtnMsg = outbt.String()
}
