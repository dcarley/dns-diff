package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

const (
	PassChar = "✔"
	FailChar = "✘"
)

func Compare(in io.Reader, priNS, secNS string) (pass bool) {
	pass = true
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		rr := strings.Fields(scanner.Text())
		if len(rr) != 2 {
			log.Errorln(FailChar, scanner.Text())
			log.Errorln("couldn't parse record")
			pass = false
			continue
		}

		var typeCode uint16
		name, typeStr := rr[0], rr[1]
		if code, ok := dns.StringToType[typeStr]; ok {
			typeCode = code
		} else {
			log.Errorln(FailChar, scanner.Text())
			log.Errorln("unknown record type")
			pass = false
			continue
		}

		log.Println(name, typeCode)
	}

	if err := scanner.Err(); err != nil {
		log.Error(err)
		pass = false
	}

	return pass
}
