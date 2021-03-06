package grape

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Submission struct {
	Title           string
	Url             string
	NumComments     int `json:"num_comments"`
	Author          string
	IsSelf          bool `json:"is_self"`
	IsNSFW          bool `json:"over_18"`
	SelfText        string
	Created         float64 `json:"created_utc"`
	Score           int
	Ups             int
	Downs           int
	Sub             string `json:"subreddit"`
	AuthorFlairCSS  string `json:"author_flair_css_class"`
	AuthorFlairText string `json:"author_flair_text"`
	LinkFlairCSS    string `json:"link_flair_css_class"`
	LinkFlairText   string `json:"link_flair_text"`
	Clicked         bool
	Domain          string
	Hidden          bool
	Likes           bool
	Permalink       string
	Saved           bool
	Edited          int
	Distinguished   string
	Stickied        bool
	*Thing
}

func (r *Submission) String() string {
	str := fmt.Sprintf(
		"Title: %s\n\t%d Up \n\t%d Down\n\tAuthor: %s\n\tSub: %s\n",
		r.Title,
		r.Ups,
		r.Downs,
		r.Author,
		r.Sub)
	return str
}

func (r *Submission) GetUrl() string {
	return fmt.Sprintf(Config.GetUrl("comment"), r.Sub, r.Id)
}

func (r *Submission) GetComments(s sort, t period) []Comment {
	data := url.Values{
		"sort": {string(s)},
		"t":    {string(t)},
	}
	b, err := makeGetRequest(r.GetUrl(), &data)
	if err != nil {
		panic(err)
	}
	cresp := make([]*commentsResponse, 2)
	err = json.Unmarshal(b, &cresp)
	comments := make([]Comment, len(cresp[1].Data.Children))
	for i, comment := range cresp[1].Data.Children {
		comments[i] = comment.Data.toComment()
	}
	return comments
}

func (r *Submission) PostComment(user *Redditor, body string) error {
	if !user.IsLoggedIn() {
		return notLoggedInError
	}
	data := &url.Values{
		"api_type": {"json"},
		"text":     {body},
		"uh":       {user.ModHash},
		"thing_id": {"t6_" + r.Id},
	}
	b, err := makePostRequest(Config.GetApiUrl("comment"), data)
	if err != nil {
		return err
	}
	return parseSimpleErrorResponse(b)
}
