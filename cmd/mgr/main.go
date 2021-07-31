package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {
	// TODO: optimize stdin, bzip compress, io.readcloser
	scriptBytes, err := ioutil.ReadFile("cmd/mgr/scripts/multiply.js")
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("go", "run", "cmd/iso/main.go")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdin.Close()

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stdout = &stderr

	fmt.Println(scriptBytes)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	defer stdin.Close()
	if _, err := stdin.Write(scriptBytes); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(stdout.String())
	fmt.Println(stderr.String())
}
