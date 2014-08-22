//Using for testing commands
package butler

import (
	"fmt"
	"testing"
)

func TestLocalListDirectory(t *testing.T) {

	commands := []string{"mkdir -p /home/kmislav/tmp/repos/", "ls -al"}

	for _, c := range commands {
		t.Log(c)
		result, err := Local(c, "/home/kmislav/tmp", true)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(fmt.Sprintf("%s", result))
	}

}
