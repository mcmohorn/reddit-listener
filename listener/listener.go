package listener

import (
	"context"
	"fmt"
	"time"

	"github.com/sethjones/go-reddit/v2/reddit"
)

type RedditPoster interface {
	NewPosts(ctx context.Context, subreddit string, opts *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error)
}

// Listener is the class that will connect to reddit and read periodically
type Listener struct {
	RedditClient RedditPoster
	Subreddit string
	Users map[string]Details
	MaxRequestsPerSecond float32
}

type Details struct {
	Posts map[string]*reddit.Post
}

func (d *Details) CountAllLikes()int {
	res := 0

	for _,p :=range d.Posts {
		res = res + p.Score
	}

	return res
}

func (d *Details) CountAllComments()int {
	res := 0

	for _,p :=range d.Posts {
		res = res + p.NumberOfComments
	}

	return res
}


func (a *Listener) StartListening() {
	fmt.Println("Listener started")
	a.Users = make(map[string]Details)
	a.listen()
}

func (a *Listener)listen() {

	firstId := ""
	firstTime := true
    for range time.Tick(time.Second * 1) {
		id, err := a.GetLatest(firstId)
		if err != nil {
			panic(err)
		}
		if firstTime {
			firstId = id
			firstTime = false
		}
		

    }
}

func (a *Listener) GetLatest(mostRecentId string) (string, error){
	lim := 1
	if mostRecentId != "" {
		lim = 100
	}
	posts, _, err := a.RedditClient.NewPosts(context.Background(), a.Subreddit,  &reddit.ListOptions{
			Limit: lim,
			Before: mostRecentId,
		},
	)
	if err != nil {
		return "", err
	}

	for _, p := range posts {
		_, exists := a.Users[p.Author]
		if exists {
			a.Users[p.Author].Posts[p.FullID] = p
		} else {
			// user doesn't exist so initialize user with the single post
			ps := make(map[string]*reddit.Post)
			ps[p.FullID] = p
			a.Users[p.Author] = Details{Posts: ps}
		}

	}

	/*
	// used for debugging
	if len(posts) > 0 {
		fmt.Printf("Received %d posts in %s.\n",  len(posts), a.Subreddit)
	}
	*/
	result := mostRecentId
	if len(posts) > 0 {
		result = posts[0].FullID
	}

	return result, nil
}