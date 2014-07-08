//Mimics fabric commands
package butler

import (
	"os/exec"
	"strings"
)

//Execute locally
func Local(command string, workingDir string, capture bool) string {
	split := strings.Split(command, " ")
	// fmt.Println(split)
	cmd := exec.Command(split[0:1][0], split[1:]...)
	cmd.Dir = workingDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err.Error())
	}
	if capture == true {
		return strip(string(out))
	}
	return ""
}

//Execute remotely on hosts
func Run(client *SSHClient, command string) string {
	commandFull := "cd " + client.CurrPath + " && " + command
	logger.Infof("run: %s\n", commandFull)

	output := client.command(commandFull)

	return output
}

//cd into directory and keep state
func Cd(client *SSHClient, path string) {
	logger.Infof("cd: %s\n", path)
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
