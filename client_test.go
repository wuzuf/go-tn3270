package tn3270_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/wuzuf/go-tn3270"
)

const certPem = `-----BEGIN CERTIFICATE-----
MIIC5TCCAc2gAwIBAgIJANYloZOeF2jKMA0GCSqGSIb3DQEBCwUAMBQxEjAQBgNV
BAMMCWxvY2FsaG9zdDAeFw0xOTA1MjIxODEzMTdaFw0xOTA2MjExODEzMTdaMBQx
EjAQBgNVBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBANYTySoD3TOVfSAUyuWaZvkJMzcT3n1pIeCi++qW3bRPg55YtFlE7W7ULlUF
uyKJh8Cr1eEDYsqb6aZi6f7Sio3nuDWG6bGqi9MqKT5M8YlWyKeGt5sVgR9Ie3Yv
Yi+Tn9UfW7Zqp8wZOYIUfHW7653kqeUZ+FTUPtoG8VrAOEDoY4Jkw28aGpQReqQy
i2N1qgQZrnyYR7TNgsgX2mnLCjCTyj8+GXI4JNIwq2VngYhu2Xxq22DgSKa8P/13
dhqZ2F0kOcdeNBF5ubcmxW8VY5oQXnRSzHRU8lZ5GMtQpX3mk0P8Sk1ayMK5/Qrm
srlm0qBGAsQ2l3EjQavYdGOymV8CAwEAAaM6MDgwFAYDVR0RBA0wC4IJbG9jYWxo
b3N0MAsGA1UdDwQEAwIHgDATBgNVHSUEDDAKBggrBgEFBQcDATANBgkqhkiG9w0B
AQsFAAOCAQEAvJPTTdQ/XliWcGw1jhPcG8PDm/dP+HWaHLMWmAomCr9OD5FsNWVh
NqMoJViJGfmtAOCfjp9WEzCHBmu2u7bNHat5FFR1CAmh46buxg1jEfUKZC8DeDea
rmwqi//H9QmgI3e+pux9/EBU6qtYOyKBWfoV7zhLLAcUdESpch70kAKcZqVRyeXR
BFFyfEUgOV0QCQ6TvtBl+5MBg0jVBWKGH4K7rCREZl1ENROZBMAOyf6FVfpB+q1E
raPdizujeOMnk53wCzs2dSXrWVJFJWJ+gRIw7fiepe6Ag7k3l5vxjTwBCIk8euqP
k6elUBjuzSLJLtVKJ4lXj9oACNZ8fp91MQ==
-----END CERTIFICATE-----
`

const keyPem = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDWE8kqA90zlX0g
FMrlmmb5CTM3E959aSHgovvqlt20T4OeWLRZRO1u1C5VBbsiiYfAq9XhA2LKm+mm
Yun+0oqN57g1humxqovTKik+TPGJVsinhrebFYEfSHt2L2Ivk5/VH1u2aqfMGTmC
FHx1u+ud5KnlGfhU1D7aBvFawDhA6GOCZMNvGhqUEXqkMotjdaoEGa58mEe0zYLI
F9ppywowk8o/PhlyOCTSMKtlZ4GIbtl8attg4EimvD/9d3YamdhdJDnHXjQRebm3
JsVvFWOaEF50Usx0VPJWeRjLUKV95pND/EpNWsjCuf0K5rK5ZtKgRgLENpdxI0Gr
2HRjsplfAgMBAAECggEAezMfzgIzRwB85f2RVtmo8SEOTGAu5tWeHX0upS71oFMy
V+qMv+MfEW0izONjctjbO1Ko37nnHNrleE/sgi4AdvIn3QYrb4fKuFfHLIdzaap8
B41MFQAnMy2vf7B9PQwkK67ERRLURm0t32KSzd68Fj4GWEa913PeR/M+6w88UH1e
8LWyz1gHdFJ4TTHl+LzCelyM2/rz7gz2rAQDVKvEmsd5bzk/8N5k3fWGngrhyTHy
AIbEk9DltvYyyq48s43JZMgOeOJ2b0xYnAEVgdUKYEg5jlxXcG2y24AEUSYltG4o
XFXwGtdRn+hzmpoaOUAPXjx+rgGWZpS/ExTIpGa3oQKBgQD4+DnqBz4o17ie8bna
3J0NkdhRISRYlTNfZ+Bmso7+fCGUro5qpRQNhqgBIT7gfNrdf+DdbVSIsfA+Cadn
JRQ4aoO5/B04KKERR5N7Ca/6DEbMPKcqcm3yRukasK5Hwxaq+2dzbXB00E0t99jo
+dTBehli3isb4q+Rds4xKMHOEwKBgQDcH1Oh2WHigIWf8RCnxDbh2TBK3UJ3hQeM
hbFja8NMvVnTq4v8YeievtMpXBFjvUdKTFsKqtPSzKtccDVVSxUwZaAtr2Qs52Nz
aLHg86V1WvHErBro3UvJdBSYhuSQ+NMGNp6VVB9ZaNA1ma42pwnSmz9MO0uAD+/4
l0xETuaBBQKBgBjMoPguwIJQ+pLagSjL0NkJLgLmyjgIpJVcQ333S0cOFko5GPaG
evjd8N4r8Zdq2GI32q4ztbfoAYYscABWMS1tbrGX61Esut59wrL+xAikMGknoX8Y
5tq7NXzzHGkJhbiCUkutGwaHuShbB8AtIoQjJWQzvReJ/PMAYomDBcsBAoGANptG
4gXNdKUxgQYKfbP9cXWxt0DAdmn3/3JDGUjogCcRG6OY7JlVXdw1AjOm1Lll8BaV
F0ZdmhPQBvSHJoujzAfJ/std7I3SbBTy271VtJFFHOcdHduYK3eyjEwac6RmpUnz
eVQPGt0XmdRwFXrGwwpkX4LuLezGOUM/VkrEgAkCgYEAudno57MI+YeuHIHBDCXP
ozwjh8DJmw358q5XgOssAoZrNKFUm+pZZvX7vUQGmhibCV038pxSuE49FvLDt7ob
5jJu3L/qGj2nRyaB51FihsxeB5SRT0CNqVyBZlDyNo0uH94B1hReeACO2aB8vodo
ogSdyradpbTTZVcdT7T39bE=
-----END PRIVATE KEY-----
`

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

var _ = Describe("TN3270 Client", func() {
	var server *tn3270.Server
	var addr string
	var rootCAs *x509.CertPool

	Describe("Telnet connection", func() {
		BeforeEach(func() {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			Expect(err).To(Succeed())
			addr = listener.Addr().String()
			server = (&tn3270.Server{Handler: &MyHandler{}})
			go server.Serve(listener)
		})

		AfterEach(func() {
			server.Close()
		})

		It("Should connect and get a welcome screen", func() {
			client := tn3270.NewClient("09123456")
			recv, err := client.Connect(addr)
			Expect(err).To(Succeed())
			output := <-recv
			Expect(output).To(Equal("WELCOME TO MY TN3270 SERVER"))
		})

		It("Should get an echo reply", func() {
			client := tn3270.NewClient("09123456")
			recv, err := client.Connect(addr)
			Expect(err).To(Succeed())
			<-recv
			output := client.SendRecv("Hello")
			Expect(output).To(Equal("ECHO: Hello"))
		})
	})

	Describe("Telnet over TLS connection", func() {
		BeforeEach(func() {
			cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
			if err != nil {
				log.Fatal(err)
			}
			listener, err := tls.Listen("tcp", "localhost:0", &tls.Config{Certificates: []tls.Certificate{cert}})
			Expect(err).To(Succeed())
			addr = fmt.Sprintf("localhost:%s", strings.Split(listener.Addr().String(), ":")[1])
			server = (&tn3270.Server{Handler: &MyHandler{}})
			go server.Serve(listener)

			rootCAs, _ = x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}
			Expect(rootCAs.AppendCertsFromPEM([]byte(certPem))).To(BeTrue())
		})

		AfterEach(func() {
			server.Close()
		})

		It("Should connect and get a welcome screen", func() {
			client := tn3270.NewClient("09123456")
			recv, err := client.ConnectTLS(addr, &tls.Config{RootCAs: rootCAs})
			Expect(err).To(Succeed())
			output := <-recv
			Expect(output).To(Equal("WELCOME TO MY TN3270 SERVER"))
		})

		It("Should get an echo reply", func() {
			client := tn3270.NewClient("09123456")
			recv, err := client.ConnectTLS(addr, &tls.Config{RootCAs: rootCAs})
			Expect(err).To(Succeed())
			<-recv
			output := client.SendRecv("Hello")
			Expect(output).To(Equal("ECHO: Hello"))
		})
	})
})
