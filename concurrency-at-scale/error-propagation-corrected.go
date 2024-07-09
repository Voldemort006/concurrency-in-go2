package concurrency_at_scale

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError1 struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func (err MyError) Error1() string {
	return err.Message
}

func wrapError1(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{Inner: err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

type LowLevelErr1 struct {
	error
}

func isGloballyExec1(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{wrapError1(err, err.Error())}
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

type intermediateErr1 struct {
	error
}

func runJob1(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec1(jobBinPath)
	if err != nil {
		return intermediateErr{
			error: wrapError1(err, "cannot run job %q: requisite binaries are not available", id),
		}
	} else if isExecutable == false {
		return wrapError1(nil, "cannot run job %q: requisite binaries are not executable", id)
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}

func handleError1(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logId: %v]: \n", key))
	log.Printf("%#v\n", err)
	fmt.Printf("[%v] %v\n", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob1("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(intermediateErr); ok {
			msg = err.Error()
		}

		handleError1(1, err, msg)
	}
}
