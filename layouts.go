package datetime

/*
kinda useless now
*/
import (
	"fmt"
	"strings"
	"time"
)

var layoutCollection = []string{"2006-01-02T15:04:05Z", "2006-01-02", "2006/01/02", "15:04:05",
	"2006-01-02 15:04:05", "2006-01-02T15:04:05Z0700", "2006-01-02T15:04:05Z07:00", time.RFC3339, time.RFC3339Nano,
}

var layoutConversionMap = map[string]string{
	"YYYY": "2006",
	"YY":   "06",
	"MM":   "01",
	"DD":   "02",
	"hh":   "15",
	"mm":   "04",
	"ss":   "05",
	"zh":   "07",
	"zm":   "00",
	"nn":   "999999999",
}

var layoutLengthMap = map[int][]string{}

func init() {
	//for efficiency, create a map which maps length of layouts to a list of layouts of that length
	for _, layout := range layoutCollection {
		layoutLengthMap[len(layout)] = append(layoutLengthMap[len(layout)], layout)
	}
}

//ConvertLayoutToGolangLayout converts a YYYY MM DD hh mm ss
func convertLayoutToGolangLayout(inputLayout string) string {
	for k, v := range layoutConversionMap {
		inputLayout = strings.Replace(inputLayout, k, v, 1)
	}
	return inputLayout
}

func smartDetectLayout(input string) (string, error) {
	_, exists := layoutLengthMap[len(input)]
	if !exists {
		return "", fmt.Errorf("(smartDetectLayout)failed as length of input string is not in collection")
	}
	for _, layout := range layoutLengthMap[len(input)] {
		_, err := time.Parse(layout, input)
		if err == nil {
			return layout, nil
		}
	}
	return "", fmt.Errorf("(smartDetectLayout)failed to detect valid layout for given string")
}

//SmartParse automatically detects input format and returns corresponding time, will return error if not found
func smartParse(input string) (time.Time, error) {
	layout, err := smartDetectLayout(input)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(layout, input)
}
