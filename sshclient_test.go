//Using for testing commands
package butler

import (
	"fmt"
	"testing"
)

var localhost = &SSHClient{User: "mislav", Host: "localhost", Port: "22"}

func TestListDirectory(t *testing.T) {
	client := localhost
	client.Connect()

	fmt.Println(fmt.Sprintf("%s", client.command("mkdir /home/mislav/tmp/repos/dplr.git")))
	fmt.Println(fmt.Sprintf("%s", client.command("cd /home/mislav/tmp/repos/dplr.git")))
	fmt.Println(fmt.Sprintf("%s", client.command("ls -al")))

}
