package boggle

import (
	"testing"

	"github.com/minio/pkg/trie"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_service_solveStartPosition(t *testing.T) {
	testTrie := trie.NewTrie()
	// using "alt" as the other permutations of these letters are not in the word list
	testTrie.Insert("alt")

	rawLog, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		t.Error(err)
		return
	}
	sugarLog := rawLog.Sugar()

	type args struct {
		pos     int
		board   []rune
		current string
		results map[string]struct{}
	}
	tests := []struct {
		name        string
		args        args
		wantResults map[string]struct{}
	}{
		// {
		// 	name: "no words",
		// 	args: args{
		// 		pos: 0,
		// 		board: []rune{
		// 			'x', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{},
		// },
		// {
		// 	name: "finds words written top-left to bottom-right",
		// 	args: args{
		// 		pos: 0,
		// 		board: []rune{
		// 			'a', 'x', 'x', 'x',
		// 			'x', 'l', 'x', 'x',
		// 			'x', 'x', 't', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{
		// 		"alt": {},
		// 	},
		// },
		// {
		// 	name: "finds words written top to bottom",
		// 	args: args{
		// 		pos: 0,
		// 		board: []rune{
		// 			'a', 'x', 'x', 'x',
		// 			'l', 'x', 'x', 'x',
		// 			't', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{
		// 		"alt": {},
		// 	},
		// },
		// {
		// 	name: "finds words written top-right to bottom-left",
		// 	args: args{
		// 		pos: 2,
		// 		board: []rune{
		// 			'x', 'x', 'a', 'x',
		// 			'x', 'l', 'x', 'x',
		// 			't', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{
		// 		"alt": {},
		// 	},
		// },
		// {
		// 	name: "finds words written right to left",
		// 	args: args{
		// 		pos: 2,
		// 		board: []rune{
		// 			't', 'l', 'a', 'x',
		// 			'x', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{
		// 		"alt": {},
		// 	},
		// },
		// {
		// 	name: "finds words written bottom-right to top-left",
		// 	args: args{
		// 		pos: 10,
		// 		board: []rune{
		// 			't', 'x', 'x', 'x',
		// 			'x', 'l', 'x', 'x',
		// 			'x', 'x', 'a', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{
		// 		"alt": {},
		// 	},
		// },
		// {
		// 	name: "finds words written bottom to top",
		// 	args: args{
		// 		pos: 8,
		// 		board: []rune{
		// 			't', 'x', 'x', 'x',
		// 			'l', 'x', 'x', 'x',
		// 			'a', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{
		// 		"alt": {},
		// 	},
		// },
		// {
		// 	name: "finds words written bottom-left to top-right",
		// 	args: args{
		// 		pos: 8,
		// 		board: []rune{
		// 			'x', 'x', 't', 'x',
		// 			'x', 'l', 'x', 'x',
		// 			'a', 'x', 'x', 'x',
		// 			'x', 'x', 'x', 'x',
		// 		},
		// 		current: "",
		// 		results: map[string]struct{}{},
		// 	},
		// 	wantResults: map[string]struct{}{
		// 		"alt": {},
		// 	},
		// },
		{
			name: "finds words written in a zig zag",
			args: args{
				pos: 0,
				board: []rune{
					'a', 'x', 'x', 'x',
					'l', 't', 'x', 'x',
					'x', 'x', 'x', 'x',
					'x', 'x', 'x', 'x',
				},
				current: "",
				results: map[string]struct{}{},
			},
			wantResults: map[string]struct{}{
				"alt": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				tr:  testTrie,
				log: sugarLog,
			}
			s.solveStartPosition(tt.args.pos, tt.args.board, tt.args.current, tt.args.results)
			assert.Equal(t, tt.wantResults, tt.args.results)
		})
	}
}
