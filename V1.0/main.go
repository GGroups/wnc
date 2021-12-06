package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func recvClientConn(conn net.Conn, out io.Writer) {
	for {
		var buf [512]byte
		n, err := conn.Read(buf[:])
		if n == 0 {
			if err != nil {
				fmt.Fprintln(os.Stderr, "#4", err)
			}
			break
		}
		out.Write(buf[:n])
	}
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Error: input -l port")
		return
	}
	var port = flag.Int("l", 0, "listen port on all ips")
	var verb = flag.Bool("v", false, "show ips listening")
	_ = flag.String("nc client mode", "", "wnc.exe ip port")

	flag.Parse()
	var conn net.Conn
	var err error

	//as client
	if *port == 0 {
		a := flag.Arg(0)
		p := flag.Arg(1)
		fmt.Println(a, p)
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", a, p))
		if err != nil {
			panic(err)
		}
		defer conn.Close()
	} else { //as server
		var ls net.Listener
		ls, err = net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			erro := "#1" + err.Error()
			panic(erro)
		}
		defer ls.Close()
		if *verb {
			showAddr()
		}
		conn, err = ls.Accept()
		if err != nil {
			erro := "#2" + err.Error()
			panic(erro)
		}
		defer conn.Close()
		fmt.Fprintln(os.Stdout, "#open conn->", conn.RemoteAddr().String())
	}

	//open a receive buffer
	go recvClientConn(conn, os.Stdout)

	//input buffer
	for {
		var buf [4096]byte
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "3#", err)
			break
		}
		conn.Write(buf[:n])
	}
}

func showAddr() {
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	fmt.Println("#listen on ips: {")
	for _, if1 := range ifs {
		addrs, err := if1.Addrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range addrs {
			if strings.HasPrefix(addr.String(), "127.") {
				continue
			} else if strings.Contains(addr.String(), ":") {
				continue
			} else {
				vs := strings.Split(addr.String(), "/")

				fmt.Println(vs[0])
			}
		}

	}
	fmt.Println("}")

}
