package cmd

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHConn ssh连接结构
type SSHConn struct {
	User       string
	Password   string
	Host       string
	Key        string
	CipherList []string
	Port       int
	client     *ssh.Client
}

// Conn 用于连接主机
func (s *SSHConn) Conn() error {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(s.Password))

	var config ssh.Config
	if len(s.CipherList) == 0 {
		config = ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		}
	} else {
		config = ssh.Config{
			Ciphers: s.CipherList,
		}
	}

	clientConfig := &ssh.ClientConfig{
		User:    s.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	client, err := ssh.Dial("tcp", addr, clientConfig)

	if err != nil {
		return err
	}
	s.client = client

	return nil
}

//Client 获取ssh client
func (s *SSHConn) Client() *ssh.Client {
	return s.client
}

//Session 创建session
func (s *SSHConn) Session() (*ssh.Session, error) {
	// create session
	session, err := s.client.NewSession()
	if err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return nil, err
	}
	return session, nil
}

//Close 关闭连接
func (s *SSHConn) Close() error {
	if s.client != nil {
		err := s.client.Close()
		return err
	}
	return nil
}
