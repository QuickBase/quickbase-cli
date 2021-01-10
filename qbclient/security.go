package qbclient

import "regexp"

var reUserTokenMask *regexp.Regexp

// MaskUserToken masks user tokens in a byte slice.
func MaskUserToken(b []byte) []byte {
	return reUserTokenMask.ReplaceAll(b, []byte(`${1}_${2}********************${3}`))
}

// MaskUserTokenString masks user tokens in a string.
func MaskUserTokenString(s string) string {
	return reUserTokenMask.ReplaceAllString(s, "${1}_${2}********************${3}")
}

func init() {
	reUserTokenMask = regexp.MustCompile(`([0-9a-z]+_[0-9a-z]+)_([0-9a-z]{4})[0-9a-z]+([0-9a-z]{4})`)
}
