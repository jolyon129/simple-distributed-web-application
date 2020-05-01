package mux_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "log"
    "zl2501-final-project/raftcluster/mux"
)

var _ = Describe("Mux", func() {
    var trie *mux.Trie
    BeforeSuite(func() {
        trie = mux.NewTrie()
    })
    Describe("Play with the trie", func() {
        Context("Parse a fixed route", func() {
            It("Should return an endpoint with right name", func() {
                trie.Parse("/url/a/b")
                node, err := trie.Lookup("/url/a/b")
                Expect(err).Should(BeNil())
                Expect(node.Name).Should(Equal("b"))
                Expect(node.IsEndPoint).Should(BeTrue())
                n2, err2 := trie.Lookup("/url/a")
                log.Print(err2)
                Expect(err2).ShouldNot(BeNil())
                Expect(n2).Should(BeNil())
            })
        })
        Context("Parse a dynamical route with named parameter", func() {
            It("should return an endpoint with right name(the name of parameter)", func() {
                trie.Parse("/url/a/b/:id")
                node, err := trie.Lookup("/url/a/b/:id")
                Expect(err).Should(BeNil())
                Expect(node.Name).Should(Equal("id"))
                Expect(node.IsEndPoint).Should(BeTrue())
            })
        })
        Context("Lookup a route with suffix /", func() {
            It("should trim the suffix", func() {
                node, err := trie.Lookup("/url/a/b/:id/")
                Expect(err).Should(BeNil())
                Expect(node.Name).Should(Equal("id"))
                Expect(node.IsEndPoint).Should(BeTrue())
            })
        })
    })
})
