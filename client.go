// Author: yangzq80@gmail.com
// Date: 2021-10-08
//
package ssh

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"time"
)

type SSHClient struct {
	Hostname string
	Username string
	Client   *ssh.Client
}

func NewSSHClient(hostname string, username string, password string) (*SSHClient, error) {
	return NewSSHClientWithPort(hostname, username, password, "22")
}

func NewSSHClientWithPort(hostname string, username string, password string, port string) (*SSHClient, error) {
	conf := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         time.Duration(time.Second * 30),
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	client, err := ssh.Dial("tcp", hostname+":"+port, conf)
	if err != nil {
		log.Println("Get ssh channel error:" + err.Error())
		return nil, err
	}

	return &SSHClient{hostname, username, client}, nil
}

func NewSSHClientWithKey(hostname string, username string, privateKeyfile string) (*SSHClient, error) {
	key, err := ioutil.ReadFile(privateKeyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	conf := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	client, err := ssh.Dial("tcp", hostname+":22", conf)
	if err != nil {
		return nil, err
	}

	return &SSHClient{hostname, username, client}, nil
}

func (s *SSHClient) Close() error {
	if s == nil || s.Client == nil {
		return nil
	}
	return s.Client.Close()
}

func (s *SSHClient) ExecuteCmd(cmd string) (stdout, stderr string, err error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return
	}

	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(cmd)

	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	return
}
