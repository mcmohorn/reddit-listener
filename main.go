package main

import (
	"fmt"
	"os"
	"os/exec"
	"reddit/listener"
	"runtime"
	"time"

	"github.com/sethjones/go-reddit/v2/reddit"
	"github.com/sourcegraph/conc"
)

func main() {
	fmt.Println("Starting Reddit listener :)")

	// to get started, replace the placeholder variables on the next two lines
	credentials := reddit.Credentials{ID: "{{reddit_app_id}}", Secret: "{{reddit_app_secret}}", Username: "{{reddit_username}}", Password: "{{reddit_password}}"}
    client, _ := reddit.NewClient(credentials, reddit.WithUserAgent("{{reddit_app_name}}"), reddit.WithApplicationOnlyOAuth(true))

	// you can list your subreddits here
	var subreddits = [...]string {"askreddit", "wallstreetbets", "memes", "pics", "nba", "nhl", "mlb"}
	var listeners = []*listener.Listener{}

	var wg conc.WaitGroup
	// start a separate listener object for each subreddit we listen to
    for _, v := range subreddits {
		listener := &listener.Listener{
			RedditClient: client.Subreddit,
			Subreddit: v,
	   }
	   listeners = append(listeners, listener)
        wg.Go(listener.StartListening)
    }


	// showing results
	for range time.Tick(time.Second * 1) {
		CallClear()
		for _, l := range listeners {
			
			fmt.Printf("%+v\n", l.Subreddit)
			for k,u := range l.Users {

				fmt.Printf("%+v made %+v posts with %d total likes and %d total comments\n", k, len(u.Posts), u.CountAllLikes(), u.CountAllComments())
			}
			fmt.Println("")
		}
    }

    wg.Wait()


}


// https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
var clear map[string]func() //create a map for storing clear funcs that differ based on platform

func init() {
	fmt.Println("Initializing Reddit Listener App")
	clear = make(map[string]func()) //Initialize clear function map

    clear["linux"] = func() { 
        cmd := exec.Command("clear")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
	clear["darwin"] = func() { 
        cmd := exec.Command("clear") 
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

func CallClear() {
    value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	fmt.Printf("%v", runtime.GOOS)
    if ok { //if we defined a clear func for that platform:
        value()  //we execute it
    } else { //unsupported platform
        panic("Your platform is unsupported! I can't clear terminal screen :(")
    }
}




