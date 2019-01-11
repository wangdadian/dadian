package jsonconf

import (
	"testing"
)

func TestReadConf(t *testing.T) {
	byBuff, err := ReadConf("json_test.conf")
	if err != nil {
		t.Errorf("read conf file failed: %s\n", err.Error())
		return
	}
	s := string(byBuff)
	t.Logf(s)
}
