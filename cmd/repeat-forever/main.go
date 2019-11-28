package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const commandName = "repeat-forever"

func main() {
	var timeout time.Duration
	var repeatEvery time.Duration
	var quiet bool

	app := cli.NewApp()
	app.Name = commandName
	app.Usage = "Repeatedly run a command"
	app.Flags = []cli.Flag{
		cli.DurationFlag{
			Name:        "e, every,",
			Destination: &repeatEvery,
			Required:    true,
		},
		cli.DurationFlag{
			Name:        "t, timeout",
			Destination: &timeout,
		},
		cli.BoolFlag{
			Name:        "quiet, q",
			Destination: &quiet,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			return errors.New("Expected command")
		}
		command := c.Args()

		log.SetPrefix(fmt.Sprintf("[%s] ", commandName))
		log.SetFlags(log.Ltime)
		if quiet {
			log.SetOutput(ioutil.Discard)
		}

		repeatCommandForever(command, repeatEvery, timeout)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func repeatCommandForever(command []string, repeatEvery, timeout time.Duration) {
	for {
		log.Printf("Running %#v\n", strings.Join(command, " "))

		startedAt := time.Now()

		timedOut := func() bool {
			ctx := context.Background()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}

			cmd := exec.CommandContext(ctx, command[0], command[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()

			return ctx.Err() == context.DeadlineExceeded
		}()

		finishedAt := time.Now()

		elapsed := finishedAt.Sub(startedAt)
		elapsedStr := fmt.Sprintf("%.1f", float64(elapsed)/float64(time.Second))

		if timedOut {
			log.Printf("Command timed out after %s seconds\n", elapsedStr)
		} else {
			log.Printf("Command took %s seconds\n", elapsedStr)
		}

		sleepFor := repeatEvery - elapsed
		if sleepFor > 0 {
			log.Printf("Waiting %.1f seconds\n", float64(sleepFor)/float64(time.Second))
			time.Sleep(sleepFor)
		}
	}
}
