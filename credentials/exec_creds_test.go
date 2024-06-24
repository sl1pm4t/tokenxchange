package credentials

import (
	"testing"
	"time"
)

func TestFormatExecCredential(t *testing.T) {
	now := time.Now()
	type args struct {
		token      string
		expiration time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJzdWIiOiJzdWJqZWN0In0.OFD0iVfPczqWBA_TRi1jGB5PF699eekcHt4D6qNoimc",
				expiration: now,
			},
			want: FormatExecCredential("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpc3N1ZXIiLCJzdWIiOiJzdWJqZWN0In0.OFD0iVfPczqWBA_TRi1jGB5PF699eekcHt4D6qNoimc", now),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatExecCredential(tt.args.token, tt.args.expiration); got != tt.want {
				t.Errorf("FormatExecCredential() = %v, want %v", got, tt.want)
			}
		})
	}
}
