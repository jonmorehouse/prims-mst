package main

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomList() []int {
	// random list between 500 - 1000 items
	len := int(rand.Float32()*500) + 500
	len = 100
	list := make([]int, 0, len)
	for i := 0; i < len; i++ {
		list = append(list, rand.Intn(100))
	}

	return list
}

type minHeapTestNode struct {
	value int
}

func (m minHeapTestNode) lessThan(node minHeapNode) bool {
	testNode := node.(minHeapTestNode)
	return m.value < testNode.value
}

func TestBuildMinHeap(t *testing.T) {
	for i := 0; i < 1; i++ {
		nodes := make([]minHeapNode, 0)
		for _, val := range randomList() {
			nodes = append(nodes, minHeapTestNode{val})
		}

		minHeap := NewMinHeap(nodes)
		processed := 1
		lastNode, err := minHeap.Pop()
		if err != nil {
			t.Fatalf("Did not build heap properly")
		}

		for {
			node, err := minHeap.Pop()
			if err != nil {
				break
			}

			if node.lessThan(lastNode) {
				t.Fatalf("MinHeap does not fulfill the MinHeap property")
			}

			if err != nil && processed != len(nodes) {
				t.Fatalf("Did not pop all nodes off properly")
				break
			}

			lastNode = node
		}
	}
}

func TestGraphInsertion(t *testing.T) {
	mapGraph := make(map[[2]string]int, 0)
	mapGraph[[2]string{"a", "b"}] = 5
	mapGraph[[2]string{"a", "c"}] = 8
	mapGraph[[2]string{"a", "d"}] = 1
	mapGraph[[2]string{"a", "f"}] = 5

	mapGraph[[2]string{"b", "g"}] = 9
	mapGraph[[2]string{"b", "h"}] = 1
	mapGraph[[2]string{"b", "c"}] = 5

	mapGraph[[2]string{"c", "d"}] = 9
	mapGraph[[2]string{"c", "e"}] = 1
	mapGraph[[2]string{"c", "f"}] = 5

	mapGraph[[2]string{"d", "h"}] = 1
	mapGraph[[2]string{"d", "e"}] = 9

	mapGraph[[2]string{"e", "f"}] = 9
	mapGraph[[2]string{"e", "h"}] = 2

	graph := NewGraph()
	for nodes, weight := range mapGraph {
		graph.Insert(nodes[0], nodes[1], weight)
	}

	nodes := make(map[string]interface{}, 0)
	iter := func(n Node) bool {
		node, ok := n.(string)
		if !ok {
			t.Fatalf("Programming error, the node was not found correctly")
		}

		nodes[node] = struct{}{}
		return true
	}

	graph.IterNodes(iter)
	if len(nodes) != 8 {
		t.Fatalf("Graph did not store and house nodes correctly")
	}
}

func TestTree(t *testing.T) {
	testCases := []struct {
		nodes     [2]string
		weight    int
		shouldErr bool
	}{
		{[2]string{"a", "b"}, 2, false},
		{[2]string{"b", "c"}, 4, false},
		{[2]string{"c", "d"}, 2, false},
		{[2]string{"c", "d"}, 2, true},
		{[2]string{"d", "e"}, 4, false},
		{[2]string{"e", "f"}, 2, false},
		{[2]string{"c", "e"}, 2, true},
	}

	expectedWeight := 0
	tree := NewTree()
	nodes := make(map[string]bool, 0)

	for _, testCase := range testCases {
		err := tree.Insert(
			testCase.nodes[0],
			testCase.nodes[1],
			&edge{
				[2]Node{testCase.nodes[0], testCase.nodes[1]}, testCase.weight,
			})

		if err != nil && !testCase.shouldErr {
			t.Fatalf("Tree insertion failed when it shouldn't have")
		}

		if err == nil && testCase.shouldErr {
			t.Fatalf("Tree insertion did not fail when expected")
		}

		if !testCase.shouldErr {
			expectedWeight += testCase.weight
			nodes[testCase.nodes[0]] = false
			nodes[testCase.nodes[1]] = false
		}
	}

	iter := func(edge *edge) bool {
		nodes[edge.nodes[0].(string)] = true
		nodes[edge.nodes[1].(string)] = true
		expectedWeight -= edge.weight
		return true
	}
	tree.IterEdges(iter)

	if expectedWeight != 0 {
		t.Fatalf("Tree iteration did not return correct overall weight")
	}

	for _, value := range nodes {
		if value != true {
			t.Fatalf("Some nodes were not stored properly")
		}
	}
}

func TestMSTImplementation(t *testing.T) {
	graph := buildGraph(100)
	_, err := FindMST(graph)
	if err != nil {
		t.Fatalf("%s", err)
	}
}
