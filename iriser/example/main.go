package main

import (
	"errors"

	"github.com/gxxgle/go-utils/iriser"
	"github.com/gxxgle/go-utils/log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

type defaultHandler struct{}

func (*defaultHandler) HelloPanic(ctx iris.Context) {
	_, _ = ctx.JSON(iris.Map{
		"name": "world",
	})
	panic(errors.New("custom panic"))
}

func (*defaultHandler) Hello(ctx iris.Context) {
	_, _ = ctx.JSON(iris.Map{
		"name": "world",
	})
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(iriser.NewLogger(nil))

	iriser.Register(app.Party("/v1"), new(defaultHandler))

	log.LogIfError(app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed)))
}
