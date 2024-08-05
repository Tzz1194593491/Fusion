package fweb

import "testing"

var (
	urls = []string{"/a/b/c", "/e/f/g", "/h/i/j", "/hihi/*", "/hihi/nono/hihi"}
)

func createNode() *node {
	root := &node{}
	for _, url := range urls {
		root.insert(url, parsePattern(url), 0)
	}
	return root
}

func TestNode(t *testing.T) {
	root := createNode()
	root.search(parsePattern(urls[0]), 0)
	root.search(parsePattern("/hihi/nono"), 0)
}
