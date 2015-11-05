package shell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Runner struct{}

type ConnectionOptions struct {
	Username      string
	Port          int
	PrivateKeyPEM []byte
}

func (r *Runner) ConnectAndRun(host, command string, options *ConnectionOptions) (string, error) {
	signer, err := ssh.ParsePrivateKey(options.PrivateKeyPEM)
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: options.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, options.Port), config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %s", err)
	}
	defer client.Close()

	scpSession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Failed to create SCP session: " + err.Error())
	}
	defer scpSession.Close()

	scriptReader := strings.NewReader("command")
	scpError := copy(int64(len(command)), os.FileMode(0777), "test-script", scriptReader, "/tmp/", scpSession)
	if scpError != nil {
		return "", fmt.Errorf("Failed to scp: %s", scpError)
	}

	execSession, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Failed to create session: " + err.Error())
	}
	defer execSession.Close()

	var stdoutBytes bytes.Buffer
	execSession.Stdout = &stdoutBytes
	if err := execSession.Run("/tmp/scripts/test-script"); err != nil {
		return "", fmt.Errorf("Failed running script: " + err.Error())
	}
	return stdoutBytes.String(), nil
}

func copy(size int64, mode os.FileMode, fileName string, contents io.Reader, destination string, session *ssh.Session) error {
	defer session.Close()
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintf(w, "C%#o %d %s\n", mode, size, fileName)
		io.Copy(w, contents)
		fmt.Fprint(w, "\x00")
	}()
	cmd := fmt.Sprintf("scp -t %s", destination)
	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}
