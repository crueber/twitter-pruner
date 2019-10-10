package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/mkideal/cli"
)

// Pruner is for each of the pruning functions
type Pruner func(*twitter.Client, *TwitterEnv) error

func main() {
	cli.Run(new(TwitterEnv), func(ctx *cli.Context) error {
		twitterEnv := ctx.Argv().(*TwitterEnv)
		client := twitterEnv.GenerateTwitterClient()

		fns := []interface{}{Verify, PruneTimeline, PruneLikes, PruneArchive}

		for _, fn := range fns {
			fmt.Printf("Started %v\n", strings.Replace(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), "main.", "", 1))

			f := fn.(func(*twitter.Client, *TwitterEnv) error)
			err := Pruner(f)(client, twitterEnv)
			if err != nil {
				fmt.Println(err)
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
