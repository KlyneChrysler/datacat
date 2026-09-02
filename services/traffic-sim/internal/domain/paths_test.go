package domain

import "testing"

func TestLoopingPathsStaysWithinSet(t *testing.T) {
	source := NewLoopingPaths([]string{"/a", "/b"})
	allowed := map[string]bool{"/a": true, "/b": true}

	for range 50 {
		if path := source.Next(); !allowed[path] {
			t.Fatalf("path %q outside the configured set", path)
		}
	}
}

func TestCrawlingPathsNeverRevisits(t *testing.T) {
	source := NewCrawlingPaths("/item-")
	seen := make(map[string]bool, 50)

	for range 50 {
		path := source.Next()
		if seen[path] {
			t.Fatalf("path %q visited twice", path)
		}
		seen[path] = true
	}
}
