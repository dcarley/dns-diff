package main

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"sort"
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
		query.SetQuestion(dns.Fqdn(name), typeCode)

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

		if reflect.DeepEqual(priResp.Answer, secResp.Answer) {
			log.Infoln(PassChar, name, typeStr)
		} else {
			log.Warnln(FailChar, name, typeStr)
			for _, answer := range priResp.Answer {
				log.Warnln("-", answer)
			}
			for _, answer := range secResp.Answer {
				log.Warnln("+", answer)
			}
			pass = false
		}
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

	sort.Slice(resp.Answer, func(i, j int) bool {
		return resp.Answer[i].String() < resp.Answer[j].String()
	})

	return resp, nil
}
