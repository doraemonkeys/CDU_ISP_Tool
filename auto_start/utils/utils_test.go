package utils

import "testing"

func TestReadStartWithLastLine(t *testing.T) {
	type args struct {
		filename string
		n        int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{filename: "test.txt", n: 1}, want: "", wantErr: false},
		{name: "2", args: args{filename: "test.txt", n: 2}, want: "555", wantErr: false},
		{name: "3", args: args{filename: "test.txt", n: 3}, want: "444", wantErr: false},
		{name: "4", args: args{filename: "test.txt", n: 4}, want: "333", wantErr: false},
		{name: "5", args: args{filename: "test.txt", n: 5}, want: "222", wantErr: false},
		{name: "6", args: args{filename: "test.txt", n: 6}, want: "111", wantErr: false},
		{name: "7", args: args{filename: "test.txt", n: 7}, want: "", wantErr: false},
		{name: "8", args: args{filename: "test.txt", n: 8}, want: "", wantErr: true},
		{name: "9", args: args{filename: "test.txt", n: 9}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadStartWithLastLine(tt.args.filename, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStartWithLastLine() error = %v, wantErr:%v\n", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadStartWithLastLine() = %v, want:%v\n", got, tt.want)
			}
		})
	}
}
