package flag

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/urfave/cli/v2"
)

var (
	app     = cli.NewApp()
	actions []cli.ActionFunc
)

func init() {
	app.Name = filepath.Base(os.Args[0])
	app.HideHelp = true
	app.Action = action
}

func action(ctx *cli.Context) error {
	if ctx.Bool("help") {
		cli.ShowAppHelpAndExit(ctx, 0)
	}

	for _, fn := range actions {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}

func Add(flag cli.Flag, acts ...cli.ActionFunc) {
	app.Flags = append(app.Flags, flag)
	actions = append(actions, acts...)
}

func Run() error {
	app.Flags = append(app.Flags, cli.HelpFlag)
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	return app.Run(os.Args)
}
