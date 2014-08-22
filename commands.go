//Mimics fabric commands
package butler

import (
	"errors"
	"os/exec"
	"strings"
)

//Execute locally
func Local(command string, workingDir string, capture bool) (string, error) {
	split := strings.Split(command, " ")
	// fmt.Println(split)
	cmd := exec.Command(split[0:1][0], split[1:]...)
	cmd.Dir = workingDir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		var stderrBuf []byte
		if _, err := stderr.Read(stderrBuf); err != nil {
			return "", err
		}
		return "", errors.New(string(stderrBuf))
	}
	if capture == true {
		var stdoutBuf []byte
		if _, err := stdout.Read(stdoutBuf); err != nil {
			return "", err
		}
		return strip(string(stdoutBuf)), nil
	}
	return "", nil
}

// //Execute locally
// func Local(command string, workingDir string, capture bool) (string, error) {
// 	split := strings.Split(command, " ")
// 	// fmt.Println(split)
// 	cmd := exec.Command(split[0:1][0], split[1:]...)
// 	cmd.Dir = workingDir
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		logger.Error(err.Error())
// 	}
// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		logger.Error(err.Error())
// 	}
// 	if err := cmd.Wait(); err != nil {
// 		var stderrBuf bytes.Buffer
// 		if _, err := stderr.Read(stderrBuf); err != nil {
// 			return "", err
// 		}
// 		return "", errors.New(stderrBuf.String())
// 	}
// 	if capture == true {
// 		var stdoutBuf bytes.Buffer
// 		if _, err := stdout.Read(stdoutBuf); err != nil {
// 			return "", err
// 		}
// 		return strip(stdoutBuf.String()), nil
// 	}
// 	return "", nil
// }

//Execute remotely on hosts
func Run(client *SSHClient, command string) (string, error) {
	commandFull := "cd " + client.CurrPath + " && " + command

	output, err := client.command(commandFull)

	return output, err
}

//cd into directory and keep state
func Cd(client *SSHClient, path string) {
	if string(path[0]) == "/" {
		client.CurrPath = path
	} else {
		addition := "" //New addition to string
		if string(client.CurrPath[len(client.CurrPath)-1]) != "/" {
			addition += "/"
		}
		addition += path
		client.CurrPath += addition
	}
	// output := client.command("cd " + path)
	// return output
}
