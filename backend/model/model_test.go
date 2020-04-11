package model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"zl2501-final-project/backend/model"
	"zl2501-final-project/backend/model/repository"
)

var _ = Describe("Model", func() {
	Context("Call GetUserRepo Multiple times concurrently", func() {
		It("should return the singleton", func() {
			var wg sync.WaitGroup
			arr := make([]*repository.UserRepo, 20)
			wg.Add(10)
			for i := 0; i < 10; i++ {
				go func(i int) {
					defer wg.Done()
					arr[i] = model.GetUserRepo()
				}(i)
			}
			wg.Wait()
			t := arr[0]
			for i := 1; i < 10; i++ {
				Expect(t).Should(BeIdenticalTo(arr[i]))
			}
		})
	})
	Context("Call GetTweetRepo Multiple times concurrently", func() {
		It("should return the singleton", func() {
			var wg sync.WaitGroup
			arr := make([]*repository.TweetRepo, 20)
			wg.Add(10)
			for i := 0; i < 10; i++ {
				go func(i int) {
					defer wg.Done()
					arr[i] = model.GetTweetRepo()
				}(i)
			}
			wg.Wait()
			t := arr[0]
			for i := 1; i < 10; i++ {
				Expect(t).Should(BeIdenticalTo(arr[i]))
			}
		})
	})
})
