package main

import (
	"os"

	"github.com/ggiamarchi/linuxboot-ci-api/api"
	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("linuxboot-ci-api", "Linuxboot CI API")

	app.Command("run", "Run API server", func(cmd *cli.Cmd) {

		port := cmd.IntOpt("p port", 1234, "TCP port to bind on")

		cmd.Action = func() {
			api.Run(*port)
		}
	})

	app.Run(os.Args)
}
