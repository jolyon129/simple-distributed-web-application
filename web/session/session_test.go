package session_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"zl2501-final-project/web/session/sessmanager"
)

var _ = Describe("Session Manager", func() {
	var fakeSessId string
	BeforeEach(func() {
		fakeSessId = "1j28v6loBj65ypDacf5VJxRDXDcDRU8y1RkdXNOu4qo%3D"
	})
	Describe("Call Manager.SessionStart to start a session", func() {
		var manager *sessmanager.Manager
		manager, _ = sessmanager.GetManagerSingleton("memory")
		Context("When a request without cookie", func() {
			It("should inject a new session id into the cookie", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeW := httptest.NewRecorder()
				sess := manager.SessionStart(fakeW, fakeReq)
				cookies, _ := url.QueryUnescape(fakeW.Header().Get("Set-Cookie"))
				Expect(cookies).Should(ContainSubstring(sessmanager.CookieName))
				Expect(cookies).Should(ContainSubstring(sess.SessionID()))
			})
		})
		Context("When a request with illegal session id in cookie ", func() {
			It("should replace sessionId with a new one", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeReq.Header.Set("Cookie", sessmanager.CookieName+"="+fakeSessId)
				fakeW := httptest.NewRecorder()
				manager.SessionStart(fakeW, fakeReq)
				cookies, _ := url.QueryUnescape(fakeW.Header().Get("Set-Cookie"))
				fmt.Fprintln(GinkgoWriter, cookies)
				Expect(cookies).Should(Not(ContainSubstring(sessmanager.CookieName + "=" + fakeSessId)))
			})
		})
		Context("When a request with legal session id in cookie ", func() {
			It("should reuse the same session Id", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeW := httptest.NewRecorder()
				oldSess := manager.SessionStart(fakeW, fakeReq)
				cookies, _ := url.QueryUnescape(fakeW.Header().Get("Set-Cookie"))
				oldSessid := strings.Split(cookies, ";")[0]
				fakeW2 := httptest.NewRecorder()
				fakeReq.Header.Set("Cookie", oldSessid)
				newSess := manager.SessionStart(fakeW2, fakeReq)
				c2 := fakeW2.Header().Get("Set-Cookie")
				sessIdInCookie := strings.Split(c2, ";")[0]
				fmt.Fprintln(GinkgoWriter, "Old Session Id:", oldSessid)
				fmt.Fprintln(GinkgoWriter, "New Session Id:", sessIdInCookie)
				Expect(sessIdInCookie).Should(BeEmpty())
				Expect(newSess.SessionID()).Should(Equal(oldSess.SessionID()))
			})
		})
		Context("With 10 concurrent requests)", func() {
			It("should return 10 different sessions", func() {
				sessArr := make(map[string]bool)
				var mu sync.Mutex
				for i := 0; i < 10; i++ {
					go func() {
						defer GinkgoRecover()
						fakeReq := httptest.NewRequest("GET", "/login", nil)
						fakeW := httptest.NewRecorder()
						sess := manager.SessionStart(fakeW, fakeReq)
						cookies, _ := url.QueryUnescape(fakeW.Header().Get("Set-Cookie"))
						Expect(cookies).Should(ContainSubstring(sessmanager.CookieName))
						Expect(cookies).Should(ContainSubstring(sess.SessionID()))
						mu.Lock()
						Expect(sessArr[sess.SessionID()]).Should(BeZero())
						sessArr[sess.SessionID()] = true
						mu.Unlock()
					}()
				}
			})
		})

	})
	Describe("Modified values in session", func() {
		manager, _ := sessmanager.GetManagerSingleton("memory")
		Context("When setting values in a session", func() {
			It("should be fine", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeW := httptest.NewRecorder()
				sess := manager.SessionStart(fakeW, fakeReq)
				var wg sync.WaitGroup
				wg.Add(3)
				go func() {
					sess.Set("Name", "Zhuolun Li")
					wg.Done()
				}()
				go func() {
					sess.Set("Handle", "jolyon129")
					wg.Done()
				}()
				go func() {
					sess.Set("Subject", "Distributed System")
					wg.Done()
				}()
				wg.Wait()
				name := sess.Get("Name").(string)
				handle := sess.Get("Handle").(string)
				subject := sess.Get("Subject").(string)
				Expect(name).Should(Equal("Zhuolun Li"))
				Expect(handle).Should(Equal("jolyon129"))
				Expect(subject).Should(Equal("Distributed System"))
			})
		})
		Context("When setting values in 2 sessions concurrently", func() {
			It("Should be synchronized", func() {
				fakeReq1 := httptest.NewRequest("GET", "/login", nil)
				fakeReq2 := httptest.NewRequest("GET", "/login", nil)
				fakeW1 := httptest.NewRecorder()
				fakeW2 := httptest.NewRecorder()
				s1 := manager.SessionStart(fakeW1, fakeReq1)
				s2 := manager.SessionStart(fakeW2, fakeReq2)
				var wg sync.WaitGroup
				wg.Add(10)
				for i := 0; i < 5; i++ {
					go func(i int) {
						s1.Set(i, i)
						wg.Done()
					}(i)
					go func(i int) {
						s2.Set(i, i)
						wg.Done()
					}(i)
				}
				wg.Wait()
				for i := 0; i < 5; i++ {
					Expect(s1.Get(i).(int)).Should(Equal(i))
				}
			})
		})
		Context("When delete existed keys", func() {
			It("should be wiped out", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeW := httptest.NewRecorder()
				sess := manager.SessionStart(fakeW, fakeReq)
				sess.Set("Name", "Zhuolun Li")
				sess.Delete("Name")
				Expect(sess.Get("Name")).Should(BeNil())
			})
		})
	})
})
