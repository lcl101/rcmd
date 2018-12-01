package core

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/juju/errors"
	"golang.org/x/crypto/ssh"
)

//ErrorAssert 错误断言
func ErrorAssert(err error, assert string) bool {
	return strings.Contains(err.Error(), assert)
}

//Clear 清屏
func Clear() {
	var cmd exec.Cmd
	if "windows" == runtime.GOOS {
		cmd = *exec.Command("cmd", "/c", "cls")
	} else {
		cmd = *exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

//GetExecPath 获取当前路径
func GetExecPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\"`)
	}
	return string(path[0 : i+1]), nil
}

//ZhLen 计算字符宽度（中文）
func ZhLen(str string) int {
	length := 0
	for _, c := range str {
		if unicode.Is(unicode.Scripts["Han"], c) {
			length += 2
		} else {
			length++
		}
	}
	return length
}

//SplitString 根据','或者';'拆分字符串
func SplitString(str string) (strList []string) {
	if strings.Contains(str, ",") {
		strList = strings.Split(str, ",")
	} else if strings.Contains(str, ";") {
		strList = strings.Split(str, ";")
	} else {
		strList = []string{str}
	}
	return
}

//ParseAuthMethods ssh解析鉴权方式
func ParseAuthMethods(passwd, key string) ([]ssh.AuthMethod, error) {
	sshs := []ssh.AuthMethod{}

	if passwd != "" {

		sshs = append(sshs, ssh.Password(passwd))
		return sshs, nil
	}
	method, err := pemKey(key)
	if err != nil {
		return nil, err
	}
	sshs = append(sshs, method)
	return sshs, nil
}

// 解析密钥
func pemKey(key string) (ssh.AuthMethod, error) {
	sshKey := key
	if sshKey == "" {
		sshKey = "~/.ssh/id_rsa"
	}
	sshKey, _ = ParsePath(sshKey)

	pemBytes, err := ioutil.ReadFile(sshKey)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

//ParsePath 解析路径
func ParsePath(path string) (string, error) {
	str := []rune(path)
	firstKey := string(str[:1])

	if firstKey == "~" {
		home, err := home()
		if err != nil {
			return "", err
		}

		return home + string(str[1:]), nil
	} else if firstKey == "." {
		p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		return p + "/" + path, nil
	} else {
		return path, nil
	}
}

func home() (string, error) {
	u, err := user.Current()
	if nil == err {
		return u.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
