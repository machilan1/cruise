package cors

import (
	"strings"
	"testing"
)

func TestWildcard(t *testing.T) {
	t.Parallel()

	// foo*bar
	w := wildcard{"foo", "bar"}
	if !w.match("foobar") {
		t.Error("foo*bar should match foobar")
	}
	if !w.match("foobazbar") {
		t.Error("foo*bar should match foobazbar")
	}
	if w.match("foobaz") {
		t.Error("foo*bar should not match foobaz")
	}

	// foo*oof
	w = wildcard{"foo", "oof"}
	if w.match("foof") {
		t.Error("foo*oof should not match foof")
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	got := convert([]string{"A", "b", "C"}, strings.ToLower)
	want := []string{"a", "b", "c"}
	if got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseHeaderList(t *testing.T) {
	t.Parallel()

	got := parseHeaderList("header, second-header, THIRD-HEADER, Numb3r3d-H34d3r, Header_with_underscore Header.with.full.stop")
	want := []string{"Header", "Second-Header", "Third-Header", "Numb3r3d-H34d3r", "Header_with_underscore", "Header.with.full.stop"}
	if got[0] != want[0] || got[1] != want[1] || got[2] != want[2] || got[3] != want[3] || got[4] != want[4] || got[5] != want[5] {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseHeaderListEmpty(t *testing.T) {
	t.Parallel()

	if len(parseHeaderList("")) != 0 {
		t.Error("should be empty slice")
	}
	if len(parseHeaderList(" , ")) != 0 {
		t.Error("should be empty slice")
	}
}

func BenchmarkParseHeaderList(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseHeaderList("header, second-header, THIRD-HEADER")
	}
}

func BenchmarkParseHeaderListSingle(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseHeaderList("header")
	}
}

func BenchmarkParseHeaderListNormalized(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseHeaderList("Header1, Header2, Third-Header")
	}
}

func BenchmarkWildcard(b *testing.B) {
	w := wildcard{"foo", "bar"}

	b.Run("match", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			w.match("foobazbar")
		}
	})
	b.Run("too short", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			w.match("fobar")
		}
	})
}
