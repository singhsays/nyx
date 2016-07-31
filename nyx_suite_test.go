package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNyx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nyx Suite")
}
