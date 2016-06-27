package main

import (
	"net"
	"net/http"

	"github.com/steemwatch/status/checks"
	"github.com/steemwatch/status/views"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

func startServer(listener net.Listener, runner *checks.Runner) {
	// New Echo instance.
	e := echo.New()

	// Echo middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.AddTrailingSlash())

	// Template rendering.
	renderer, err := views.NewRenderer("./views/*.html")
	if err != nil {
		panic(err)
	}
	e.SetRenderer(renderer)

	// Index.
	e.GET("/", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "home.html", &views.Context{
			Sections: runner.Results(),
		})
	})

	// Run.
	e.Run(fasthttp.WithConfig(engine.Config{
		Listener: listener,
	}))
}
