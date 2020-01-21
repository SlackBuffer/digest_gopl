// uses a binary tree to implement an insertion sort
package main

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

// sorts values in place
func Sort(values []int) {
	var root *tree
	// fmt.Printf("%#v\n", root == nil) // true

	// 生成节点树 root
	for _, v := range values {
		root = add(root, v)
	}
	// sort in place
	appendValues(values[:0], root)
}

// appends the elemens of t to values in order and returns the resulting slice
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		// append here
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	fmt.Printf("%#v\n", values)
	return values
}

func add(t *tree, value int) *tree {
	// 不断往左右两侧递归，直至 t 为 nil 时创建新的叶子节点
	if t == nil {
		// 创建叶子节点
		// equivalent to `return &tree{value: value}`
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func main() {
	a := []int{3, 4, 2, 424, 23, 1}
	Sort(a)
	fmt.Println(a)
}
