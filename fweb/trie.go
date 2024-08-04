package fweb

import (
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，例如 /p/:lang， :lang 就是精确匹配
}

// matchChild 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// mathChildren 匹配多个节点，用于查询
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}

// insert 插入对应的节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.sortInsert(child)
	}
	child.insert(pattern, parts, height+1)
}

// search 查询对应的节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// sortInsert 排序并插入，该函数目的在于降低":"和"*"的优先级
func (n *node) sortInsert(child *node) {
	if strings.HasPrefix(child.part, "*") || strings.HasPrefix(child.part, ":") {
		n.children = append(n.children, child)
		return
	}
	n.children = append([]*node{child}, n.children...)
}
