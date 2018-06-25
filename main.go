package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	var (
		priNS = flag.String("pri", "", "Primary nameserver")
		secNS = flag.String("sec", "", "Secondary nameserver")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *priNS == "" || *secNS == "" {
		flag.Usage()
		os.Exit(2)
	}

	if !Compare(os.Stdin, EnsurePort(*priNS), EnsurePort(*secNS)) {
		os.Exit(1)
	}
}

func EnsurePort(addr string) string {
	if !strings.Contains(addr, ":") {
		addr = net.JoinHostPort(strings.TrimSuffix(addr, "."), "53")
	}

	return addr
}
