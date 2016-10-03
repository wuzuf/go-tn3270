// Copyright 2016 Gabriel de Labachelerie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/wuzuf/go-tn3270"
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

	client := tn3270.NewClient("09AA0C72")
	recv, _ := client.Connect("tst-offc.tn3270.1a.amadeus.net:23")
	fmt.Println(<-recv)
	stdin := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := stdin.ReadLine()
		fmt.Println(<-client.Send(string(line)))
	}
}
