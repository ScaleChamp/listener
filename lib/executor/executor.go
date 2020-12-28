package executor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/pkg/errors"
	"gitlab.com/scalablespace/listener/app/models"
	"golang.org/x/crypto/ssh"
	"gopkg.in/alessio/shellescape.v1"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Executor struct {
	sshClientConfig map[string]*ssh.ClientConfig
}

const maxTimesWait = 100

func (s *Executor) Wait(node *models.Node) (err error) {
	for i := 0; i < maxTimesWait; i += 1 {
		if err = s.Run(node, "echo"); err != nil {
			log.Println("waiting", err)
			time.Sleep(2 * time.Second)
			continue
		}
		return
	}
	return
}

func (s *Executor) MultiRun(node *models.Node, commands []string) error {
	for _, command := range commands {
		if err := s.Run(node, command); err != nil {
			return err
		}
	}
	return nil
}

func (s *Executor) RetryRunNT(node *models.Node, command string, n int, t time.Duration) (err error) {
	for i := 0; i < n; i += 1 {
		time.Sleep(t)
		if err = s.Run(node, command); err != nil {
			continue
		}
		return
	}
	return
}

func (s *Executor) RetryRun(node *models.Node, command string) (err error) {
	return s.RetryRunNT(node, command, 90, time.Minute)
}

func (s *Executor) Run(node *models.Node, command string, a ...interface{}) error {
	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", node.Metadata.IPv4), s.sshClientConfig[node.Cloud])
	if err != nil {
		return errors.Wrapf(err, "executor failed to dial for instance %s, node %s: %s", node.InstanceId, node.Id, command)
	}
	session, err := c.NewSession()
	if err != nil {
		return errors.Wrapf(err, "executor failed to new-session for instance %s, node %s: %s", node.InstanceId, node.Id, command)
	}
	escaped := make([]interface{}, len(a))
	for i, k := range a {
		switch key := k.(type) {
		case string:
			k = shellescape.Quote(key)
		}
		escaped[i] = k
	}
	if !sudo(node) {
		command = strings.Replace(command, "sudo ", "", -1)
	}
	data, data2, err := execCmd(session, fmt.Sprintf(command, escaped...))
	log.Println("ssh-data: ", data.String())
	log.Println("ssh-err: ", data2.String())
	if err != nil {
		return errors.Wrapf(err, "executor failed to exec cmd for instance %s, node %s: %s", node.InstanceId, node.Id, command)
	}
	return err
}

func sudo(node *models.Node) bool {
	switch node.Cloud {
	case "gcp", "azure", "exoscale":
		return true
	default:
		return false
	}
}

type LenReader interface {
	io.Reader
	Len() int
}

func (s *Executor) PutString(node *models.Node, contents string, remotePath string) error {
	return s.Put(node, bytes.NewBufferString(contents), remotePath)
}

func (s *Executor) PutBytes(node *models.Node, contents []byte, remotePath string) error {
	return s.Put(node, bytes.NewBuffer(contents), remotePath)
}

func (s *Executor) Put(node *models.Node, contents LenReader, remotePath string) error {
	client := scp.NewClient(fmt.Sprintf("%s:22", node.Metadata.IPv4), s.sshClientConfig[node.Cloud])
	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Close()
	return client.Copy(contents, remotePath, "0655", int64(contents.Len()))
}

const (
	ECHO          = 53
	TTY_OP_ISPEED = 128
	TTY_OP_OSPEED = 129
)

func execCmd(session *ssh.Session, cmd string) (*bytes.Buffer, *bytes.Buffer, error) {
	modes := ssh.TerminalModes{
		ECHO:          0,
		TTY_OP_ISPEED: 14400,
		TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 100, 80, modes); err != nil {
		return new(bytes.Buffer), new(bytes.Buffer), err
	}
	stdoutbuf := new(bytes.Buffer)
	stdoutErr := new(bytes.Buffer)
	session.Stdout = stdoutbuf
	session.Stderr = stdoutErr
	err := session.Run(cmd)
	if err != nil {
		return stdoutbuf, stdoutErr, err
	}
	return stdoutbuf, stdoutErr, nil
}

const sshTimeout = 35 * time.Second

func NewExecutor(env models.Environment) (*Executor, error) {
	key, err := ioutil.ReadFile(env.SecretKeyPath)
	if err != nil {
		return nil, err
	}
	var signer ssh.Signer

	if env.SecretKeyPassword != "" {
		data, err := base64.StdEncoding.DecodeString(env.SecretKeyPassword)
		if err != nil {
			return nil, err
		}
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, data)
		if err != nil {
			return nil, err
		}
	} else {
		signer, err = ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
	}
	defaultConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         sshTimeout,
	}
	return &Executor{
		sshClientConfig: map[string]*ssh.ClientConfig{
			"ibm":     defaultConfig,
			"do":      defaultConfig,
			"scw":     defaultConfig,
			"linode":  defaultConfig,
			"upcloud": defaultConfig,
			"hetzner": defaultConfig,
			"aws":     defaultConfig,
			"vultr":   defaultConfig,
			"gcp": {
				User:            "mikefaraponov",
				Auth:            defaultConfig.Auth,
				HostKeyCallback: defaultConfig.HostKeyCallback,
				Timeout:         defaultConfig.Timeout,
			},
			"azure": {
				User:            "scalechamp",
				Auth:            defaultConfig.Auth,
				HostKeyCallback: defaultConfig.HostKeyCallback,
				Timeout:         defaultConfig.Timeout,
			},
			"exoscale": {
				User:            "ubuntu",
				Auth:            defaultConfig.Auth,
				HostKeyCallback: defaultConfig.HostKeyCallback,
				Timeout:         defaultConfig.Timeout,
			},
		},
	}, nil
}
