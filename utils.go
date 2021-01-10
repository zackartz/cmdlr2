package cmdlr2

import (
	"sort"
	"strings"
)

// StringHasPrefix checks whether or not the string contains one of the given prefixes and returns the string without the prefix
func StringHasPrefix(str string, prefixes []string, ignoreCase bool) (bool, string) {
	for _, prefix := range prefixes {
		stringToCheck := str
		if ignoreCase {
			stringToCheck = strings.ToLower(stringToCheck)
			prefix = strings.ToLower(prefix)
		}
		if strings.HasPrefix(stringToCheck, prefix) {
			return true, string(str[len(prefix):])
		}
	}
	return false, str
}

func StringTrimPreSuffix(str string, preSuffix string) string {
	if !(strings.HasPrefix(str, preSuffix) && strings.HasSuffix(str, preSuffix)) {
		return str
	}
	return strings.TrimPrefix(strings.TrimSuffix(str, preSuffix), preSuffix)
}

func Equals(str1, str2 string, ignoreCase bool) bool {
	if !ignoreCase {
		return str1 == str2
	}
	return strings.ToLower(str1) == strings.ToLower(str2)
}

func StringArrayContains(array []string, str string, ignoreCase bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
	}
	for _, value := range array {
		if ignoreCase {
			value = strings.ToLower(value)
		}
		if value == str {
			return true
		}
	}
	return false
}

func BuildCheckPrefixes(command *Command) []string {
	toCheck := make([]string, len(command.Aliases)+1)
	toCheck[0] = command.Name
	for i, alias := range command.Aliases {
		toCheck[1+i] = alias
	}

	sort.Slice(toCheck, func(i, j int) bool {
		return len(toCheck[i]) > len(toCheck[j])
	})

	return toCheck
}
