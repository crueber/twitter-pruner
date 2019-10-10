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
type Pruner func(*twitter.Client, *twitter.User, *PrunerEnv) error

func main() {
	cli.Run(new(PrunerEnv), func(ctx *cli.Context) error {
		twitterEnv := ctx.Argv().(*PrunerEnv)
		client, err := twitterEnv.GenerateTwitterClient()
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		user, err := Verify(client, twitterEnv)
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		fns := []interface{}{PruneTimeline, PruneLikes, PruneArchive}

		for _, fn := range fns {
			fmt.Printf("Started %v\n", strings.Replace(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), "main.", "", 1))

			f := fn.(func(*twitter.Client, *twitter.User, *PrunerEnv) error)
			err := Pruner(f)(client, user, twitterEnv)
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
