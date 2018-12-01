package core

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	"golang.org/x/crypto/ssh"
)

//Scp 实现scp功能
func Scp(client *ssh.Client, spath, dpath string) error {
	if dpath == "" || spath == "" {
		flag.PrintDefaults()
		return errors.New("dpath,spath为空")
	}
	file, err := os.Open(spath)
	if err != nil {
		return err
	}
	info, _ := file.Stat()
	defer file.Close()

	filename := filepath.Base(dpath)
	dirname := strings.Replace(filepath.Dir(dpath), "\\", "/", -1)

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	go func() {
		w, _ := session.StdinPipe()
		fmt.Fprintln(w, "C0644", info.Size(), filename)
		io.CopyN(w, file, info.Size())
		fmt.Fprint(w, "\x00")
		w.Close()
	}()

	if err := session.Run(fmt.Sprintf("/usr/bin/scp -qrt %s", dirname)); err != nil {
		return err
	}

	fmt.Printf("%s 发送成功.\n", client.RemoteAddr())
	session.Close()

	if session, err = client.NewSession(); err == nil {
		defer session.Close()
		buf, err := session.Output(fmt.Sprintf("/usr/bin/md5sum %s", dpath))
		if err != nil {
			return err
		}
		fmt.Printf("%s 的MD5:\n%s\n", client.RemoteAddr(), string(buf))
	}
	return nil
}
