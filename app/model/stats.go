package model

type Statis struct {
	Comments int
	Articles int64
	Pages    int64
	Files    int
	Version  int
	Readers  int
}

func NewStatis() *Statis {
	s := new(Statis)
	postNum, _ := GetNumberOfPosts(false, false)
	pageNum, _ := GetNumberOfPosts(true, false)
	// s.Comments = len(commentsIndex)
	s.Articles = postNum
	s.Pages = pageNum
	// s.Pages = len(contentsIndex["page"])
	// s.Files = len(files)
	// s.Version = GetVersion().Version
	// s.Readers = len(GetReaders())
	return s
}
