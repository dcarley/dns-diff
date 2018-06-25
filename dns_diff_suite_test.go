package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDnsDiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DnsDiff Suite")
}
