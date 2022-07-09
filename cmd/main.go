package main

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/pkg/trie"
	"github.com/t-margheim/boggler/internal/boggle"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	rawLog, err := zap.NewProduction()
	if err != nil {
		e.Logger.Panicf("failed to create zapLogger: %s", err.Error())
	}
	sugarLog := rawLog.Sugar()

	start := time.Now()
	tr := generateTrie()
	sugarLog.Infow("trie generated",
		"duration", time.Since(start))

	s := boggle.NewService(tr, sugarLog)
	h := boggle.NewHandler(s, sugarLog)
	e.GET("/board/:letters/solve", h.GetWords)

	e.Logger.Print(e.Start(":80"))
}

func generateTrie() *trie.Trie {
	tr := trie.NewTrie()

	for _, word := range words {
		tr.Insert(word)
	}

	return tr
}
