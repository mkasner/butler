//Mimics fabric commands
package butler

import (
	"bytes"
	"os/exec"
	"strings"
)

//Execute locally
func Local(command string, workingDir string, capture bool) (string, error) {
	split := strings.Split(command, " ")

	cmd := exec.Command(split[0:1][0], split[1:]...)
	cmd.Dir = workingDir
	var o, e bytes.Buffer
	cmd.Stdout = &o
	cmd.Stderr = &e
	//if command exits with 0 no error is written

	// if capture {
	// 	return o.String(), nil
	// }
	err := cmd.Run()

	return o.String(), err
}

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
