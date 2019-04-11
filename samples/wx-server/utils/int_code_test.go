package utils

import (
	"testing"
)

func Test_int2code(t *testing.T) {
	i1 := uint64(1234678)
	s, err := Int2Code(i1)
	if err != nil {
		t.Fatalf("converting %d error: %v\n", i1, err)
	}
	t.Logf("%d => %s\n", i1, s)

	i2, err := Code2Int(s)
	if err != nil {
		t.Fatalf("failed: %v\n", err)
	}
	t.Logf("original int: %d\n", i2)
}
