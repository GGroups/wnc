package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var m_ls_str string
var m_client_socket string
var m_verb bool

//整个程序启动复用一个链接对象
var mlistener net.Listener
var mhasLis bool

//websocket链接是双向的，只定义一个变量
var mCnn net.Conn

func showAddr() {
	ifs, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	fmt.Println("#listen on ips: \n{")
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

				fmt.Println("\t" + vs[0])
			}
		}

	}
	fmt.Println("}")
}

func outFun(out io.Writer) {
	fmt.Println("#server output ready")

	for {
		var buf [512]byte
		n, err := mCnn.Read(buf[:])
		if n == 0 {
			if err != nil {
				fmt.Fprintln(os.Stderr, "#OutBuffer=>", err)
				mCnn.Close()
				// ls := <-mls
				// ls.Close()
				if len(m_ls_str) > 1 {
					creServer()
				}
			}
		}
		out.Write(buf[:n])
	}
}

func inputFun(in io.Reader) {
	fmt.Println("#server input ready")
	for {
		var buf [4096]byte
		n, err := in.Read(buf[:])
		if err != nil {
			fmt.Fprintln(os.Stderr, "#InputBuffer=>", err)
			os.Exit(0)
		}
		mCnn.Write(buf[:n])
	}
}

func creServer() {
	var err error

	if !mhasLis {
		mlistener, err = net.Listen("tcp", m_ls_str)
	}

	if err != nil {
		erro := "#creServer " + err.Error()
		panic(erro)
	}
	mhasLis = true
	if m_verb {
		showAddr()
	}

	fmt.Println("#server listening on " + m_ls_str + "...")
	mCnn, err = mlistener.Accept() //block
	if err != nil {
		erro := "#2" + err.Error()
		panic(erro)
	}

	fmt.Fprintln(os.Stdout, "#open <->", mCnn.RemoteAddr().String())
}

func creClient() {
	var err error

	mCnn, err = net.Dial("tcp", m_client_socket)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Error: input -l port")
		return
	}
	var port = flag.Int("l", 0, "listen port on all ips")
	var verb = flag.Bool("v", false, "show ips listening")
	_ = flag.String("nc client mode", "", "wnc.exe ip port")

	flag.Parse() //do the parse on this func

	if *port == 0 {
		//as client
		a := flag.Arg(0)
		p := flag.Arg(1)
		if len(p) <= 0 {
			panic("not enough args")
		}
		m_client_socket = fmt.Sprintf("%s:%s", a, p)

		creClient()
	} else {
		//as server
		m_ls_str = fmt.Sprintf(":%d", *port)
		m_verb = *verb

		creServer()
	}

	//receive buffer
	go outFun(os.Stdout)

	//input buffer
	inputFun(os.Stdin)
}
