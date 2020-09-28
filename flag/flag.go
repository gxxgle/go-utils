package flag

import (
	"os"
	"path/filepath"

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

func Add(flag cli.Flag, fn cli.ActionFunc) {
	app.Flags = append(app.Flags, flag)
	if fn != nil {
		actions = append(actions, fn)
	}
}

func Run() error {
	app.Flags = append(app.Flags, cli.HelpFlag)
	return app.Run(os.Args)
}
