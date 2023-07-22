package auth

import "testing"

func Test_getExpireTime(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			"test",
			1800000000000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getExpireTime(); got != tt.want {
				t.Errorf("getExpireTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
