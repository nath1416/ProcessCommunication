package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	// "os"
)

func startWebServer() {
	// // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// // 	w.Write([]byte(html))
	// // })
	// // err := http.ListenAndServe(":3000", nil)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	fmt.Println("EmulatorGui is running on port 3000")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println(err)
	}
}

func parseStd(line string) {
	words := strings.Fields(line)

	switch words[0] {
	case "Register":
		fmt.Printf("Register: %s\n", words[1:])
	case "StartCycle":
		fmt.Print("Start Cycle\n")
	case "EndCycle":
		fmt.Print("End cycle\n")
	default:
		fmt.Print(words)
	}

}

func startEmulatorGui() {
	cmd := exec.Command("./build/emulator_cpu")
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	buf := bufio.NewReader(stdout)
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		parseStd(string(line))
		// fmt.Println(string(line))
	}
}
func main() {
	startEmulatorGui()
}
