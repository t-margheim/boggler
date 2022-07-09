package boggle

import (
	"fmt"
	"strings"

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
	resultsMap := make(map[string]struct{})
	// for i, startCharacter := range board {
	for i := range board {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)

		// s.solveStartPosition(i, boardCopy, toString(startCharacter), resultsMap)
		s.solveStartPosition(i, boardCopy, "", resultsMap)
	}

	var words []string
	for k := range resultsMap {
		words = append(words, k)
	}

	return words, nil
}

func (s *service) validateBoard(board []rune) error {
	gotChars := len(board)
	if gotChars != boardSize {
		return fmt.Errorf("invalid number of characters, expected %d, got %d",
			boardSize,
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

func (s *service) solveStartPosition(pos int, board []rune, current string, results map[string]struct{}) {
	s.log.Debugw("start function",
		"board", board,
		"position", pos,
		"current_word", current,
	)

	newWord := strings.Builder{}
	newWord.WriteString(current)
	newWord.WriteRune(board[pos])

	matches := s.tr.PrefixMatch(current)
	if matches == nil {
		return
	}

	for _, w := range matches {
		if current == w {
			results[current] = struct{}{}
			break
		}
	}

	board[pos] = '_'

	// move right
	rightIdx := pos + 1
	if isSameRow(pos, rightIdx) {
		s.solveStartPosition(rightIdx, board, newWord.String(), results)
	}

	// move left
	leftIdx := pos - 1
	if leftIdx >= 0 && isSameRow(pos, leftIdx) {
		s.solveStartPosition(leftIdx, board, newWord.String(), results)
	}

	// move down
	downIdx := pos + numCols
	if downIdx < boardSize {
		s.solveStartPosition(downIdx, board, newWord.String(), results)
	}

	// move up
	upIdx := pos - numCols
	if upIdx >= 0 {
		s.solveStartPosition(upIdx, board, newWord.String(), results)
	}
}

func isSameRow(pos1, pos2 int) bool {
	return pos1/numRows == pos2/numRows
}

// toString converts runes to string and also handles special case
// of converting to 'q' -> "qu"
func toString(l rune) string {
	if l == 'q' {
		return "qu"
	}
	return string(l)
}
