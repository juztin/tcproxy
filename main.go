package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func proxy(in net.Conn, target string) {
	out, err := net.Dial("tcp", target)
	if err != nil {
		in.Close()
		log.Printf("failed to establish outgoing connection; %s\n", err)
		return
	}

	go func() {
		_, err = io.Copy(in, out)
		if err != nil {
			log.Printf("\tin→out; %v\n", err)
		}
		err = in.Close()
		if err != nil {
			log.Printf("\t%v\n", err)
		}
		err = out.Close()
		if err != nil {
			log.Printf("\t%v\n", err)
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Printf("out→in; %v\n", err)
	}
}

func usage(bin string) string {
	return fmt.Sprintf(`
Usage:
  %s source target

Example:
  %s localhost:8080 remote.com:80
`, bin, bin)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "invalid args; expected %d but got %d\n%s\n", 2, len(os.Args)-1, usage(os.Args[0]))
		os.Exit(1)
	}

	source, target := os.Args[1], os.Args[2]
	listener, err := net.Listen("tcp", source)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Fprintf(os.Stdout, "Listening on %s\n", source)
	for {
		in, err := listener.Accept()
		if err != nil {
			log.Printf("failed to establish incoming connection; %v\n", err)
			continue
		}
		go proxy(in, target)
	}
}
