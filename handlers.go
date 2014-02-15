package tn3270

import (
	"bytes"
	"fmt"
	"strings"
	"errors"
)

var e2a string = " \x01\x02\x03\x9c\t\x86\x7f\x97\x8d\x8e\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x9d\x85\x08\x87\x18\x19\x92\x8f\x1c\x1d\x1e\x1f\x80\x81\x82\x83\x84\n\x17\x1b\x88\x89\x8a\x8b\x8c\x05\x06\x07\x90\x91\x16\x93\x94\x95\x96\x04\x98\x99\x9a\x9b\x14\x15\x9e\x1a \xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8[.<(+!&\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1]$*);^-/\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9|,%_>?\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2`:#@'=\"\xc3abcdefghi\xc4\xc5\xc6\xc7\xc8\xc9\xcajklmnopqr\xcb\xcc\xcd\xce\xcf\xd0\xd1~stuvwxyz\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7{ABCDEFGHI\xe8\xe9\xea\xeb\xec\xed}JKLMNOPQR\xee\xef\xf0\xf1\xf2\xf3\\\x9fSTUVWXYZ\xf4\xf5\xf6\xf7\xf8\xf90123456789\xfa\xfb\xfc\xfd\xfe\xff"
var e2d string = " \x01\x02\x03\x9c\t\x86\x7f\x97\x8d\x8e\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x9d\x85\x08\x87\x18\x19\x92\x8f\x1c \x1e\x1f\x80\x81\x82\x83\x84\n\x17\x1b\x88\x89\x8a\x8b\x8c\x05\x06\x07\x90\x91\x16\x93\x94\x95\x96\x04\x98\x99\x9a\x9b\x14\x15\x9e\x1a \xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8[.<(+!&\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1]$*);^-/\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9|,%_>?\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2`:#@'=\"\xc3abcdefghi\xc4\xc5\xc6\xc7\xc8\xc9\xcajklmnopqr\xcb\xcc\xcd\xce\xcf\xd0\xd1~stuvwxyz\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7{ABCDEFGHI\xe8\xe9\xea\xeb\xec\xed}JKLMNOPQR\xee\xef\xf0\xf1\xf2\xf3\\\x9fSTUVWXYZ\xf4\xf5\xf6\xf7\xf8\xf90123456789\xfa\xfb\xfc\xfd\xfe\xff"
var a2e string = " \x01\x02\x037-./\x16\x05%\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13<=2&\x18\x19?'\x1c\x1d\x1e\x1f@O\x7f{[lP}M]\\Nk`Ka\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9z^L~no|\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9J\xe0Z_my\x81\x82\x83\x84\x85\x86\x87\x88\x89\x91\x92\x93\x94\x95\x96\x97\x98\x99\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xc0j\xd0\xa1\x07 !\"#$\x15\x06\x17()*+,\t\n\x1b01\x1a3456\x0889:;\x04\x14>\xe1ABCDEFGHIQRSTUVWXYbcdefghipqrstuvwx\x80\x8a\x8b\x8c\x8d\x8e\x8f\x90\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xca\xcb\xcc\xcd\xce\xcf\xda\xdb\xdc\xdd\xde\xdf\xea\xeb\xec\xed\xee\xef\xfa\xfb\xfc\xfd\xfe\xff"

func E2A(src []byte) (res []byte) {
	res = make([]byte, len(src))
	for i, b := range src {
		res[i] = e2a[b]
	}
	return
}

func A2E(src []byte) (res []byte) {
	res = make([]byte, len(src))
	for i, b := range src {
		res[i] = a2e[b]
	}
	return
}

func E2D(src []byte) (res []byte) {
	res = make([]byte, len(src))
	for i, b := range src {
		res[i] = e2d[b]
	}
	return
}

type TNHandler interface {
	OnTNCommand(byte)
	OnTNArgCommand(byte, byte)
}

type TN3270NegoHandler interface {
	OnTN3270DeviceTypeRequest([]byte, []byte, []byte)
	OnTN3270DeviceTypeIs([]byte, []byte)
	OnTN3270DeviceTypeReject(byte)
	OnTN3270FunctionsIs([]byte)
	OnTN3270FunctionsRequest([]byte)
	OnTN3270SendDeviceType()
}

type TN3270Handler interface {
	OnTN3270Command(byte)
	OnTN3270Text([]byte)
	OnTN3270WCC(byte)
	OnTN3270AID(byte)
	OnTN3270PT()
	OnTN3270IC()
	OnTN3270SF(byte)
	OnTN3270SFE(byte)
	OnTN3270RA(int, byte)
	OnTN3270SBA(int)
	OnTN3270EUA(int)
	OnTN3270Message()
}

type ErrorHandler interface {
	OnError([]byte, int) error
}

type VerboseErrorHandler struct {
}

type Screen []byte

func (h *VerboseErrorHandler) OnError(data []byte, position int) error {
	return errors.New("Decoding error")
}

// TextTN3270Handler is a handler that keeps only track of text display
type TextTN3270Handler struct {
	lines []string
	line  []string
	rows   int

	HandleMessage func(string)
}
func (h *TextTN3270Handler) lineFeed() {
	if len(h.line) > 0 {
		h.lines = append(h.lines, strings.Join(h.line, ""))
	}
	h.line = h.line[0:0]
}
func (h *TextTN3270Handler) OnTN3270Command(byte) {
	// Clear screen for each command
	h.lines = h.lines[0:0]
	h.line = h.line[0:0]
}
func (h *TextTN3270Handler) OnTN3270Text(text []byte) {
	h.line = append(h.line, string(E2A(text)))
}
func (h *TextTN3270Handler) OnTN3270WCC(byte) {
	// Do nothing
}
func (h *TextTN3270Handler) OnTN3270AID(byte) {
	// Do nothing
}
func (h *TextTN3270Handler) OnTN3270PT() {
	h.lineFeed()
}
func (h *TextTN3270Handler) OnTN3270IC() {
	// Do nothing
}
func (h *TextTN3270Handler) OnTN3270SF(byte) {
	// Do nothing
}
func (h *TextTN3270Handler) OnTN3270SFE(byte) {
	// Do nothing
}
func (h *TextTN3270Handler) OnTN3270RA(int, byte) {
	// Do nothing
}
func (h *TextTN3270Handler) OnTN3270SBA(int) {
	h.lineFeed()
}
func (h *TextTN3270Handler) OnTN3270EUA(int) {
	// Do nothing
}
func (h *TextTN3270Handler) OnTN3270Message() {
	h.lineFeed()
	h.HandleMessage(h.String())
}
func (h *TextTN3270Handler) String() string {
	return strings.Join(h.lines, "\n")
}

// MultiHandler is a TN3270 handler that wraps several handlers into one
// all handlers are called for each function
type MultiHandler struct {
	handlers []*MultiHandler
}

func (h* MultiHandler) OnTN3270Command(b byte) {
	for _, h1 := range h.handlers {
		h1.OnTN3270Command(b)
	}
}
func (h* MultiHandler) OnTN3270Text(text []byte) {
	for _, h1 := range h.handlers {
		h1.OnTN3270Text(text)
	}
}
func (h* MultiHandler) OnTN3270WCC(b byte) {
	for _, h1 := range h.handlers {
		h1.OnTN3270WCC(b)
	}
}
func (h* MultiHandler) OnTN3270AID(b byte) {
	for _, h1 := range h.handlers {
		h1.OnTN3270AID(b)
	}
}
func (h* MultiHandler) OnTN3270PT() {
	for _, h1 := range h.handlers {
		h1.OnTN3270PT()
	}
}
func (h* MultiHandler) OnTN3270IC() {
	for _, h1 := range h.handlers {
		h1.OnTN3270IC()
	}
}
func (h* MultiHandler) OnTN3270SF(b byte) {
	for _, h1 := range h.handlers {
		h1.OnTN3270SF(b)
	}
}
func (h* MultiHandler) OnTN3270SFE(b byte) {
	for _, h1 := range h.handlers {
		h1.OnTN3270SFE(b)
	}
}
func (h* MultiHandler) OnTN3270RA(addr int, b byte) {
	for _, h1 := range h.handlers {
		h1.OnTN3270RA(addr, b)
	}
}
func (h* MultiHandler) OnTN3270SBA(addr int) {
	for _, h1 := range h.handlers {
		h1.OnTN3270SBA(addr)
	}
}
func (h* MultiHandler) OnTN3270EUA(addr int) {
	for _, h1 := range h.handlers {
		h1.OnTN3270EUA(addr)
	}
}
func (h* MultiHandler) OnTN3270Message() {
	for _, h1 := range h.handlers {
		h1.OnTN3270Message()
	}
}


// VirtualScreenTN3270Handler is a TN3270 handler that simulates a terminal
// and keeps track of the terminal display as if it was a GUI
type VirtualScreenTN3270Handler struct {
	screen           Screen
	rows, cols       int
	position, cursor int

	HandleMessage func(string)
}

func (h *VirtualScreenTN3270Handler) String() string {
	rows := make([][]byte, h.rows)
	for i := 0; i < h.rows; i++ {
		rows[i] = bytes.TrimRight(E2D(h.screen[i*h.cols:(i+1)*h.cols]), " ")
	}
	// Trim empty lines
	var beg, end int
	var started bool = false
	for i := 0; i < h.rows; i++ {
		if len(rows[i]) != 0 {
			if !started {
				beg = i
			}
			started = true
			end = i
		}
	}
	// Generate screen string
	var sep []byte = []byte{'\n'}
	return string(bytes.Join(rows[beg:end+1], sep))
}

func (h *VirtualScreenTN3270Handler) OnTN3270Command(b byte) {
	switch b {
	case 0x05, 0xf5:
		// Clear screen
		h.screen = make([]byte, h.rows*h.cols)
		h.position = 0
	}
}

func (h *VirtualScreenTN3270Handler) OnTN3270WCC(b byte) {
	// Nothing to be done on WCC
}

func (h *VirtualScreenTN3270Handler) OnTN3270AID(b byte) {
	// Nothing to be done on AID
}

func (h *VirtualScreenTN3270Handler) OnTN3270SF(b byte) {
	h.screen[h.position] = 0x1d
	h.position = (h.position + 1) % (h.rows * h.cols)
}

func (h *VirtualScreenTN3270Handler) OnTN3270SFE(b byte) {
	h.screen[h.position] = 0x1d
	h.position = (h.position + 1) % (h.rows * h.cols)
}

func (h *VirtualScreenTN3270Handler) OnTN3270PT() {
	// Move to next start field
	var idx int
	for i := 0; i < h.rows*h.cols; i++ {
		idx = (h.position + i) % (h.rows * h.cols)
		if h.screen[idx] == 0x1d {
			break
		}
		h.screen[idx] = 0x00
	}
	h.position = (idx + 1) % (h.rows * h.cols)
}

func (h *VirtualScreenTN3270Handler) OnTN3270IC() {
	h.cursor = h.position
}

func (h *VirtualScreenTN3270Handler) OnTN3270Text(b []byte) {
	copy(h.screen[h.position:h.position+len(b)], b)
	h.position = (h.position + len(b)) % (h.rows * h.cols)
}

func (h *VirtualScreenTN3270Handler) OnTN3270RA(addr int, b byte) {
	for ; h.position != addr; h.position = (h.position + 1) % (h.rows * h.cols) {
		h.screen[h.position] = b
	}
}

func (h *VirtualScreenTN3270Handler) OnTN3270SBA(addr int) {
	h.position = addr
}

func (h *VirtualScreenTN3270Handler) OnTN3270EUA(addr int) {
	for ; h.position != addr; h.position = (h.position + 1) % (h.rows * h.cols) {
		h.screen[h.position] = 0x00
	}
}

func (h *VirtualScreenTN3270Handler) OnTN3270Message() {
	h.HandleMessage(h.String())
}

// VerboseTN3270Handler is a handler that prints all the TN3270 commands to
// stdout
type VerboseTN3270Handler struct {
}

func (h *VerboseTN3270Handler) OnTN3270Command(b byte) {
	fmt.Println("TN3270 Command: ", b)
}

func (h *VerboseTN3270Handler) OnTN3270WCC(b byte) {
	fmt.Println("TN3270 WCC: ", b)
}

func (h *VerboseTN3270Handler) OnTN3270AID(b byte) {
	fmt.Println("TN3270 AID: ", b)
}

func (h *VerboseTN3270Handler) OnTN3270SF(b byte) {
	fmt.Println("TN3270 SF: ", b)
}

func (h *VerboseTN3270Handler) OnTN3270SFE(b byte) {
	fmt.Println("TN3270 SFE: ", b)
}

func (h *VerboseTN3270Handler) OnTN3270PT() {
	fmt.Println("TN3270 PT")
}

func (h *VerboseTN3270Handler) OnTN3270IC() {
	fmt.Println("TN3270 IC")
}

func (h *VerboseTN3270Handler) OnTN3270Text(b []byte) {
	fmt.Println("TN3270 Text: ", string(E2A(b)))
}

func (h *VerboseTN3270Handler) OnTN3270RA(addr int, b byte) {
}
func (h *VerboseTN3270Handler) OnTN3270SBA(addr int) {
	fmt.Println("TN3270 SBA: ", addr)
}
func (h *VerboseTN3270Handler) OnTN3270EUA(addr int) {
}
func (h *VerboseTN3270Handler) OnTN3270Message() {
	fmt.Println("End of Message")
}

func NewParser(tnh TNHandler, tn3270negoh TN3270NegoHandler, tn3270h TN3270Handler, errorh ErrorHandler) Parser {
	p := new(parser)
	p.Init()
	p.tnh = tnh
	p.tn3270negoh = tn3270negoh
	p.errorh = errorh
	screen := tn3270h
	p.tn3270h = screen

	return p
}
