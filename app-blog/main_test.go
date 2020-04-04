package main

import (
	"net/http"
	"testing"

	"github.com/deepaksinghkushwah/projects/app-blog/blog"
)

func BenchmarkList10(b *testing.B) {
	var r *http.Request
	var w http.ResponseWriter
	for n := 0; n < b.N; n++ {
		blog.List(w, r)
	}
}
