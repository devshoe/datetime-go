package datetime

import (
	"fmt"
	"reflect"
	"time"
)

//BucketTimeArrayByInterval buckets time array by interval provided, ie it starts at first sample and splits into arrays each of the duration supplied
//if startTime is provided, that is used as reference to start splitting
//this is useful if you have a struct which is indexed by time, splitting data columns by length of arrays returned
func BucketTimeArrayByInterval(index []time.Time, bucketingInterval time.Duration, startTime ...time.Time) ([][]time.Time, error) {
	var startAtTime time.Time
	var bucketedTimes [][]time.Time

	if len(index) == 0 {
		return bucketedTimes, fmt.Errorf("(BucketArrayByInterval) cannot proceed as length of time array is 0")
	}
	if startTime != nil {
		startAtTime = startTime[0]
	} else {
		startAtTime = index[0]
	}

	nextSplitAt := startAtTime.Add(bucketingInterval)
	presentBucket := []time.Time{}
	for i := 0; i < len(index); i++ {
		for index[i].After(nextSplitAt) || index[i].Equal(nextSplitAt) {
			if len(presentBucket) != 0 {
				bucketedTimes = append(bucketedTimes, presentBucket)
				presentBucket = []time.Time{}
			}
			nextSplitAt = nextSplitAt.Add(bucketingInterval)
		}
		if index[i].Before(nextSplitAt) {
			presentBucket = append(presentBucket, index[i])
		}
	}
	if len(presentBucket) != 0 {
		bucketedTimes = append(bucketedTimes, presentBucket)
	}
	return bucketedTimes, nil
}

//BucketDataSliceByBucketedTimeArray converts data of any array type to multi row array each of lengths matching bucketedTimes
//error if total lengths dont match
//After recieving output, convert to to type using  output.([][]Typename)
func BucketDataSliceByBucketedTimeArray(bucketedTimes [][]time.Time, data interface{}) (interface{}, error) {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	if dataType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("(BucketDataSliceByBucketedTimeArray) failed because data is not a slice type")
	}
	totalLen := 0
	bucketedData := reflect.MakeSlice(reflect.SliceOf(dataType), 0, len(bucketedTimes))
	for _, bucket := range bucketedTimes {
		bucketedData = reflect.Append(bucketedData, dataValue.Slice(totalLen, totalLen+len(bucket)))
		totalLen += len(bucket)
	}
	if totalLen != dataValue.Len() {
		return bucketedData.Interface(), fmt.Errorf("(BucketFloat64DataByBucketedTimeArray)failed because of length mismatch")
	}
	return bucketedData, nil
}

//BucketFloat64SliceByBucketedTimeArray converts data to multi row array corresponding to lengths of bucketedTimes provided
func BucketFloat64SliceByBucketedTimeArray(bucketedTimes [][]time.Time, data []float64) ([][]float64, error) {
	totalLen := 0
	bucketedData := [][]float64{}
	for _, bucket := range bucketedTimes {
		bucketedData = append(bucketedData, data[totalLen:totalLen+len(bucket)])
		totalLen += len(bucket)
	}
	if totalLen != len(data) {
		return bucketedData, fmt.Errorf("(BucketFloat64DataByBucketedTimeArray)failed because of length mismatch")
	}
	return bucketedData, nil
}

//GenerateTimeRange creates a range of times by adding interval to startTime `length` number of times
func GenerateTimeRange(startTime time.Time, interval time.Duration, length int) []time.Time {
	t := []time.Time{}
	for i := 0; i < length; i++ {
		t = append(t, startTime)
		startTime = startTime.Add(interval)
	}
	return t
}

//GenerateTimeRangeBetween creates a range of times of size interval between start and end times
func GenerateTimeRangeBetween(startTime time.Time, endTime time.Time, interval time.Duration) []time.Time {
	t := []time.Time{}
	for endTime.After(startTime) {
		t = append(t, startTime)
		startTime = startTime.Add(interval)
	}
	return t
}
