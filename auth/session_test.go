package auth_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync"
	"time"
	. "zl2501-final-project/auth/sessmanager"
)

var _ = Describe("Session", func() {
	var fakeSessId string
	BeforeEach(func() {
		fakeSessId = "1j28v6loBj65ypDacf5VJxRDXDcDRU8y1RkdXNOu4qo%3D"
	})
	Describe("Call Manager.SessionStart to start a session", func() {
		var manager *Manager
		manager, _ = GetManagerSingleton("memory")
		Context("When the context timeout", func() {
			It("should return timeout err", func() {
				ctx, _ := context.WithTimeout(context.Background(), 400*time.Millisecond)
				time.Sleep(1000*time.Millisecond)
				_,err:=manager.SessionStart(ctx,"")
				Expect(err).NotTo(BeNil())
			})
		})
		Context("When a request with illegal session id ", func() {
			It("should replace sessionId with a new one", func() {
				ctx, _ := context.WithTimeout(context.Background(), 400*time.Millisecond)
				sId, err := manager.SessionStart(ctx, fakeSessId)
				if err != nil {
					Fail(err.Error())
				}
				Expect(sId).ShouldNot(Equal(fakeSessId))
			})
		})
		Context("When a request with legal session id in cookie ", func() {
			It("should reuse the same session Id", func() {
				ctx, _ := context.WithTimeout(context.Background(), 400*time.Millisecond)
				oldSess, _ := manager.SessionStart(ctx, "")
				newSess, _ := manager.SessionStart(ctx, oldSess)
				fmt.Fprintln(GinkgoWriter, "Old Session Id:", oldSess)
				fmt.Fprintln(GinkgoWriter, "New Session Id:", newSess)
				Expect(oldSess).Should(Equal(newSess))
			})
		})
		Context("With 10 concurrent requests", func() {
			It("should return 10 different sessions", func() {
				ctx, _ := context.WithTimeout(context.Background(), 4000*time.Millisecond)
				sessArr := make(map[string]bool)
				var mu sync.Mutex
				for i := 0; i < 10; i++ {
					go func() {
						defer GinkgoRecover()
						sessId, _ := manager.SessionStart(ctx, "")
						Expect(sessId).Should(Not(BeZero()))
						mu.Lock()
						Expect(sessArr[sessId]).Should(BeZero())
						sessArr[sessId] = true
						mu.Unlock()
					}()
				}
			})
		})

	})
	Describe("Modified values in session", func() {
		manager, _ := GetManagerSingleton("memory")
		Context("When the context timeout", func() {
			It("should return timeout err", func() {
				ctx, _ := context.WithTimeout(context.Background(), 400*time.Millisecond)
				time.Sleep(1000*time.Millisecond)
				_,err:= manager.SetValue(ctx,fakeSessId,"Test","tst")
				Expect(err).NotTo(BeNil())
			})
		})
		Context("When setting values in a session", func() {
			It("should be fine", func() {
				ctx, _ := context.WithTimeout(context.Background(), 4000*time.Millisecond)
				sessId, _ := manager.SessionStart(ctx, "")
				var wg sync.WaitGroup
				wg.Add(3)
				go func() {
					manager.SetValue(ctx, sessId, "Name", "Zhuolun Li")
					wg.Done()
				}()
				go func() {
					manager.SetValue(ctx, sessId, "Handle", "jolyon129")
					wg.Done()
				}()
				go func() {
					manager.SetValue(ctx, sessId, "Subject", "Distributed System")
					wg.Done()
				}()
				wg.Wait()
				name, _ := manager.GetValue(ctx, sessId, "Name")
				handle, _ := manager.GetValue(ctx, sessId, "Handle")
				subject, _ := manager.GetValue(ctx, sessId, "Subject")
				Expect(name.(string)).Should(Equal("Zhuolun Li"))
				Expect(handle.(string)).Should(Equal("jolyon129"))
				Expect(subject.(string)).Should(Equal("Distributed System"))
			})
		})
		Context("When setting values in 2 sessions concurrently", func() {
			It("Should be synchronized", func() {
				ctx, _ := context.WithTimeout(context.Background(), 4000*time.Millisecond)
				s1,_ := manager.SessionStart(ctx,"")
				s2,_ := manager.SessionStart(ctx, "")
				var wg sync.WaitGroup
				wg.Add(10)
				for i := 0; i < 5; i++ {
					go func(i int) {
						manager.SetValue(ctx, s1,i,i)
						wg.Done()
					}(i)
					go func(i int) {
						manager.SetValue(ctx, s2,i,i)
						wg.Done()
					}(i)
				}
				wg.Wait()
				for i := 0; i < 5; i++ {
					value,_:= manager.GetValue(ctx, s1, i)
					Expect(value.(int)).Should(Equal(i))
				}
			})
		})
		Context("When delete existed keys", func() {
			It("should be wiped out", func() {
				ctx, _ := context.WithTimeout(context.Background(), 4000*time.Millisecond)
				sessId,_ := manager.SessionStart(ctx,"")
				manager.SetValue(ctx, sessId,"Name", "Zhuolun Li")
				manager.DeleteValue(ctx, sessId,"Name")
				v,err:= manager.GetValue(ctx, sessId,"Name")
				Expect(err).Should(BeNil())
				Expect(v).Should(BeNil())
			})
		})
	})
})