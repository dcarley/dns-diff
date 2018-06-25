package main_test

import (
	"bytes"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	log "github.com/sirupsen/logrus"

	. "github.com/dcarley/dns-diff"
)

func getEnv(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return fallback
}

var _ = Describe("Compare", func() {
	var (
		pass     bool
		in, out  *bytes.Buffer
		pri, sec string
	)

	BeforeEach(func() {
		pass = false
		in = &bytes.Buffer{}
		out = &bytes.Buffer{}
		pri = getEnv("DNS_SEC_HOST", "127.0.0.1") + ":" + getEnv("DNS_SEC_PORT", "10053")
		sec = getEnv("DNS_SEC_HOST", "127.0.0.1") + ":" + getEnv("DNS_SEC_PORT", "20053")

		format.TruncatedDiff = false
		log.SetOutput(out)
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
	})

	JustBeforeEach(func() {
		pass = Compare(in, pri, sec)
	})

	Describe("invalid inputs", func() {
		BeforeEach(func() {
			in = bytes.NewBuffer([]byte(`notype.example.com.
unknowntype.example.com. X
extrafield.example.com. IN A
`))
		})

		It("should log errors", func() {
			Expect(out.String()).To(Equal(`level=error msg="✘ notype.example.com."
level=error msg="couldn't parse record"
level=error msg="✘ unknowntype.example.com. X"
level=error msg="unknown record type"
level=error msg="✘ extrafield.example.com. IN A"
level=error msg="couldn't parse record"
`))
			Expect(pass).To(BeFalse())
		})
	})

	Describe("records that don't exist", func() {
		BeforeEach(func() {
			in = bytes.NewBuffer([]byte(`pri-missing.example.com. A
sec-missing.example.com. A
`))
		})

		It("should log errors", func() {
			Expect(out.String()).To(Equal(`level=error msg="✘ pri-missing.example.com. A"
level=error msg="primary response: NXDOMAIN"
level=error msg="✘ sec-missing.example.com. A"
level=error msg="secondary response: NXDOMAIN"
`))
			Expect(pass).To(BeFalse())
		})
	})
})
