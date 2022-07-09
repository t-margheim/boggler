package boggle

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type mockService struct {
	solveBoardCalledTimes int
	solveBoardCalledWith  []rune
	mockSolveBoard        func() []string

	validateBoardCalledTimes int
	validateBoardCalledWith  []rune
	mockValidateBoard        func() error
}

func (m *mockService) solveBoard(board []rune) []string {
	m.solveBoardCalledTimes++
	m.solveBoardCalledWith = board
	return m.mockSolveBoard()
}

func (m *mockService) validateBoard(board []rune) error {
	m.validateBoardCalledTimes++
	m.validateBoardCalledWith = board
	return m.mockValidateBoard()
}

func Test_handler_GetWords(t *testing.T) {
	testBoard := "testboardinput"
	tests := []struct {
		name                    string
		svc                     *mockService
		board                   string
		wantRespStatus          int
		wantRespBody            string
		wantLogs                []observer.LoggedEntry
		wantValidateCalledTimes int
		wantValidateCalledWith  []rune
		wantSolveCalledTimes    int
		wantSolveCalledWith     []rune
	}{
		{
			name: "validation fails",
			svc: &mockService{
				mockValidateBoard: func() error {
					return errors.New("unit test error")
				},
			},
			board:          testBoard,
			wantRespStatus: http.StatusBadRequest,
			wantRespBody:   `{"err":"unit test error","msg":"invalid board submitted"}`,
			wantLogs: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zap.WarnLevel,
						Message: "invalid board submitted",
					},
					Context: []zapcore.Field{
						zap.String("board", testBoard),
						zap.String("error", "unit test error"),
					},
				},
			},
			wantValidateCalledTimes: 1,
			wantValidateCalledWith:  []rune(testBoard),
			wantSolveCalledTimes:    0,
			wantSolveCalledWith:     nil,
		},
		{
			name: "success",
			svc: &mockService{
				mockSolveBoard: func() []string {
					return []string{"def", "abc"}
				},
				mockValidateBoard: func() error {
					return nil
				},
			},
			board:                   testBoard,
			wantRespStatus:          http.StatusOK,
			wantRespBody:            `{"words":["abc","def"]}`,
			wantLogs:                []observer.LoggedEntry{},
			wantValidateCalledTimes: 1,
			wantValidateCalledWith:  []rune(testBoard),
			wantSolveCalledTimes:    1,
			wantSolveCalledWith:     []rune(testBoard),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obsCore, obsLogs := observer.New(zap.InfoLevel)
			obsLogger := zap.New(obsCore)
			h := &handler{
				svc: tt.svc,
				log: obsLogger.Sugar(),
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/testinput", nil)
			c := echo.New().NewContext(r, w)

			// This is normally handled by the echo router
			c.SetParamNames("letters")
			c.SetParamValues(testBoard)

			err := h.GetWords(c)
			assert.NoErrorf(t, err, "got unexpected error")

			assert.Equalf(t, tt.wantRespStatus, w.Code, "unexpected status code")
			respBody := w.Body.String()
			respBody = strings.TrimSuffix(respBody, "\n")
			assert.Equalf(t, tt.wantRespBody, respBody, "unexpected response body")

			// service spies
			assert.Equalf(t, tt.wantSolveCalledTimes, tt.svc.solveBoardCalledTimes,
				"solveBoard called unexpected number of times")
			assert.Equalf(t, tt.wantSolveCalledWith, tt.svc.solveBoardCalledWith,
				"solveBoard called with unexpected input")
			assert.Equalf(t, tt.wantValidateCalledTimes, tt.svc.validateBoardCalledTimes,
				"validateBoard called unexpected number of times")
			assert.Equalf(t, tt.wantValidateCalledWith, tt.svc.validateBoardCalledWith,
				"validateBoard called with unexpected input")

			// log spies
			assert.Equalf(t, tt.wantLogs, obsLogs.AllUntimed(), "unexpected logs")
		})
	}
}
