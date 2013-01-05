package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/wuzuf/go-tn3270/tn3270"
	"os"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	client := tn3270.NewClient("09123456")
	fmt.Println(client.Connect("localhost:8023"))
	stdin := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := stdin.ReadLine()
		client.Send(string(line))
		fmt.Println(client.Recv())
	}
}
