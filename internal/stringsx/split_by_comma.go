package stringsx

import "strings"

func SplitByComma(s string) []string {
	// When s is empty, strings.Split returns [""], which is neither intuitive
	// nor desirable. This if statement ensures that an empty slice is returned
	// instead.
	if len(s) == 0 {
		return []string{}
	}
	return strings.Split(s, ", ")
}
