package uri

import (
	"strconv"
	"time"
)

func GetStringParam(params map[string][]string, name string) string {
	values, exists := params[name]
	if !exists || len(values) == 0 || values[0] == "" {
		return ""
	}
	return values[0]
}

func GetStringArrayParam(params map[string][]string, name string) []string {
	values, exists := params[name]
	if !exists || len(values) == 0 {
		return []string{}
	}

	nonEmpty := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" {
			nonEmpty = append(nonEmpty, v)
		}
	}
	return nonEmpty
}

func GetIntParam(params map[string][]string, name string) int {
	strValue := GetStringParam(params, name)
	if strValue == "" {
		return 0
	}

	val, err := strconv.Atoi(strValue)
	if err != nil {
		return 0
	}
	return val
}

func GetTimeParam(params map[string][]string, name string, format string) time.Time {
	strValue := GetStringParam(params, name)
	if strValue == "" {
		return time.Time{}
	}

	date, err := time.Parse(format, strValue)
	if err != nil {
		return time.Time{}
	}
	return date
}
