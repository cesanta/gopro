gopro is a simple and stupid protocol analyzer written in
order to see what gdb sends to the serial port (or to another
gdb server via tcp).
I didn't manage to make socat work for me.

It turns out that it's still useful when debugging ESP8266 remotely

The main problem with this tool is that it doesn't close the socket
(or quits) when the inbound http connection is dropped, and it will
happily accept more incoming connections, which is fine for tcp->tcp
forwarding, but it doesn't play well with sockets.
Thus, you should ctrl-c the tool after disconnecting from gdb.

usage:
    gopro -s :1234 -d /dev/tty.SLAB_USBtoUART -b 115200
