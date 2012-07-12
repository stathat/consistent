package consistent

import (
	"sort"
	"testing"
	"testing/quick"
)

func checkNum(num, expected int, t *testing.T) {
	if num != expected {
		t.Errorf("expected %d, got %d", expected, num)
	}
}

func TestNew(t *testing.T) {
	x := New()
	if x == nil {
		t.Errorf("expected obj")
	}
	checkNum(x.NumberOfReplicas, 20, t)
}

func TestAdd(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	checkNum(len(x.circle), 20, t)
	checkNum(len(x.sortedHashes), 20, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
	x.Add("qwer")
	checkNum(len(x.circle), 40, t)
	checkNum(len(x.sortedHashes), 40, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
}

func TestRemove(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Remove("abcdefg")
	checkNum(len(x.circle), 0, t)
	checkNum(len(x.sortedHashes), 0, t)
}

func TestRemoveNonExisting(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Remove("abcdefghijk")
	checkNum(len(x.circle), 20, t)
}

func TestGetEmpty(t *testing.T) {
	x := New()
	_, err := x.Get("asdfsadfsadf")
	if err == nil {
		t.Errorf("expected error")
	}
	if err != ErrEmptyCircle {
		t.Errorf("expected empty circle error")
	}
}

func TestGetSingle(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	f := func(s string) bool {
		y, err := x.Get(s)
		if err != nil {
			t.Logf("error: %q", err)
			return false
		}
		t.Logf("s = %q, y = %q", s, y)
		return y == "abcdefg"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

func TestGetMultiple(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	result, err := x.Get("ggg")
	if err != nil {
		t.Fatal(err)
	}
	if result != "opqrstu" {
		t.Errorf("expected 'opqrstu', got %q", result)
	}
	result, err = x.Get("hhh")
	if err != nil {
		t.Fatal(err)
	}
	if result != "abcdefg" {
		t.Errorf("expected 'abcdefg', got %q", result)
	}
	result, err = x.Get("iiiiii")
	if err != nil {
		t.Fatal(err)
	}
	if result != "hijklmn" {
		t.Errorf("expected 'hijklmn', got %q", result)
	}
}

func TestGetMultipleQuick(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	f := func(s string) bool {
		y, err := x.Get(s)
		if err != nil {
			t.Logf("error: %q", err)
			return false
		}
		t.Logf("s = %q, y = %q", s, y)
		return y == "abcdefg" || y == "hijklmn" || y == "opqrstu"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

func TestGetMultipleRemove(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	result, err := x.Get("ggg")
	if err != nil {
		t.Fatal(err)
	}
	if result != "opqrstu" {
		t.Errorf("expected 'opqrstu', got %q", result)
	}
	result, err = x.Get("hhh")
	if err != nil {
		t.Fatal(err)
	}
	if result != "abcdefg" {
		t.Errorf("expected 'abcdefg', got %q", result)
	}
	result, err = x.Get("iiiiii")
	if err != nil {
		t.Fatal(err)
	}
	if result != "hijklmn" {
		t.Errorf("expected 'hijklmn', got %q", result)
	}
	x.Remove("hijklmn")
	result, err = x.Get("ggg")
	if err != nil {
		t.Fatal(err)
	}
	if result != "opqrstu" {
		t.Errorf("expected 'opqrstu', got %q", result)
	}
	result, err = x.Get("hhh")
	if err != nil {
		t.Fatal(err)
	}
	if result != "abcdefg" {
		t.Errorf("expected 'abcdefg', got %q", result)
	}
	result, err = x.Get("iiiiii")
	if err != nil {
		t.Fatal(err)
	}
	if result != "opqrstu" {
		t.Errorf("expected 'opqrstu', got %q", result)
	}
}

func TestGetMultipleRemoveQuick(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	x.Remove("opqrstu")
	f := func(s string) bool {
		y, err := x.Get(s)
		if err != nil {
			t.Logf("error: %q", err)
			return false
		}
		t.Logf("s = %q, y = %q", s, y)
		return y == "abcdefg" || y == "hijklmn"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

func TestGetTwo(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	a, b, err := x.GetTwo("99999999")
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Errorf("a shouldn't equal b")
	}
	if a != "abcdefg" {
		t.Errorf("wrong a: %q", a)
	}
	if b != "opqrstu" {
		t.Errorf("wrong b: %q", b)
	}
}

func TestGetTwoQuick(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	f := func(s string) bool {
		a, b, err := x.GetTwo(s)
		if err != nil {
			t.Logf("error: %q", err)
			return false
		}
		if a == b {
			t.Logf("a == b")
			return false
		}
		if a != "abcdefg" && a != "hijklmn" && a != "opqrstu" {
			t.Logf("invalid a: %q", a)
			return false
		}

		if b != "abcdefg" && b != "hijklmn" && b != "opqrstu" {
			t.Logf("invalid b: %q", b)
			return false
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}

func TestGetTwoOnlyOneInCircle(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	a, b, err := x.GetTwo("99999999")
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Errorf("a shouldn't equal b")
	}
	if a != "abcdefg" {
		t.Errorf("wrong a: %q", a)
	}
	if b != "" {
		t.Errorf("wrong b: %q", b)
	}
}
