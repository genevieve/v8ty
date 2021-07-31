package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"rogchap.com/v8go"
)

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	var buf []byte
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			buf = append(buf, scanner.Bytes()...)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("received sript: %s!\n", buf)
	} else {
		log.Fatal("no input")
	}

	iso, _ := v8go.NewIsolate()    // creates a new JavaScript VM
	ctx, _ := v8go.NewContext(iso) // new context within the VM

	vals := make(chan *v8go.Value, 1)
	errs := make(chan error, 1)

	go func() {
		_, err := ctx.RunScript(string(buf), "math.js")
		if err != nil {
			errs <- err
			return
		}

		val, err := ctx.RunScript("multiply(3, 4)", "main.js")
		if err != nil {
			errs <- err
			return
		}
		vals <- val
	}()

	select {
	case val := <-vals:
		fmt.Printf("result: %v\n", val)
	case err := <-errs:
		if err != nil {
			log.Fatal(err)
		}
	case <-time.After(200 * time.Millisecond):
		vm, _ := ctx.Isolate()  // get the Isolate from the context
		vm.TerminateExecution() // terminate the execution
		err := <-errs           // will get a termination error back from the running script
		if err != nil {
			log.Fatal(err)
		}
	}
}
