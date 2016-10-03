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
    "flag"
    "fmt"
    "io"
    "github.com/wuzuf/go-tn3270"
    "log"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type MyHandler struct {

}

func (*MyHandler) ServeTN3270(w tn3270.ResponseWriter, r *tn3270.Request) {
    io.WriteString(w, "ECHO: ")
    io.WriteString(w, r.Text)
    log.Printf("Text: %s", r.Text)
}

func (*MyHandler) ServeWelcomeScreen(w tn3270.ResponseWriter) {
    io.WriteString(w, "WELCOME TO MY TN3270 SERVER")
}

func main() {
    flag.Parse() // Scan the arguments list

    if *versionFlag {
        fmt.Println("Version:", APP_VERSION)
    }
    h := &MyHandler{}
    server := tn3270.Server{Addr: ":10023", Handler: h}
    e := server.ListenAndServe()
    if(e != nil) {
        fmt.Println(e)
    }
}
