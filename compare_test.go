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
		pass    bool
		in, out *bytes.Buffer
	)

	BeforeEach(func() {
		pass = false
		in = &bytes.Buffer{}
		out = &bytes.Buffer{}

		format.TruncatedDiff = false
		log.SetOutput(out)
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
	})

	JustBeforeEach(func() {
		pass = Compare(in, "", "")
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
})
