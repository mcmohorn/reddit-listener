package listener

import (
	"context"
	"errors"
	"testing"

	"github.com/sethjones/go-reddit/v2/reddit"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestStartListening(t *testing.T) {
	mock := &RedditClientMockGood{}
	lst := Listener{
		RedditClient: mock,
		Subreddit: "memes",
   }

   lst.Users = make(map[string]Details)

   lst.GetLatest("")
}

func TestStartListeningBadService(t *testing.T) {
	mock := &RedditClientMockBad{}
	lst := Listener{
		RedditClient: mock,
		Subreddit: "memes",
   }

   lst.Users = make(map[string]Details)

   lst.GetLatest("")
}

func TestCounting(t *testing.T) {
	p := reddit.Post{NumberOfComments: 1, Score: 1}
	ps := make(map[string]*reddit.Post)
	ps["abc"] = &p
	d := Details{
		Posts: ps,
	}
	d.CountAllComments()
	d.CountAllLikes()
}


type RedditClientMockBad struct {

}

func (c *RedditClientMockBad) NewPosts(ctx context.Context, subreddit string, opts *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error) {
    return []*reddit.Post{}, &reddit.Response{}, errors.New("This is an important error, reddit must have failed")
}

type RedditClientMockGood struct {

}

func (c *RedditClientMockGood) NewPosts(ctx context.Context, subreddit string, opts *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error) {
	
    return []*reddit.Post{&reddit.Post{}, &reddit.Post{}}, &reddit.Response{}, nil
}
