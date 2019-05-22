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
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/juju/errors"
	"gopkg.in/tomb.v2"
)

// debugServerConnections controls whether all server connections are wrapped
// with a verbose logging wrapper.
const debugServerConnections = false

// noLimit is an effective infinite upper bound for io.LimitedReader
const noLimit int64 = (1 << 63) - 1

type Request struct {
	Text string
}

type ResponseWriter interface {
	io.Writer
}

type defaultResponseWriter struct {
	headerWrote  bool
	trailerWrote bool
	buf          *bufio.ReadWriter
}

func (w *defaultResponseWriter) Write(s []byte) (n int, e error) {
	var n1 int
	if !w.headerWrote {
		w.headerWrote = true
		n1, e = w.buf.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00})
		n += n1
		n1, e = w.buf.Write([]byte{0xf5, 0xc3, 0x11, 0xc1, 0x50})
		n += n1
	}
	n1, e = w.buf.Write(A2E(s))
	n += n1
	return
}

func (w *defaultResponseWriter) finishRequest() {
	w.buf.Write([]byte{0xff, 0xef})
	w.buf.Flush()
}

// Handler implements the Handler interface can be
// registered to serve a particular path or subtree
// in the HTTP server.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return.  Returning signals that the request is finished
// and that the HTTP server can move on to the next request on
// the connection.
type Handler interface {
	ServeWelcomeScreen(ResponseWriter)
	ServeTN3270(ResponseWriter, *Request)
}

type Server struct {
	Addr      string  // TCP address to listen on, ":telnet" if empty
	Handler   Handler // handler to invoke
	TLSConfig *tls.Config

	t          tomb.Tomb // Manages the go routine spawned by the server
	mu         sync.Mutex
	inShutdown int32 // accessed atomically (non-zero means we're in Shutdown)
	listeners  map[*net.Listener]struct{}
	activeConn map[*conn]struct{}
}

type conn struct {
	remoteAddr string            // network address of remote side
	server     *Server           // the Server on which the connection arrived
	rwc        net.Conn          // i/o connection
	lr         *io.LimitedReader // io.LimitReader(sr)
	buf        *bufio.ReadWriter // buffered(lr,rwc)
	parser     Parser
}

func (c *conn) serve() {
	c.buf.Write([]byte{0xff, 0xfd, 0x28})
	c.buf.Flush()
	for {
		recv_buf := make([]byte, 1024)
		n, _ := c.buf.Read(recv_buf)
		if n == 0 {
			break
		}
		c.parser.Parse(recv_buf[:n])
	}
}

func (c *conn) recv() {
}

type defaultTNHandler struct {
	c    *conn
	text []string
}

func (*defaultTNHandler) OnTNCommand(byte) {
}

func (h *defaultTNHandler) OnTNArgCommand(c byte, a byte) {
	if c == 0xfb && a == 0x28 { // WILL TN3270
		h.c.buf.Write([]byte{0xff, 0xfa, 0x28, 0x08, 0x02, 0xff, 0xf0}) // SEND DEVICE TYPE
	}
	h.c.buf.Flush()
}

func (h *defaultTNHandler) OnTN3270DeviceTypeRequest(device_type []byte, device_name []byte, resource_name []byte) {
	h.c.buf.Write([]byte{0xff, 0xfa, 0x28, 0x02, 0x04})
	h.c.buf.Write(device_type)
	h.c.buf.Write([]byte{0x01})
	h.c.buf.Write(resource_name)
	h.c.buf.Write([]byte{0xff, 0xf0})
	h.c.buf.Flush()
}

func (h *defaultTNHandler) OnTN3270DeviceTypeIs([]byte, []byte) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270DeviceTypeReject(byte) {
	// Not applicable for server
}

func (h *defaultTNHandler) OnTN3270SendDeviceType() {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270FunctionsIs([]byte) {
	w := &defaultResponseWriter{buf: h.c.buf}
	h.c.server.Handler.ServeWelcomeScreen(w)
	w.finishRequest()
}

func (h *defaultTNHandler) OnTN3270FunctionsRequest([]byte) {
	h.c.buf.Write([]byte{0xff, 0xfa, 0x28, 0x03, 0x07, 0xff, 0xf0})
	h.c.buf.Flush()
}

func (h *defaultTNHandler) OnError([]byte, int) error {
	return nil
}

func (h *defaultTNHandler) OnTN3270Command(byte) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270Text(text []byte) {
	h.text = append(h.text, string(E2A(text)))
}

func (h *defaultTNHandler) OnTN3270WCC(byte) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270AID(byte) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270PT() {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270IC() {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270SF(byte) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270SFE(byte) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270RA(int, byte) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270SBA(int) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270EUA(int) {
	// Not applicable for servers
}

func (h *defaultTNHandler) OnTN3270Message() {
	w := &defaultResponseWriter{buf: h.c.buf}
	h.c.server.Handler.ServeTN3270(w, &Request{Text: strings.Join(h.text, "")})
	h.text = h.text[0:0]
	w.finishRequest()
}

func (s *Server) newConn(rwc net.Conn) *conn {
	c := new(conn)
	c.remoteAddr = rwc.RemoteAddr().String()
	c.server = s
	c.rwc = rwc
	h := &defaultTNHandler{c: c, text: make([]string, 0)}
	c.parser = NewParser(h, h, h, h)

	if debugServerConnections {
		c.rwc = newLoggingConn("server", c.rwc)
	}
	c.lr = io.LimitReader(c.rwc, noLimit).(*io.LimitedReader)
	br := bufio.NewReader(c.lr)
	bw := bufio.NewWriter(c.rwc)
	c.buf = bufio.NewReadWriter(br, bw)
	return c
}

// ListenAndServe listens on the TCP network address s.Addr and then
// calls Serve to handle requests on incoming connections.  If
// s.Addr is blank, ":telnet" is used.
func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":telnet"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

func (s *Server) shuttingDown() bool {
	// TODO: replace inShutdown with the existing atomicBool type;
	// see https://github.com/golang/go/issues/20239#issuecomment-381434582
	return atomic.LoadInt32(&s.inShutdown) != 0
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if s.shuttingDown() {
		return errors.New("Server closed")
	}
	addr := s.Addr
	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	return s.ServeTLS(ln, certFile, keyFile)
}

// cloneTLSConfig returns a shallow clone of cfg, or a new zero tls.Config if
// cfg is nil. This is safe to call even if cfg is in active use by a TLS
// client or server.
func cloneTLSConfig(cfg *tls.Config) *tls.Config {
	if cfg == nil {
		return &tls.Config{}
	}
	return cfg.Clone()
}

// ServeTLS accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines perform TLS
// setup and then read requests, calling srv.Handler to reply to them.
func (srv *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {

	config := cloneTLSConfig(srv.TLSConfig)

	configHasCert := len(config.Certificates) > 0 || config.GetCertificate != nil
	if !configHasCert || certFile != "" || keyFile != "" {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
	}

	tlsListener := tls.NewListener(l, config)
	return srv.Serve(tlsListener)
}

var ErrServerClosed = errors.New("Server closed")

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call s.Handler to reply to them.
func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	if !s.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer s.trackListener(&l, false)

	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-s.t.Dying():
				return nil
			default:
			}
			return err
		}
		c := s.newConn(rw)
		s.t.Go(func() error {
			c.serve()
			return nil
		})
	}
}

// trackListener adds or removes a net.Listener to the set of tracked
// listeners.
//
// We store a pointer to interface in the map set, in case the
// net.Listener is not comparable. This is safe because we only call
// trackListener via Serve and can track+defer untrack the same
// pointer to local variable there. We never need to compare a
// Listener from another caller.
//
// It reports whether the server is still up (not Shutdown or Closed).
func (s *Server) trackListener(ln *net.Listener, add bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listeners == nil {
		s.listeners = make(map[*net.Listener]struct{})
	}
	if add {
		if s.shuttingDown() {
			return false
		}
		s.listeners[ln] = struct{}{}
	} else {
		delete(s.listeners, ln)
	}
	return true
}

func (s *Server) trackConn(c *conn, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeConn == nil {
		s.activeConn = make(map[*conn]struct{})
	}
	if add {
		s.activeConn[c] = struct{}{}
	} else {
		delete(s.activeConn, c)
	}
}

func (s *Server) closeListenersLocked() error {
	var err error
	for ln := range s.listeners {
		if cerr := (*ln).Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(s.listeners, ln)
	}
	return err
}

func (s *Server) Close() error {
	atomic.StoreInt32(&s.inShutdown, 1)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.t.Kill(nil)
	err := s.closeListenersLocked()
	for c := range s.activeConn {
		c.rwc.Close()
		delete(s.activeConn, c)
	}
	return err
}

// loggingConn is used for debugging.
type loggingConn struct {
	name string
	net.Conn
}

func newLoggingConn(baseName string, c net.Conn) net.Conn {
	return &loggingConn{
		name: fmt.Sprintf("%s", baseName),
		Conn: c,
	}
}

func (c *loggingConn) Write(p []byte) (n int, err error) {
	log.Printf("%s.Write(%d) = %x", c.name, len(p), p)
	n, err = c.Conn.Write(p)
	log.Printf("%s.Write(%d) = %d, %v", c.name, len(p), n, err)
	return
}

func (c *loggingConn) Read(p []byte) (n int, err error) {
	log.Printf("%s.Read(%d) = ....", c.name, len(p))
	n, err = c.Conn.Read(p)
	log.Printf("%s.Read(%d) = %d, %v", c.name, len(p), n, err)
	return
}

func (c *loggingConn) Close() (err error) {
	log.Printf("%s.Close() = ...", c.name)
	err = c.Conn.Close()
	log.Printf("%s.Close() = %v", c.name, err)
	return
}
