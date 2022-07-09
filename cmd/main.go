package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/t-margheim/boggler/internal/boggle"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	rawLog, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		e.Logger.Panicf("failed to create zapLogger: %s", err.Error())
	}
	sugarLog := rawLog.Sugar()

	s := boggle.NewService(nil, sugarLog)
	h := boggle.NewHandler(s, sugarLog)
	e.GET("/board/:letters/solve", h.GetWords)

	e.Logger.Print(e.Start(":80"))
}
