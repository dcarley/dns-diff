package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/dcarley/dns-diff"
)

var _ = Describe("Main", func() {
	DescribeTable("EnsurePort",
		func(in, out string) {
			Expect(EnsurePort(in)).To(Equal(out))
		},
		Entry("no port suffix", "8.8.8.8", "8.8.8.8:53"),
		Entry("existing port", "8.8.8.8:123", "8.8.8.8:123"),
		Entry("tailing dot", "ns.example.com.", "ns.example.com:53"),
	)
})
