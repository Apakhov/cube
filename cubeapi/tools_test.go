package cubeapi_test

import "testing"

func printDiffInfo(t *testing.T, expected, got interface{}, addInfo ...string) {
	info := ""
	for _, s := range addInfo {
		info += s + "\n"
	}
	t.Errorf(`
%s
expected:
	%+v
got:
	%+v`, info, expected, got)
}
