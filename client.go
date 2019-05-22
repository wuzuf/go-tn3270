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

package tn3270

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
)

type Client struct {
	luname string
	parser Parser
	screen TextTN3270Handler
	read   chan []byte
	write  chan []byte
	msgin  chan string
	msgout chan string
}

func (c *Client) recv(conn io.Reader) {
	recv_buf := make([]byte, 2048)
	for {
		n, _ := conn.Read(recv_buf)
		if n == 0 {
			break
		}
		err := c.parser.Parse(recv_buf[:n])
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
	}
}

func (c *Client) send(conn io.Writer) {
	for {
		data := <-c.write
		conn.Write(data)
	}
}

func (c *Client) handle(conn net.Conn) {
	go c.recv(conn)
	go c.send(conn)
}

func (c *Client) Connect(addr string) (chan string, error) {
	var conn net.Conn
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	if debugServerConnections {
		conn = newLoggingConn("server", conn)
	}
	go c.handle(conn)
	return c.msgin, nil
}

func (c *Client) ConnectTLS(addr string, tslconfig *tls.Config) (chan string, error) {
	var conn net.Conn
	var err error
	conn, err = tls.Dial("tcp", addr, tslconfig)
	if err != nil {
		return nil, err
	}
	if debugServerConnections {
		conn = newLoggingConn("server", conn)
	}
	go c.handle(conn)
	return c.msgin, nil
}

func (c *Client) Send(s string) chan string {
	c.write <- []byte{0x00, 0x00, 0x00, 0x00, 0x00}
	c.write <- []byte{0x7d, 0xc1, 0x50, 0x11, 0xc1, 0x50}
	c.write <- A2E([]byte(s))
	c.write <- []byte{0xff, 0xef}
	return c.msgin
}

func (c *Client) SendRecv(s string) string {
	return <-c.Send(s)
}

func (c *Client) OnTNCommand(b byte) {

}

func (c *Client) OnTNArgCommand(b byte, arg byte) {
	if b == 0xfd && arg == 0x28 {
		c.write <- []byte{0xff, 0xfb, 0x28}
	}
}

func (c *Client) OnError([]byte, int) error {
	log.Printf("Error occured")
	return errors.New("Unknown error")
}

func (c *Client) OnTN3270DeviceTypeRequest([]byte, []byte, []byte) {
}

func (c *Client) OnTN3270DeviceTypeIs(model []byte, name []byte) {
	c.write <- []byte("\xff\xfa\x28\x03\x07\x00\x02\x04\xff\xf0")
}

func (c *Client) OnTN3270DeviceTypeReject(byte) {
}

func (c *Client) OnTN3270FunctionsIs([]byte) {
}

func (c *Client) OnTN3270FunctionsRequest(functions []byte) {
	c.write <- []byte("\xff\xfa\x28\x03\x04")
	c.write <- functions
	c.write <- []byte("\xff\xf0")
}

func (c *Client) OnTN3270SendDeviceType() {
	c.write <- []byte{0xff, 0xfa, 0x28, 0x02, 0x07}
	c.write <- []byte("IBM-3278-2-E")
	c.write <- []byte{0x01}
	c.write <- []byte(c.luname)
	c.write <- []byte{0xff, 0xf0}
}

func NewClient(luname string) (c *Client) {
	c = new(Client)
	c.luname = luname
	c.parser = NewParser(c, c, &c.screen, c)
	c.screen.rows = 24
	c.screen.HandleMessage = func(s string) { c.msgin <- s }
	c.read = make(chan []byte)
	c.write = make(chan []byte)
	c.msgin = make(chan string)
	c.msgout = make(chan string)
	return
}
