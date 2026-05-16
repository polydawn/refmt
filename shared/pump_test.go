package shared

import (
	"errors"
	"strings"
	"testing"

	"github.com/polydawn/refmt/tok"
)

type scriptedSource struct {
	tokens []tok.Token
	errAt  int
	i      int
}

func (s *scriptedSource) Step(fillme *tok.Token) (bool, error) {
	if s.errAt == s.i {
		return true, errors.New("source error")
	}
	*fillme = s.tokens[s.i]
	done := s.i == len(s.tokens)-1
	s.i++
	return done, nil
}

type recordingSink struct {
	doneAt int
	errAt  int
	i      int
	got    []tok.Token
}

func (s *recordingSink) Step(token *tok.Token) (bool, error) {
	if s.errAt == s.i {
		return true, errors.New("sink error")
	}
	s.got = append(s.got, *token)
	done := s.i == s.doneAt
	s.i++
	return done, nil
}

func TestTokenPump(t *testing.T) {
	tokens := []tok.Token{
		{Type: tok.TMapOpen, Length: -1},
		{Type: tok.TString, Str: "k"},
		{Type: tok.TInt, Int: 1},
		{Type: tok.TMapClose},
	}

	tests := []struct {
		name       string
		source     *scriptedSource
		sink       *recordingSink
		wantErr    string
		wantTokens int
	}{
		{
			name:       "success",
			source:     &scriptedSource{tokens: tokens, errAt: -1},
			sink:       &recordingSink{doneAt: len(tokens) - 1, errAt: -1},
			wantTokens: len(tokens),
		},
		{
			name:       "source error stops before sink",
			source:     &scriptedSource{tokens: tokens, errAt: 2},
			sink:       &recordingSink{doneAt: len(tokens) - 1, errAt: -1},
			wantErr:    "source error",
			wantTokens: 2,
		},
		{
			name:       "sink error is returned",
			source:     &scriptedSource{tokens: tokens, errAt: -1},
			sink:       &recordingSink{doneAt: len(tokens) - 1, errAt: 2},
			wantErr:    "sink error",
			wantTokens: 2,
		},
		{
			name:       "source done before sink",
			source:     &scriptedSource{tokens: tokens, errAt: -1},
			sink:       &recordingSink{doneAt: len(tokens), errAt: -1},
			wantErr:    "src at end of item",
			wantTokens: len(tokens),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TokenPump{TokenSource: tt.source, TokenSink: tt.sink}.Run()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Run returned error: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Run error = %v, want containing %q", err, tt.wantErr)
			}
			if len(tt.sink.got) != tt.wantTokens {
				t.Fatalf("sink saw %d tokens, want %d", len(tt.sink.got), tt.wantTokens)
			}
			for i, token := range tt.sink.got {
				if !tok.IsTokenEqual(token, tokens[i]) {
					t.Fatalf("sink token %d = %s, want %s", i, token, tokens[i])
				}
			}
		})
	}
}
