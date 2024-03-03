package web

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanderv/traincontrol-prototype/internal/traintracks"
)

func Init(ctx context.Context, c *traintracks.TrackService) error {
	// Add file server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Static("internal/web/static"))

	// Add route for getting chunked data
	rt := NewRouter()
	lad := NewLayoutAdapter(c, rt)
	go lad.Handle(ctx)

	e.POST("/switchSet", lad.SetSwitch)
	e.GET("/chunk", RouteWithMessageRouter(rt))
	// Start the server

	return e.Start(":9898")
}
