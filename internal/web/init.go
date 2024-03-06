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
	lad := NewLayoutJSONAdapter(c, rt)
	go lad.Handle(ctx)

	rt2 := NewRouter()
	lad2 := NewLayoutHTTP(c, rt2)
	go lad2.Handle(ctx)

	e.POST("/switchSet", lad2.SetSwitch)
	e.GET("/chunk", RouteWithMessageRouter(rt))
	e.GET("/chunkHTML", RouteWithMessageRouter(rt2))
	// Start the server

	return e.Start(":9898")
}
