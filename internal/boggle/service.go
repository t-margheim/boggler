package boggle

import (
	"errors"
	"fmt"

	"github.com/minio/pkg/trie"
	"go.uber.org/zap"
)

type Service interface {
	solveBoard(board []rune) (wordList []string, err error)
	validateBoard(board []rune) error
}

func NewService(wordTrie *trie.Trie, l *zap.SugaredLogger) Service {
	return &service{
		tr:  wordTrie,
		log: l,
	}
}

type service struct {
	tr *trie.Trie

	log *zap.SugaredLogger
}

func (s *service) solveBoard(board []rune) ([]string, error) {

	return nil, errors.New("not yet implemented")
}

func (s *service) validateBoard(board []rune) error {
	wantChars := numCols * numRows
	gotChars := len(board)
	if gotChars != wantChars {
		return fmt.Errorf("invalid number of characters, expected %d, got %d",
			wantChars,
			gotChars,
		)
	}

	var invalidChars []rune
	for _, l := range board {
		if l < validCharStart || l > validCharEnd {
			invalidChars = append(invalidChars, l)
		}
	}
	if len(invalidChars) > 0 {
		return fmt.Errorf("invalid characters submitted: %s", string(invalidChars))
	}

	return nil
}
