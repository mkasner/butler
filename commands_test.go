//Using for testing commands
package butler

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/wsxiaoys/terminal/color"
)

func TestLocalListDirectorySuccess(t *testing.T) {

	commands := []string{"ls"}

	for _, c := range commands {
		// t.Log(c)
		result, err := Local(c, "/home/kmislav/tmp", true)
		if err != nil {
			t.Logf(color.Sprintf("@r%v\n", err))
		}
		assert(t, len(result) > 2, "Result empty.")
		t.Log(result)
		// fmt.Println(fmt.Sprintf("%s", result))
	}

}

func TestLocalGitPushRemote(t *testing.T) {

	commands := []string{"git push ssh://vagrant@127.0.0.1:2200/home/vagrant/repos/AQD.git master"}

	for _, c := range commands {
		// t.Log(c)
		result, err := Local(c, "/home/kmislav/Projects/AQD", true)
		if err != nil {
			t.Errorf(color.Sprintf("@r%v\n", err))
		}
		equals(t, "", result)
		// t.Log(result)
		// fmt.Println(fmt.Sprintf("%s", result))
	}

}

func TestLocalListDirectoryFail(t *testing.T) {

	commands := []string{"lsa"}

	for _, c := range commands {
		// t.Log(c)
		result, err := Local(c, "/home/kmislav/tmp", true)
		if err != nil {
			t.Logf(color.Sprintf("@r%v\n", err))
		}
		equals(t, "", result)
		// t.Log(result)
		// fmt.Println(fmt.Sprintf("%s", result))
	}

}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
