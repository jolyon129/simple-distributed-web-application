package memory_test

import (
    "encoding/json"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    authstorage "zl2501-final-project/raftcluster/store/authstore"
    _ "zl2501-final-project/raftcluster/store/authstore/memory"
)

var _ = Describe("Memory", func() {
    var provider authstorage.ProviderInterface
    fakeSessId := "1j28v6loBj65ypDacf5VJxRDXDcDRU8y1RkdXNOu4qo%3D"
    fakeSessId2 := "1j28v6loBj65ypDacf5VJxRDXDcDRU8y1RkdXNOu4qo%2F"
    BeforeSuite(func() {
        p, _ := authstorage.GetProvider("memory")
        provider = p
    })
    Describe("Marshal session provider", func() {
        Context("When have a list of sessions", func() {
            It("should return a JSON list with map object inside", func() {
                provider.SessionInit(fakeSessId)
                provider.SessionInit(fakeSessId2)
                sess, _ := provider.SessionRead(fakeSessId)
                sess.Set("Name", "Zhuolun")
                sess.Set("Twitter", "jolyon129")
                sess.Set("id", 3334)
                sess1, _ := provider.SessionRead(fakeSessId2)
                sess1.Set("ins", "jolyon_z")
                sess1.Set("gender", "male")
                j, err := json.Marshal(provider)
                js := string(j)
                //println(js)
                Expect(js).Should(ContainSubstring(`"ins":"jolyon_z"`), ContainSubstring(`"id":3334`))
                Expect(err).Should(BeNil())
            })
        })
    })
})
