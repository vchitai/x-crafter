package main

import (
	"fmt"
	"strconv"
	"strings"
)

type syntaxVault []string

func (sv syntaxVault) getKeyByIdx(idx int) string {
	return fmt.Sprintf("@@@syntax_%d@@@", idx)
}
func (sv syntaxVault) lastInsertedKey() string {
	return sv.getKeyByIdx(len(sv) - 1)
}
func (sv syntaxVault) put(s string) syntaxVault {
	return append(sv, s)
}
func (sv syntaxVault) get(s string) string {
	s = strings.Replace(s, "@@@syntax_", "", 1)
	s = strings.Replace(s, "@@@", "", 1)
	idx, _ := strconv.Atoi(s)
	return sv[idx-1]
}
