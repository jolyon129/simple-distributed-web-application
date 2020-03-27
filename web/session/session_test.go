package session_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
	"strings"
	"zl2501-final-project/web/session/sessmanager"
)

var _ = Describe("Session Manager", func() {
	var manager *sessmanager.Manager
	var fakeSessId string
	BeforeEach(func() {
		fakeSessId = "1j28v6loBj65ypDacf5VJxRDXDcDRU8y1RkdXNOu4qo%3D"
		manager, _ = sessmanager.GetManagerSingleton("memory")
	})
	Describe("Create a new session", func() {
		PContext("When a request without cookie", func() {
			It("should inject a new session id into the cookie", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeW := httptest.NewRecorder()
				sess := manager.SessionStart(fakeW, fakeReq)
				cookies := fakeW.Header().Get("Set-Cookie")
				Expect(cookies).Should(ContainSubstring(sessmanager.CookieName))
				Expect(cookies).Should(ContainSubstring(sess.SessionID()))
			})
		})
		PContext("When a request with illegal session id in cookie ", func() {
			It("should replace sessionId with a new one", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeReq.Header.Set("Cookie", sessmanager.CookieName+"="+fakeSessId)
				fakeW := httptest.NewRecorder()
				manager.SessionStart(fakeW, fakeReq)
				cookies := fakeW.Header().Get("Set-Cookie")
				fmt.Fprintln(GinkgoWriter, cookies)
				Expect(cookies).Should(Not(ContainSubstring(sessmanager.CookieName + "=" + fakeSessId)))
			})
		})
		Context("When a request with legal session id in cookie ", func() {
			It("should replace sessionId with a new one", func() {
				fakeReq := httptest.NewRequest("GET", "/login", nil)
				fakeW := httptest.NewRecorder()
				oldSess := manager.SessionStart(fakeW, fakeReq)
				cookies := fakeW.Header().Get("Set-Cookie")
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
		Context("With concurrent requests)", func() {

		})
	})

})
