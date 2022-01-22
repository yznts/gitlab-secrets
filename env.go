package main

import (
	"sort"
	"strings"
)

func envencode(kvmap map[string]string, novalues bool) string {
	// Output builder
	output := strings.Builder{}
	// Make ordered list of keys
	keys := make([]string, 0, len(kvmap))
	for k := range kvmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Iterate over keys
	for i, key := range keys {
		// Extact value
		value := kvmap[key]
		// Write key
		output.WriteString(key)
		// Write value
		if !novalues {
			output.WriteString("=")
			output.WriteString(`"` + value + `"`)
		}
		// Write newline
		if i != len(keys)-1 {
			output.WriteString("\n")
		}
	}
	// Return output
	return output.String()
}

func envdecode(input string) map[string]string {
	// Output map
	output := map[string]string{}
	// Split input into lines
	lines := strings.Split(input, "\n")
	// Iterate over lines
	for _, line := range lines {
		// Split line into key-value pair
		pair := strings.Split(line, "=")
		// Extract
		key := pair[0]
		value := pair[1]
		// Trim quotes, if present
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = value[1 : len(value)-1]
		}
		// Add key-value pair to output
		output[key] = value
	}
	// Return output
	return output
}
