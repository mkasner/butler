//SSH client implementation
//Mix from different implementations found on the web
//Implementations on the web mostly use some old versions of ssh library
package butler

import (
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
func (client *SSHClient) Connect() bool {
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
	if err != nil {
		logger.Errorf("Failed to connect to %s", client.Host+":"+client.Port)
		// panic("Failed to dial: " + err.Error())
		return false
	} else {
		logger.Infof("Successful connection to %s", client.Host+":"+client.Port)
		return true
	}
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

func (client *SSHClient) command(s string) string {

	session, err := client.Client.NewSession()
	if err != nil {
		logger.Errorf("Failed to create session: %v", err.Error())
	}
	defer session.Close()

	b, err := session.CombinedOutput(s)
	if err != nil {
		logger.Errorf("Failed to run: %v", err.Error())
	}

	return strip(string(b))
}

func (client *SSHClient) Close() {
	if client.Client != nil {
		if err := client.Client.Close(); err != nil {
			logger.Error("Failed to close client: " + err.Error())
		}
	}
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

// //https://gist.github.com/kiyor/7817632
// package main

// import (
// 	"bytes"

// 	"code.google.com/p/go.crypto/ssh"

// 	"fmt"

// 	"github.com/wsxiaoys/terminal/color"

// 	"io"
// 	"io/ioutil"
// 	"runtime"
// 	"strings"
// 	"sync"
// 	"time"
// )

// func strip(v string) string {
// 	return strings.TrimSpace(strings.Trim(v, "\n"))
// }

// type keychain struct {
// 	keys []ssh.Signer
// }

// func (k *keychain) Key(i int) (ssh.PublicKey, error) {
// 	if i < 0 || i >= len(k.keys) {
// 		return nil, nil
// 	}
// 	return k.keys[i].PublicKey(), nil
// }

// func (k *keychain) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
// 	return k.keys[i].Sign(rand, data)
// }

// func (k *keychain) add(key ssh.Signer) {
// 	k.keys = append(k.keys, key)
// }

// func (k *keychain) loadPEM(file string) error {
// 	buf, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return err
// 	}
// 	key, err := ssh.ParsePrivateKey(buf)
// 	if err != nil {
// 		return err
// 	}
// 	k.add(key)
// 	return nil
// }

// func Sshcmd(hosts []string, command string, background bool, debug bool) {
// 	k := new(keychain)
// 	// Add path to id_rsa file
// 	err := k.loadPEM("¬/.ssh/id_rsa")

// 	if err != nil {
// 		panic("Cannot load key: " + err.Error())
// 	}

// 	// config := ssh.ClientAuth

// 	// auth := make(ssh.Client) ssh.ClientAuthKeyring(k)

// 	// Switch out username
// 	config := &ssh.ClientConfig{
// 		// Change to your username
// 		User: "vagrant",
// 		Auth: ssh.PublicKeys{
// 			ssh.ClientAuthKeyring(k),
// 		},
// 	}

// 	// command = fmt.Sprintf("/usr/bin/sudo bash <<CMD\nexport PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin\n%s\nCMD", command)

// 	if background {
// 		command = fmt.Sprintf("/usr/bin/nohup bash -c \\\n\"%s\" `</dev/null` >nohup.out 2>&1 &", command)
// 	}

// 	if debug {
// 		color.Printf("@{b}%s\n", command)
// 	}

// 	var wg sync.WaitGroup
// 	queue := make(chan Result)
// 	count := new(int)
// 	// Change if you need
// 	workers := 24
// 	var results []Result

// 	// Filter hosts if you need
// 	// Should be like ["hostname1:22","hostname2:22"]
// 	conns := hosts

// 	for _, conn := range conns {
// 		wg.Add(1)
// 		*count++
// 		if debug {
// 			color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
// 		}
// 		for *count >= workers {
// 			time.Sleep(10 * time.Millisecond)
// 		}
// 		go func(h string) {
// 			defer wg.Done()
// 			var r Result

// 			r.Host = h
// 			client, err := ssh.Dial("tcp", h, config)
// 			if err != nil {
// 				color.Printf("@{!r}%s: Failed to connect: %s\n", h, err.Error())
// 				*count--
// 				if debug {
// 					color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
// 				}
// 				return
// 			}

// 			session, err := client.NewSession()
// 			if err != nil {
// 				color.Printf("@{!r}%s: Failed to create session: %s\n", h, err.Error())
// 				*count--
// 				if debug {
// 					color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
// 				}
// 				return
// 			}
// 			defer session.Close()

// 			// This not working, need debug
// 			//session.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin")

// 			var b bytes.Buffer
// 			var e bytes.Buffer
// 			session.Stdout = &b
// 			session.Stderr = &e
// 			if err := session.Run(command); err != nil {
// 				color.Printf("@{!r}%s: Failed to run: %s\n", h, err.Error())
// 				color.Printf("@{!r}%s\n", strip(e.String()))
// 				*count--
// 				if debug {
// 					color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
// 				}
// 				return
// 			}

// 			if !background {
// 				r.Res = strip(b.String())
// 			} else {
// 				r.Res = "command success and out put in remote server's /home/user/nohup.out"
// 			}
// 			color.Printf("@{!g}%s\n", r.Host)
// 			fmt.Println(r.Res)

// 			*count--
// 			if debug {
// 				color.Printf("@{y}%s\t\tcounter %3d\n", conn, *count)
// 			}

// 			runtime.Gosched()
// 			queue <- r
// 		}(conn)
// 	}
// 	go func() {
// 		defer wg.Done()
// 		for r := range queue {
// 			results = append(results, r)
// 		}
// 	}()
// 	wg.Wait()
// 	// Print res if you need
// 	//fmt.Println(results)
// }
