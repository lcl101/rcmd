package core

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	// "github.com/islenbo/autossh/core"
	// "github.com/islenbo/autossh/core"
	"github.com/juju/errors"
	"golang.org/x/crypto/ssh"
)

type SessionClient struct {
	session *ssh.Session
	writer  io.WriteCloser
	reader  io.Reader
	wg      *sync.WaitGroup
	errors  chan error
}

// NewSessionClient creates a new SessionClient structure from the ssh.Session
// pointer
func NewSessionClient(s *ssh.Session) (*SessionClient, error) {
	writer, err := s.StdinPipe()
	if err != nil {
		return nil, err
	}
	reader, err := s.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	return &SessionClient{
		session: s,
		writer:  writer,
		reader:  reader,
		wg:      &wg,
		errors:  make(chan error),
	}, nil
}

func (c *SessionClient) Receive(rp, rf, lp, lf string) error {
	go c.FileSink(lp, lf)
	go func() {
		c.wg.Add(1)
		defer c.wg.Done()
		err := c.session.Run(fmt.Sprintf("/usr/bin/scp -f %s", path.Join(rp, rf)))
		if err != nil {
			c.errors <- err
			// panic(err)
		}
	}()
	for err := range c.errors {
		Log.Category("main").Info("Error receiving the file", err)
	}
	c.wg.Wait()
	return nil
}

func (c *SessionClient) Send(rp, rf, lp, lf string) error {
	go c.FileSource(lp, lf)

	Log.Category("main").Info("Beginning transfer")
	go func() {
		c.wg.Add(1)
		defer c.wg.Done()
		err := c.session.Run(fmt.Sprintf("/usr/bin/scp -t %s", path.Join(rp, rf)))
		if err != nil {
			c.errors <- err
			// panic(err)
		}
	}()
	Log.Category("main").Info("Waiting for transfer to complete...")

	for err := range c.errors {
		fmt.Println(err)
		Log.Category("main").Info("Error sending the file", err)
	}
	c.wg.Wait()
	return nil
}

// FileSink is used to receive a file from the remote machine and save it to the local machine
func (c *SessionClient) FileSink(lp, lf string) {
	// We must close the channel for the main thread to work properly. Defer ensures
	// when the function ends, this is closed. We also want to be sure we mark this
	// attempt as completed
	c.wg.Add(1)
	defer c.wg.Done()
	defer close(c.errors)

	Log.Category("main").Info("Beginning transfer")
	successfulByte := []byte{0}
	// Send a null byte saying that we are ready to receive the data
	c.writer.Write(successfulByte)

	// We want to first receive the command input from remote machine
	// e.g. C0644 113828 test.csv
	scpCommandArray := make([]byte, 500)
	bytesRead, err := c.reader.Read(scpCommandArray)
	if err != nil {
		if err == io.EOF {
			//no problem.
		} else {
			c.errors <- err
			return
		}
	}

	scpStartLine := string(scpCommandArray[:bytesRead])
	scpStartLineArray := strings.Split(scpStartLine, " ")

	filePermission := scpStartLineArray[0][1:]
	fileSize := scpStartLineArray[1]
	fileName := scpStartLineArray[2]

	Log.Category("main").Info("File with permissions: %s, File Size: %s, File Name: %s", filePermission, fileSize, fileName)

	// Confirm to the remote host that we have received the command line
	c.writer.Write(successfulByte)

	// Now we want to start receiving the file itself from the remote machine
	// one byte at a time
	fileContents := make([]byte, 1)

	var file *os.File
	if lf == "" {
		file, err = Create(path.Join(lp, fileName))
		if err != nil {
			c.errors <- err
			return
		}
	} else {
		file, err = Create(path.Join(lp, lf))
		if err != nil {
			c.errors <- err
			return
		}
	}
	defer file.Close()

	more := true
	for more {
		bytesRead, err = c.reader.Read(fileContents)
		if err != nil {
			if err == io.EOF {
				more = false
			} else {
				c.errors <- err
				return
			}
		}
		_, err = WriteBytes(file, fileContents[:bytesRead])
		if err != nil {
			c.errors <- err
			return
		}
		c.writer.Write(successfulByte)
	}
	err = file.Sync()
	if err != nil {
		c.errors <- err
		return
	}
}

// FileSource allows us to acting as the machine sending a file to the remote host
func (c *SessionClient) FileSource(lp, lf string) {
	c.wg.Add(1)
	defer c.wg.Done()
	response := make([]byte, 1)
	defer close(c.errors)
	defer c.writer.Close()

	Log.Category("main").Info("Opening file to send")
	f, err := os.Open(path.Join(lp, lf))
	if err != nil {
		c.errors <- err
		return
	}
	defer f.Close()

	Log.Category("main").Info("Getting file information")
	i, err := f.Stat()
	if err != nil {
		c.errors <- err
		return
	}
	// fmt.Println(fmt.Sprintf("C%#o %d %s\n", i.Mode(), i.Size(), i.Name()))
	begin := []byte(fmt.Sprintf("C%#o %d %s\n", i.Mode(), i.Size(), i.Name()))
	_, err = c.writer.Write(begin)
	if err != nil {
		c.errors <- err
		return
	}

	c.reader.Read(response)
	if err != nil {
		c.errors <- err
		return
	}

	io.Copy(c.writer, f)

	fmt.Fprint(c.writer, "\x00")

	_, err = c.reader.Read(response)
	if err != nil {
		c.errors <- err
	}
}

// Create is used to create a file with a specific name
func Create(fn string) (*os.File, error) {
	tfn := strings.TrimSpace(fn)
	f, err := os.Create(tfn)
	if err != nil {
		return f, err
	}

	return f, nil
}

// WriteBytes is used to write an array of bytes to a file
func WriteBytes(file *os.File, content []byte) (int, error) {
	w, err := file.Write(content)
	if err != nil {
		return 0, err
	}
	return w, nil
}

// ExpandPath is used to ensure a path that we have is fully expanded rather
// than something like ./LICENSE
func ExpandPath(f string) (string, error) {
	if !path.IsAbs(f) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		fname := path.Clean(path.Join(wd, f))
		return fname, nil
	}
	return f, nil
}

//Scp 实现scp功能 上传
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

// func ScpDown(client *ssh.Client, spath, dpath string) error {
// 	if dpath == "" || spath == "" {
// 		flag.PrintDefaults()
// 		return errors.New("dpath,spath为空")
// 	}

// 	session, err := client.NewSession()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
