package id_generator

import (
	"strings"
	"testing"
)

func TestRandomIdGenerator_Generate(t *testing.T) {
	generator := new(RandomIdGenerator)
	for i := 0; i < 100; i++ {
		res := generator.Generate()
		cuts := strings.Split(res, "-")
		if len(cuts) != 3 {
			t.Errorf("res:%s un expect format", res)
		}
	}
}

func TestRandomIdGenerator_getLastFieldSplitByDot(t *testing.T) {
	generator := new(RandomIdGenerator)

	var inputs = map[string]string{
		"a1.a2": "a2",
		"f1":    "f1",
	}
	for k, v := range inputs {
		if res := generator.getLastFieldSplitByDot(k); res != v {
			t.Errorf("input:%s res:%s != %s", k, res, v)
		}
	}
}

func TestRandomIdGenerator_generateRandomAlphameric(t *testing.T) {
	generator := new(RandomIdGenerator)
	for i := 0; i < 100; i++ {
		curRandStr := generator.generateRandomAlphameric(8)
		if len(curRandStr) != 8 {
			t.Errorf("cur:%s len != 8", curRandStr)
		}
		for _, c := range curRandStr {
			isDigit := c >= '0' && c <= '9'
			isUpperCase := c >= 'A' && c <= 'Z'
			isLowerCase := c >= 'a' && c <= 'z'
			if !isDigit && !isUpperCase && !isLowerCase {
				t.Errorf("cur:%s un expext", curRandStr)
			}
		}
	}
}

func BenchmarkRandomIdGenerator_Generate(b *testing.B) {
	generator := new(RandomIdGenerator)
	for i := 0; i < b.N; i++ {
		generator.Generate()
	}
}
