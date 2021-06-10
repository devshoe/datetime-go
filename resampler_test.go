package datetime

import (
	"reflect"
	"testing"
	"time"
)

var times = []string{"2020-12-12 09:15:00", "2020-12-12 09:17:00", "2020-12-12 09:20:00", "2020-12-12 09:37:00", "2020-12-12 09:39:00", "2020-12-13 09:00:00"}

func TestBucketDataByInterval(t *testing.T) {
	parsedTimes, _ := ParseDatetimeArray(times)
	type args struct {
		index             []time.Time
		bucketingInterval time.Duration
		startTime         []time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    [][]time.Time
		wantErr bool
	}{
		{"no offset 5 mins", args{parsedTimes, time.Minute * 5, nil}, [][]time.Time{parsedTimes[:2], {parsedTimes[2]}, parsedTimes[3:5], {parsedTimes[5]}}, false},
		{"with offset -3, 5 mins", args{parsedTimes, time.Minute * 5, []time.Time{parsedTimes[0].Add(-time.Duration(time.Minute * 3))}}, [][]time.Time{{parsedTimes[0]}, parsedTimes[1:3], parsedTimes[3:5], {parsedTimes[5]}}, false},
		{"no offset 1 day", args{parsedTimes, time.Hour * 24, nil}, [][]time.Time{parsedTimes}, false},
		{"with offset day start, 1 day", args{parsedTimes, time.Hour * 24, []time.Time{ExtractDateFromDatetime(parsedTimes[0])}}, [][]time.Time{parsedTimes[:5], {parsedTimes[5]}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BucketTimeArrayByInterval(tt.args.index, tt.args.bucketingInterval, tt.args.startTime...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BucketDataByInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BucketDataByInterval() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBucketDataSliceByBucketedTimeArray(t *testing.T) {
	parsedTimes, _ := ParseDatetimeArray(times)
	type args struct {
		bucketedTimes [][]time.Time
		data          interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{"floatTest", args{[][]time.Time{parsedTimes[:2], {parsedTimes[2]}, parsedTimes[3:5], {parsedTimes[5]}}, []float64{0, 1, 2, 3, 4, 5}}, [][]float64{{0, 1}, {2}, {3, 4}, {5}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BucketDataSliceByBucketedTimeArray(tt.args.bucketedTimes, tt.args.data)
			got = got.([][]float64)
			if (err != nil) != tt.wantErr {
				t.Errorf("BucketDataSliceByBucketedTimeArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BucketDataSliceByBucketedTimeArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}
