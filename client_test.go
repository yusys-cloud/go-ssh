// Author: yangzq80@gmail.com
// Date: 2021-10-08
//
package ssh

import (
	"log"
	"testing"
)

func initClient() (*SSHClient, error) {
	return NewSSHClient("host", "root", "")
}

var client, err = initClient()

func TestNewSSHClient(t *testing.T) {

	if err != nil {
		t.Fatal("DialWithPasswd err: ", err)
	}
	if err := client.Close(); err != nil {
		t.Fatal("client.Close err: ", err)
	}
}

func TestNewSSHClientWithKey(t *testing.T) {
	client, err := NewSSHClientWithKey("food", "root", "/Users/zqy/.ssh/id_rsa")
	if err != nil {
		t.Fatal("DialWithKey err: ", err.Error())
	}
	if err := client.Close(); err != nil {
		t.Fatal("client.Close err: ", err)
	}
}

func TestSSHClient_ExecuteCmd(t *testing.T) {

	stdout, _, _ := client.ExecuteCmd("pwd")

	if stdout != "/root\n" {
		t.Fatal("error:", stdout)
	}

	stdout, _, _ = client.ExecuteCmd("cd /root\n ls\n")

	log.Println(stdout)

	client.Close()
}

func TestSSHClient_ExecuteTarCmd(t *testing.T) {

	//client.UploadFile("/Users/zqy/Downloads/tmp/chaosblade-1.3.0-linux-amd64.tar.gz","/root/chaos/chaosblade-1.3.0-linux-amd64.tar.gz")

	stdout, stderr, err := client.ExecuteCmd("cd /root/chaos\n tar -zxvf chaosblade-1.3.0-linux-amd64.tar.gz")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(stdout, stderr)

	stdout, stderr, err = client.ExecuteCmd("cd /root/chaos/chaosblade-1.3.0\n ./blade server start --port 9998")

	log.Println(stdout, stderr)

	client.Close()
}


func TestSSHClient_UploadFile(t *testing.T) {
	client.ExecuteCmd("mkdir /root/chaos")
	stdout, stderr, err := client.UploadFile("/Users/zqy/test/tmp2/t.html", "/root/chaos/tt.html")
	if err != nil {
		t.Error(err.Error())
	}

	stdout, stderr, err=client.ExecuteCmd("cat /root/chaos/tt.html")
	log.Println(stdout, stderr)
}
