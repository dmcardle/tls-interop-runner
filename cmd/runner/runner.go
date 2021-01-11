package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

const usage = `Usage:

    $ runner [-help]

    $ util -make-root -out root.crt -key-out root.key -host root.com
`

type implementation struct {
	name   string
	client bool
	server bool
}

var implEndpoints = []implementation{
	{
		name:   "boringssl",
		client: true,
		server: true,
	},
	{
		name:   "cloudflare-go",
		client: true,
		server: true,
	},
}

var regressionEndpoints = []implementation{
	{
		name:   "tlsfuzzer",
		client: true,
		server: false,
	},
}

var testCases = []string{
	"dc",
	"ech",
}

func main() {
	log.SetFlags(0)

	buildOut := new(bytes.Buffer)

	// env SERVER_SRC=./impl-endpoints SERVER=boringssl CLIENT_SRC=./regression-endpoints CLIENT=tlsfuzzer docker-compose build
	cmd := exec.Command("docker-compose", "build")
	cmd.Stdout = buildOut
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "SERVER_SRC=./impl-endpoints")
	cmd.Env = append(cmd.Env, "SERVER=boringssl")
	cmd.Env = append(cmd.Env, "CLIENT_SRC=./regression-endpoints")
	cmd.Env = append(cmd.Env, "CLIENT=tlsfuzzer")
	err := cmd.Start()
	if err != nil {
		log.Println(buildOut.String())
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(buildOut.String())
		log.Fatal(err)
	}

	log.Printf("Build process %d complete, exiting\n", cmd.Process.Pid)

	runOutput := new(bytes.Buffer)

	// env SERVER_SRC=./impl-endpoints SERVER=boringssl CLIENT_SRC=./regression-endpoints CLIENT=tlsfuzzer TESTCASE=dc docker-compose up --abort-on-container-exit
	cmd = exec.Command("docker-compose", "up", "--abort-on-container-exit")
	cmd.Stdout = runOutput
	cmd.Env = append(cmd.Env, "SERVER_SRC=./impl-endpoints")
	cmd.Env = append(cmd.Env, "SERVER=boringssl")
	cmd.Env = append(cmd.Env, "CLIENT_SRC=./regression-endpoints")
	cmd.Env = append(cmd.Env, "CLIENT=tlsfuzzer")
	cmd.Env = append(cmd.Env, "TESTCASE=dc")
	err = cmd.Start()
	if err != nil {
		log.Println(runOutput.String())
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		exitError := err.(*exec.ExitError)
		log.Println(runOutput.String())
		log.Fatal(exitError.ProcessState.ExitCode())
	} else {
		log.Println("Success!")
	}
	// log.Printf("Ran interop test %d, exiting\n", cmd.Process.Pid)
	// log.Println(output.String())

	// TODO(caw): run docker-compose and whatnot
}
