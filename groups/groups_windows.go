//go:build !linux

package groups

func IsUserIn(groups ...string) bool {
	return true
}
