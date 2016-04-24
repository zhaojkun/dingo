package model

import "github.com/dinever/golf"

type Statis struct {
	Comments int64
	Articles int64
	Pages    int64
	Files    int
	Version  int
	Sessions int
}

func NewStatis(app *golf.Application) *Statis {
	s := new(Statis)
	postNum, _ := GetNumberOfPosts(false, false)
	pageNum, _ := GetNumberOfPosts(true, false)
	commentNum, _ := GetNumberOfComments()

	s.Articles = postNum
	s.Pages = pageNum
	s.Sessions = app.SessionManager.Count()
	s.Comments = commentNum
	// s.Pages = len(contentsIndex["page"])
	// s.Files = len(files)
	// s.Version = GetVersion().Version
	// s.Readers = len(GetReaders())
	return s
}
