// License: GPLv3 Copyright: 2023, Kovid Goyal, <kovid at kovidgoyal.net>

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

var _ = fmt.Print

var DefaultExeSearchPaths = (&Once[[]string]{Run: func() []string {
	candidates := [...]string{"/usr/local/bin", "/opt/bin", "/opt/homebrew/bin", "/usr/bin", "/bin", "/usr/sbin", "/sbin"}
	ans := make([]string, 0, len(candidates))
	for _, x := range candidates {
		if s, err := os.Stat(x); err == nil && s.IsDir() {
			ans = append(ans, x)
		}
	}
	return ans
}}).Get

func Which(cmd string, paths ...string) string {
	if strings.Contains(cmd, string(os.PathSeparator)) {
		return ""
	}
	if len(paths) == 0 {
		path := os.Getenv("PATH")
		if path == "" {
			return ""
		}
		paths = strings.Split(path, string(os.PathListSeparator))
	}
	for _, dir := range paths {
		q := filepath.Join(dir, cmd)
		if unix.Access(q, unix.X_OK) == nil {
			s, err := os.Stat(q)
			if err == nil && !s.IsDir() {
				return q
			}
		}

	}
	return ""
}
