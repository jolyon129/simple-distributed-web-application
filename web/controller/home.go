package controller

import (
	"container/list"
	"html/template"
	"log"
	"net/http"
	"zl2501-final-project/web/constant"
)

type tweet struct {
	Content   string
	CreatedAt string
	CreatedBy string
	UserId    int
}

type homeView struct {
	Name     string
	MyTweets []tweet
	Feed     []tweet
}


func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess := globalSessions.SessionStart(w, r)
		uname := sess.Get(constant.UserName).(string)
		//log.Println("THe user name in session is:",uname)
		t, _ := template.ParseFiles(constant.RelativePathForTemplate+"home.html")
		w.Header().Set("Content-Type", "text/html")
		view := homeView{Name: uname, MyTweets: make([]tweet, 0)}
		userRepo := model.GetUserRepo()
		postRepo := model.GetPostRepo()
		userE := userRepo.SelectByName(uname)
		// Iterate in reverse order because the latest one is stored in the tail in DB
		for e := userE.Posts.Back(); e != nil; e = e.Prev() { // Add the tweets into view
			pId := e.Value.(uint)
			postE := postRepo.SelectById(pId)
			view.MyTweets = append(view.MyTweets, tweet{
				Content:   postE.Content,
				CreatedAt: postE.CreatedTime.Format(constant.TimeFormat),
				CreatedBy: userE.UserName,
				UserId:    int(userE.ID),
			})
		}
		// Build feed
		followingUIdList := userE.Following
		postFeed := list.New()
		for e := followingUIdList.Front(); e != nil; e = e.Next() {
			u := userRepo.SelectById(e.Value.(uint))
			postFeed = mergeSortedPostList(postFeed, u.Posts) // Keep merging into the new feed
		}
		postFeed = mergeSortedPostList(postFeed, userE.Posts) // Merge my own tweets
		retFeed := make([]tweet, 0)
		for e := postFeed.Back(); e != nil; e = e.Prev() { // Iterate in revers oder so that the latest comes first
			pid := e.Value.(uint)
			pE := model.GetPostRepo().SelectById(pid)
			uE := model.GetUserRepo().SelectById(pE.UserID)
			retFeed = append(retFeed, tweet{
				Content:   pE.Content,
				CreatedAt: pE.CreatedTime.Format(constant.TimeFormat),
				CreatedBy: uE.UserName,
				UserId:    int(uE.ID),
			})
		}
		log.Println(retFeed)
		view.Feed = retFeed
		t.Execute(w, view)
	}
}

// Given two list of post id, merge them into a new list sorted in time order(oldest-first).
// Return a pointer of the new list of the merged post ids.
func mergeSortedPostList(l1 *list.List, l2 *list.List) *list.List {
	ret := list.New()
	e1 := l1.Front()
	e2 := l2.Front()
	for e1 != nil && e2 != nil {
		pid1 := e1.Value.(uint)
		pid2 := e2.Value.(uint)
		p1 := model.GetPostRepo().SelectById(pid1)
		p2 := model.GetPostRepo().SelectById(pid2)
		if p1.CreatedTime.Before(p2.CreatedTime) {
			ret.PushBack(pid1)
			e1 = e1.Next()
		} else {
			ret.PushBack(pid2)
			e2 = e2.Next()
		}
	}
	for e1 != nil {
		p := e1.Value.(uint)
		ret.PushBack(p)
		e1 = e1.Next()
	}
	for e2 != nil {
		p := e2.Value.(uint)
		ret.PushBack(p)
		e2 = e2.Next()
	}
	return ret
}
