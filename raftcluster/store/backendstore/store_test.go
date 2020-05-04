package backendstore_test

import (
    "bytes"
    "encoding/gob"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    . "zl2501-final-project/raftcluster/store"
)

var _ = Describe("DBStore", func() {

    Context("Encode and Decode a commandLog", func() {
        It("should get the right parameter and targetMethod", func() {
            var buf bytes.Buffer
            var command CommandLog
            command = CommandLog{
                ID: 1, TargetMethod: METHOD_SessionInit,
                Params: SessionProviderParams{Sid: "testsid"}}
            gob.Register(SessionProviderParams{})
            if err := gob.NewEncoder(&buf).Encode(command); err != nil {
                Fail(err.Error())
            }
            var nCommand CommandLog
            dec := gob.NewDecoder(&buf)
            err := dec.Decode(&nCommand)
            Expect(err).Should(BeNil())
            if sess, ok := nCommand.Params.(SessionProviderParams); !ok {
                Fail("could not convert the Params to concrete type")
            } else {
                Expect(sess.Sid).Should(Equal("testsid"))
            }
            Expect(nCommand.ID).Should(BeNumerically("==", 1))
            Expect(nCommand.TargetMethod).Should(Equal(METHOD_SessionInit))
        })

    })

})
