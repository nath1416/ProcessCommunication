package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
	"strings"
)

const EXE_FILE = "./cpu_server"

var received bool = false
var sent bool = false

func startWebServer() {
	fmt.Println("EmulatorGui is running on port 3000")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println(err)
	}
}

func parseStd(line string, isPrefix bool) {
	words := strings.Fields(line)
	ID := words[0]
	switch ID {
	// case "PID":
	// 	fmt.Printf("\t{%s} \t%s %b\n", ID, words[1:], isPrefix)
	// case "READY":
	// 	fmt.Print("Start Cycle\n")
	// case "EndCycle":
	// 	fmt.Print("End cycle\n")
	default:
		fmt.Printf("\t{%s}\t%s %b\n", ID, words[1:], isPrefix)
	}
}

func listenToChildProcess(stdout io.ReadCloser) error {
	buf := bufio.NewReader(stdout)
	fmt.Printf("Stdout:\n")
	for {
		line, isPrefix, err := buf.ReadLine()
		if err != nil {
			return err
		}
		received = true
		parseStd(string(line), isPrefix)
	}
}

func speakToChildProcess(stdin io.WriteCloser) error {
	writer := bufio.NewWriter(stdin)
	for i := 0; i < 5000; i++ {
		time.Sleep(3 * time.Second)
		if _, err := writer.WriteString(fmt.Sprintf("it:,%s\n", i)); err != nil {
			return err
		}
		if err := writer.Flush(); err != nil{
			return err
		}
		sent = true
	}
	return nil
}

func startEmulatorGui() {
	cmd := exec.Command(EXE_FILE)
	stdout, err := cmd.StdoutPipe()
	stdin, err := cmd.StdinPipe()
	defer stdin.Close()
	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		fmt.Println(err)
	}

	go func() { _ = speakToChildProcess(stdin) }()
	go func() { _ = listenToChildProcess(stdout) }()


	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for child process:", err)
	}
}

func main() {
	// startWebServer()
	startEmulatorGui()
}