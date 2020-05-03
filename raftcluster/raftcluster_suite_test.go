package raftcluster_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRaftcluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Raftcluster Suite")
}
