package mux_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mux Suite")
}
