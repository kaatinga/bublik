package bublyk

import "testing"

func TestNewDateFromTime(t *testing.T) {
	type args struct {
		t *time.Time
	}
	tests := []struct {
		name string
		args args
		want Date
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDateFromTime(tt.args.t); got != tt.want {
				t.Errorf("NewDateFromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
