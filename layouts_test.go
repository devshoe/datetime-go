package datetime

import "testing"

func Test_smartDetectLayout(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", args{"2020-12-12"}, "2006-01-02", false},
		{"", args{""}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := smartDetectLayout(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("smartDetectLayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("smartDetectLayout() = %v, want %v", got, tt.want)
			}
		})
	}
}
