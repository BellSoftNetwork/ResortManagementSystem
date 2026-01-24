package dto

import (
	"fmt"
	"strings"
	"time"
)

// JSONTime은 JSON 직렬화시 날짜 형식을 처리하는 타입
type JSONTime struct {
	time.Time
}

// UnmarshalJSON은 다양한 날짜 형식을 파싱
func (jt *JSONTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "null" || str == "" {
		return nil
	}

	// 다양한 날짜 형식 시도
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05",       // RFC3339 without timezone
		"2006-01-02",                // Date only
		"2006-01-02 15:04:05",       // DateTime with space
	}

	var parseErr error
	for _, format := range formats {
		t, err := time.Parse(format, str)
		if err == nil {
			jt.Time = t
			return nil
		}
		parseErr = err
	}

	return fmt.Errorf("unable to parse time %q: %v", str, parseErr)
}

// MarshalJSON은 시간을 JSON으로 변환
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	if jt.IsZero() {
		return []byte("null"), nil
	}
	// 날짜만 포함된 경우
	if jt.Hour() == 0 && jt.Minute() == 0 && jt.Second() == 0 {
		return []byte(fmt.Sprintf(`"%s"`, jt.Format("2006-01-02"))), nil
	}
	// 시간까지 포함된 경우
	return []byte(fmt.Sprintf(`"%s"`, jt.Format("2006-01-02T15:04:05"))), nil
}
