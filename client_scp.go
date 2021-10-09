// Author: yangzq80@gmail.com
// Date: 2021-10-08
//
package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	chunkSize = 65536 // chunk size in bytes for scp
)

func (s *SSHClient) UploadFile(localSourceFile string, remoteTargetFile string) (stdout, stderr string, err error) {
	fp, err := os.Open(localSourceFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fp.Close()

	contents, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Println(err.Error())
		return
	}

	session, err := s.Client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	cmd := "cat >'" + strings.Replace(remoteTargetFile, "'", "'\\''", -1) + "'"
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Start(cmd)
	if err != nil {
		return
	}

	for start, maxEnd := 0, len(contents); start < maxEnd; start += chunkSize {

		//todo limit scp
		//<-maxThroughputChan

		end := start + chunkSize
		if end > maxEnd {
			end = maxEnd
		}
		_, err = stdinPipe.Write(contents[start:end])
		if err != nil {
			return
		}
	}

	err = stdinPipe.Close()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	err = session.Wait()
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	return
}