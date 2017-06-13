package binaryTree

import (
	"errors"
	"strings"
	"math"
)

// Binary Tree
// References:
// 	https://appliedgo.net/bintree/
//	https://github.com/karlstroetmann/Algorithms/blob/master/SetlX/BinaryTree/binary-tree.stlx
type Tree struct {
	root *node
}

func NewTree() *Tree {
	return &Tree{}
}

// Inserts new key-Value pair to the Tree
// arguments:
//	- key int: key to find the Value
//	- Value string: Value to insert at the key
// returns: error if insertion fails
func (t *Tree) Insert(key int, value string) error {
	if t.root == nil {
		t.root = &node{Key: key, Value: value}
		t.root.restoreHeight()
		return nil
	}
	return t.root.insert(key, value)
}

// Finds Value in Tree for given key
// arguments:
//	- key int: key to find the Value for
// returns:
//	- string: Value of the key (empty string if no Value was found)
//	- bool: true if key was found in tree
func (t *Tree) Find(key int) (string, bool) {
	if t.root == nil {
		return "", false
	}
	return t.root.find(key)
}

// Deletes key-Value pair from Tree
// arguments:
//	- key int: key and a associated Value, which should be deleted
// returns: error if key was not found in Tree
func (t *Tree) Delete(key int) error {
	if t.root == nil {
		return errors.New("Cannot delete from an empty tree")
	}
	fakeParent := &node{right: t.root}
	err := t.root.delete(key, fakeParent)
	if err != nil {
		return err
	}
	if fakeParent.right == nil {
		t.root = nil
	}
	return nil
}

// Calls the given function with every key and Value ordered by the key (ascending)
// arguments:
//	- f func(int,string): the function to call, first argument is key, second Value
func (t *Tree) ForEach(f func(int, string)) {
	t.root.traverse(t.root, func(n *node) {
		f(n.Key, n.Value)
	})
}

func (t *Tree) Keys() []int {
	keys := []int{}
	t.ForEach(func(key int, value string) {
		keys = append(keys, key)
	})
	return keys
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Compare(a, e) == 0 {
			return true
		}
	}
	return false
}

func (t *Tree) Values(distinct bool) []string {
	values := []string{}
	t.ForEach(func(key int, value string) {
		if !distinct || !contains(values, value) {
			values = append(values, value)
		}
	})
	return values
}

func (t *Tree) ContainsValue(s string) bool {
	found := false
	t.ForEach(func(key int, value string) {
		if strings.Compare(s, value) == 0 {
			found = true;
			return;
		}
	})
	return found
}

func (t *Tree) ContainsKey(key int) bool {
	found := false
	t.ForEach(func(k int, value string) {
		if k == key {
			found = true;
			return;
		}
	})
	return found
}

func (t *Tree) ToMap() map[int]string {
	m := make(map[int]string)
	t.ForEach(func(key int, value string) {
		m[key] = value
	})
	return m
}

func (t *Tree) Height() int {
	max := 0
	t.root.traverse(t.root, func(n *node) {
		println(n.height)
		if n.height > max {
			max = n.height
		}
	})
	return max
}

func (t *Tree) Size() int {
	size := 0
	t.ForEach(func(key int, value string) {
		size++
	})
	return size
}

func (t *Tree) GetKeys(value string) []int {
	result := []int{}
	t.ForEach(func(k int, v string) {
		if strings.Compare(value, v) == 0 {
			result = append(result, k)
		}
	})
	return result
}

// Node structure for the Tree
// Key int: node key
// Value string: node Value
// left,right *node: the nodes left and right children
type node struct {
	Key         int
	Value       string
	left, right *node
	height      int
}

func (n *node) update(s *node) {
	n.Key, n.Value, n.left, n.right = s.Key, s.Value, s.left, s.right
}

func (n *node) delMin() (*node, int, string) {
	if n.left == nil {
		return n.right, n.Key, n.Value
	} else {
		ls, km, vm := n.left.delMin()
		n.left = ls
		return n, km, vm
	}
}

func (n *node) insert(key int, value string) error {
	if n == nil {
		return errors.New("Cannot insert a key into a nil tree")
	}

	switch {
	case key == n.Key:
		return nil
	case key < n.Key:
		if n.left == nil {
			n.left = &node{Key: key, Value: value}
			n.left.restoreHeight()
			return nil
		}
		err := n.left.insert(key, value)
		if err != nil {
			n.restore()
		}
		return err
	case key > n.Key:
		if n.right == nil {
			n.right = &node{Key: key, Value: value}
			n.right.restoreHeight()
			return nil
		}
		err := n.right.insert(key, value)
		if err != nil {
			n.restore()
		}
		return err
	}
	return nil
}

func (n *node) find(key int) (string, bool) {
	if n == nil {
		return "", false
	}
	switch {
	case key == n.Key:
		return n.Value, true
	case key < n.Key:
		return n.left.find(key)
	default:
		return n.right.find(key)
	}
}

func (n *node) findMax(parent *node) (*node, *node) {
	if n.right == nil {
		return n, parent
	}
	return n.right.findMax(n)
}

func (n *node) replaceNode(parent, replacement *node) error {
	if n == nil {
		return errors.New("replaceNode() not allowed on a nil node")
	}

	if n == parent.left {
		parent.left = replacement
		return nil
	}
	parent.right = replacement
	return nil
}

func (n *node) delete(key int, parent *node) error {
	if n == nil {
		return errors.New("Key to be deleted does not exist in the tree")
	}
	switch {
	case key < n.Key:
		err := n.left.delete(key, n)
		if err != nil {
			n.restore()
		}
		return err
	case key > n.Key:
		err := n.right.delete(key, n)
		if err != nil {
			n.restore()
		}
		return err
	default:
		if n.left == nil && n.right == nil {
			n.replaceNode(parent, nil)
			return nil
		}
		if n.left == nil {
			n.replaceNode(parent, n.right)
			return nil
		}
		if n.right == nil {
			n.replaceNode(parent, n.left)
			return nil
		}
		replacement, replParent := n.left.findMax(n)
		n.Key = replacement.Key
		n.Value = replacement.Value
		return replacement.delete(replacement.Key, replParent)
	}
}

func (n *node) restoreHeight() {
	switch {
	case n.left == nil && n.right == nil:
		n.height = 1
		return
	case n.left == nil:
		n.height = 1 + n.right.height
		return
	case n.right == nil:
		n.height = 1 + n.left.height
		return
	default:
		n.height = 1 + int(math.Max(float64(n.right.height), float64(n.left.height)))
	}
}
func (t *node) traverse(n *node, f func(*node)) {
	if n == nil {
		return
	}
	t.traverse(n.left, f)
	f(n)
	t.traverse(n.right, f)
}

func (n *node) restore() {
	if math.Abs(float64(n.left.height-n.right.height)) <= 1 {
		n.restoreHeight()
		return
	}
	if n.left.height > n.right.height {
		k1, v1, l1, r1 := n.Key, n.Value, n.left, n.right;
		k2, v2, l2, r2 := l1.Key, l1.Value, l1.left, l1.right;
		if l2.height >= r2.height {
			n.Key, n.Value, n.left, n.right = k2, v2, l2, &node{k1, v1, r2, r1, 1}
			n.right.restoreHeight()
		} else {
			k3, v3, l3, r3 := r2.Key, r2.Value, r2.left, r2.right;
			n.Key, n.Value, n.left, n.right = k3, v3, &node{k2, v2, l2, l3, 1}, &node{k1, v1, r3, r1, 1}
			n.left.restoreHeight()
			n.right.restoreHeight()
		}
	} else if n.right.height > n.left.height {
		k1, v1, l1, r1 := n.Key, n.Value, n.left, n.right
		k2, v2, l2, r2 := r1.Key, r1.Value, r1.left, r1.right;
		if r2.height >= l2.height {
			n.Key, n.Value, n.left, n.right = k2, v2, &node{k1, v1, l1, l2, 1}, r2
			n.left.restoreHeight()
		} else {
			k3, v3, l3, r3 := l2.Key, l2.Value, l2.left, l2.right
			n.Key, n.Value, n.left, n.right = k3, v3, &node{k1, v1, l1, l3, 1}, &node{k2, v2, r3, r2, 1}
			n.left.restoreHeight()
			n.right.restoreHeight()
		}
	}
	n.restoreHeight()
};
