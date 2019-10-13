package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/crueber/twitter-pruner/pruner"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/mkideal/cli"
)

// Pruner is for each of the pruning functions
type Pruner func(*pruner.Client, *twitter.User) error

func main() {
	cli.Run(new(pruner.Env), func(ctx *cli.Context) error {
		twitterEnv := ctx.Argv().(*pruner.Env)
		client, err := twitterEnv.GenerateClient()
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		user, err := Verify(client)
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		var fns []interface{}
		if twitterEnv.InclTweets {
			fns = append(fns, PruneTimeline)
		}
		if twitterEnv.InclLikes {
			fns = append(fns, PruneLikes)
		}

		for _, fn := range fns {
			if twitterEnv.MaxAPICalls <= 0 {
				fmt.Println("Max number of twitter interactions reached for this run.")
				break
			}
			fmt.Printf("Started %v\n", strings.Replace(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), "main.", "", 1))

			f := fn.(func(*pruner.Client, *twitter.User) error)
			err := Pruner(f)(client, user)
			if err != nil {
				fmt.Printf("%+v\n", err)
				os.Exit(1)
			}
		}
		fmt.Println("Done")

		return nil
	})
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}
