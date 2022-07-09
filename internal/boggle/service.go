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
		"board", string(board),
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
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(rightIdx, boardCopy, newWord.String(), results)
	}

	// move right down
	rightDownIdx := pos + 1 + numCols
	if rightDownIdx < boardSize && isSameRow(pos+numCols, rightDownIdx) {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(rightDownIdx, boardCopy, newWord.String(), results)
	}

	// move down
	downIdx := pos + numCols
	if downIdx < boardSize {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(downIdx, boardCopy, newWord.String(), results)
	}

	// move left down
	leftDownIdx := pos - 1 + numCols
	if leftDownIdx < boardSize && isSameRow(pos+numCols, leftDownIdx) {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(leftDownIdx, boardCopy, newWord.String(), results)
	}

	// move left
	leftIdx := pos - 1
	if leftIdx >= 0 && isSameRow(pos, leftIdx) {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(leftIdx, boardCopy, newWord.String(), results)
	}

	// move left up
	leftUpIdx := pos - 1 - numCols
	if leftUpIdx >= 0 && isSameRow(pos-numCols, leftUpIdx) {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(leftUpIdx, boardCopy, newWord.String(), results)
	}

	// move up
	upIdx := pos - numCols
	if upIdx >= 0 {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(upIdx, boardCopy, newWord.String(), results)
	}

	// move right up
	rightUpIdx := pos + 1 - numCols
	if rightUpIdx >= 0 && isSameRow(pos-numCols, rightUpIdx) {
		boardCopy := make([]rune, len(board))
		copy(boardCopy, board)
		s.solveStartPosition(rightUpIdx, boardCopy, newWord.String(), results)
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

func copyBoard(board []rune) []rune {
	boardCopy := make([]rune, len(board))
	copy(boardCopy, board)
	return boardCopy
}
