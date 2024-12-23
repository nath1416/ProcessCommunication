package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

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

func listenToChildProcess(stdout io.ReadCloser, receivedLine chan<- string) error {
	buf := bufio.NewReader(stdout)
	fmt.Printf("Stdout:\n")
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			return err
		}
		received = true
		receivedLine <- string(line)
	}
}

func speakToChildProcess(stdin io.WriteCloser, sendLine <-chan string) error {
	writer := bufio.NewWriter(stdin)
	for {
		line := <-sendLine

		if _, err := writer.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return err
		}
		if err := writer.Flush(); err != nil {
			return err
		}
	}
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

	receivedLine := make(chan string, 10)
	sendLine := make(chan string, 10)

	go func() { _ = speakToChildProcess(stdin, sendLine) }()
	go func() { _ = listenToChildProcess(stdout, receivedLine) }()

	for {

		select {
		case msg := <-receivedLine:
			parseStd(msg, false)
		case sendLine <- fmt.Sprintf("%s", 1):
			time.Sleep(1 * time.Second)
		default:
		}
	}

	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Println("Error waiting for child process:", err)
	// }
}

func main() {
	// startWebServer()
	startEmulatorGui()
}
