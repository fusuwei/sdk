package utils

import (
	"os"
	"path"
)

func MakeDir(elem ...string) (p string, ok bool) {
	if len(elem) < 1 {
		return "", true
	}
	ok = true
	p = path.Join(elem...)
	_, err := os.Stat(p)
	if err == nil {
		return
	}
	if os.IsExist(err) {
		return
	}
	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		ok = false
		return
	}
	return
}
