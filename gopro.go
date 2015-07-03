// gopro is a simple and stupid protocol analyzer written in
// order to see what gdb sends to the serial port (or to another
// gdb server via tcp).
// I didn't manage to make socat work for me.
//
// It turns out that it's still useful when debugging ESP8266 remotely
//
// The main problem with this tool is that it doesn't close the socket
// (or quits) when the inbound http connection is dropped, and it will
// happily accept more incoming connections, which is fine for tcp->tcp
// forwarding, but it doesn't play well with serial ports.
// Thus, you should ctrl-c the tool after disconnecting from gdb.
//
// usage:
//     gopro -s :1234 -d /dev/tty.SLAB_USBtoUART -b 115200
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/tarm/serial"
)

var (
	src  = flag.String("s", "", "source (listening) ip:port")
	dst  = flag.String("d", "", "destination ip:port or path to serial")
	baud = flag.Int("b", 0, "baud rate, omit if dest is tcp address")
)

func pipe(done chan struct{}, r io.Reader, w io.Writer) {
	t := io.TeeReader(r, os.Stdout)
	io.Copy(w, t)
}

func connect() (io.ReadWriter, error) {
	if *baud == 0 {
		return net.Dial("tcp", *dst)
	}
	c := &serial.Config{Name: *dst, Baud: *baud}
	return serial.OpenPort(c)
}

func handleConnection(s net.Conn) {
	log.Println("Got incoming connection")
	d, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go pipe(done, s, d)
	go pipe(done, d, s)
	<-done
}

func main() {
	flag.Parse()

	if *src == "" || *dst == "" {
		log.Fatalf("-s and -d are mandatory")
	}

	ln, err := net.Listen("tcp", *src)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}
