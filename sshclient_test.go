//Using for testing commands
package butler

import (
	"fmt"
	"testing"
)

var localhost = &SSHClient{User: "vagrant", Host: "localhost", Port: "2200"}

func TestListDirectory(t *testing.T) {
	client := localhost
	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}
	commands := []string{"mkdir /vagrant/tmp/repos/dplr.git", "cd /vagrant/tmp/repos/dplr.git", "ls -al"}

	for _, c := range commands {
		t.Log(c)
		result, err := client.command(c)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(fmt.Sprintf("%s", result))
	}

}
