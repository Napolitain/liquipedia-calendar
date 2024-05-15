package main

import (
	"context"
	"testing"
)

func Test_isValidGame(t *testing.T) {
	type args struct {
		ctx  context.Context
		game string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidGame(tt.args.ctx, tt.args.game); got != tt.want {
				t.Errorf("isValidGame() = %v, want %v", got, tt.want)
			}
		})
	}
}
