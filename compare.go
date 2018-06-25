package main

import (
	"bufio"
	"errors"
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
	client := &dns.Client{}
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

		query := &dns.Msg{}
		query.SetQuestion(name, typeCode)

		var (
			err              error
			priResp, secResp *dns.Msg
		)
		priResp, err = Request(client, query, priNS)
		if err != nil {
			log.Errorln(FailChar, name, typeStr)
			log.Errorln("primary response:", err)
			pass = false
			continue
		}
		secResp, err = Request(client, query, secNS)
		if err != nil {
			log.Errorln(FailChar, name, typeStr)
			log.Errorln("secondary response:", err)
			pass = false
			continue
		}

		log.Println(priResp, secResp)
	}

	if err := scanner.Err(); err != nil {
		log.Error(err)
		pass = false
	}

	return pass
}

func Request(client *dns.Client, query *dns.Msg, ns string) (*dns.Msg, error) {
	resp, _, err := client.Exchange(query, ns)
	if err != nil {
		return nil, err
	}
	if resp.Rcode != dns.RcodeSuccess {
		return nil, errors.New(dns.RcodeToString[resp.Rcode])
	}

	return resp, nil
}
