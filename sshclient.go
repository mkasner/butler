//SSH client implementation
//Mix from different implementations found on the web
//Implementations on the web mostly use some old versions of ssh library
package butler

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"strings"

	"code.google.com/p/go.crypto/ssh"
)

type SSHClient struct {
	User string
	Host string
	Port string
	//Holds current path for this session
	//I must remember this value because, path is lost after every command execute
	//Here are docs that explain this in fabric https://fabric.readthedocs.org/en/1.3.3/api/core/context_managers.html
	CurrPath string
	Client   *ssh.Client
}

//First function when creating a new client
func (client *SSHClient) Connect() error {
	pkey := parsekey(UserHomeDir() + "/.ssh/id_rsa")

	//
	config := &ssh.ClientConfig{
		User: client.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(pkey),
		},
	}
	var err error
	client.Client, err = ssh.Dial("tcp", client.Host+":"+client.Port, config)
	return err
	// if err != nil {
	// 	logger.Errorf("Failed to connect to %s", client.Host+":"+client.Port)
	// 	// panic("Failed to dial: " + err.Error())
	// 	return false
	// } else {
	// 	logger.Infof("Successful connection to %s", client.Host+":"+client.Port)
	// 	return true
	// }
}

func parsekey(file string) ssh.Signer {
	privateBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}
	return private
}

func (client *SSHClient) command(s string) (string, error) {
	var err error
	session, err := client.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var o, e bytes.Buffer
	session.Stdout = &o
	session.Stderr = &e
	if err := session.Run(s); err != nil {
		return o.String(), errors.New(e.String())
	}
	return o.String(), nil

}

func handleWait(session *ssh.Session, stderr io.Reader, errr error) {
	if err := session.Wait(); err != nil {
		var stderrBuf bytes.Buffer
		in := bufio.NewScanner(stderr)
		for in.Scan() {
			stderrBuf.Write(in.Bytes())
			stderrBuf.Write([]byte("\n"))
		}
		errr = errors.New(strip(string(stderrBuf.String())))
		// return "",

		// return "", err
	}
}

func (client *SSHClient) Close() error {
	if client.Client != nil {
		err := client.Client.Close()
		return err

	}
	return nil
}

//Constructs full host name with user@host:port
func (client *SSHClient) FullHost() string {
	hostArray := make([]string, 5)
	if client.User != "" {
		hostArray[0] = client.User
		hostArray[1] = "@"
	}
	hostArray[2] = client.Host
	hostArray[3] = ":"
	hostArray[4] = client.Port

	return strings.Join(hostArray, "")

}
