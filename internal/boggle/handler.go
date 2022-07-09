package boggle

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler interface {
	GetWords(c echo.Context) error
}

func NewHandler(s Service, l *zap.SugaredLogger) Handler {
	return &handler{
		svc: s,
		log: l,
	}
}

type handler struct {
	svc Service

	log *zap.SugaredLogger
}

// GetWords reads the query parameter passed in as the board and returns
// a list of words from the wordlist that appear in the board.
func (h *handler) GetWords(c echo.Context) error {
	input := c.Param("letters")
	board := []rune(input)

	err := h.svc.validateBoard(board)
	if err != nil {
		h.log.Warnw("invalid board submitted",
			"board", board,
			"error", err.Error())
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("board is invalid: %s", err.Error()),
		)
	}

	ww, err := h.svc.solveBoard(board)
	if err != nil {
		h.log.Errorw("failed to process board",
			"board", board,
			"error", err.Error())
		return echo.NewHTTPError(
			http.StatusInternalServerError,
		)
	}

	sort.Strings(ww)
	return c.JSON(http.StatusOK, response{Words: ww})
}
