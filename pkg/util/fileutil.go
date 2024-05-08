package util

import (
	"strings"
)

// PortablePosixValidChar is the list of valid characters allowed for portability
const PortablePosixValidChar = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789._-"

// IsPortableFilename is a function that determines if a filename is using portable posix characters.
func IsPortableFilename(filename string) bool {
	if filename == "" {
		return false
	}

	// '.' and '..' are not filenames (but they are valid characters)
	// also want to disallow hidden files
	if strings.HasPrefix(filename, ".") {
		return false
	}

	for _, ch := range filename {
		idx := strings.IndexRune(PortablePosixValidChar, ch)
		if idx == -1 {
			return false
		}
	}

	return true
}

// IsPortablePath looks at each component of a path, and verifies that each component uses portable
// characters only.  Also ensures that the separator is always a "/" (even on Windows).
// we only want relative paths
func IsPortablePath(path string) bool {
	if strings.HasPrefix(path, "/") {
		return false
	}

	// We mean "/" here even on Windows.  All part path names use "/"
	for _, p := range strings.Split(path, "/") {
		// when splitting path, if the first or last character is '/', then an empty string is part of the
		// paths to check. We handle that here by skipping that result
		if p == "" {
			continue
		}
		if !IsPortableFilename(p) {
			return false
		}
	}

	return true
}
