package main

import (
	"fmt"
	"log"
)

type minHeapNode interface {
	lessThan(minHeapNode) bool
}

type MinHeap struct {
	nodes []minHeapNode
}

func NewMinHeap(nodes []minHeapNode) *MinHeap {
	minHeap := &MinHeap{
		nodes: nodes,
	}

	if len(minHeap.nodes) <= 1 {
		return minHeap
	}

	// floyds algorithm is used to guarantee the heap property when a new
	// binary heap is created. This is down by starting at the bottom of
	// the heap and sifing nodes so the heap property is fulfilled
	for i := len(minHeap.nodes) / 2; i >= 0; i-- {
		minHeap.siftDown(i)
	}

	return minHeap
}

func (m *MinHeap) Insert(node minHeapNode) {
	m.nodes = append(m.nodes, node)
	m.siftUp(len(m.nodes) - 1)
}

func (m *MinHeap) siftUp(index int) {
	// ensure that the index is in range
	if index == 0 || index >= len(m.nodes) {
		return
	}

	parentIndex := (index - 1) / 2
	// if the child node is less than the parent, then swap the nodes in the backing array
	if m.nodes[index].lessThan(m.nodes[parentIndex]) {
		temp := m.nodes[parentIndex]
		m.nodes[parentIndex] = m.nodes[index]
		m.nodes[index] = temp
		m.siftUp(parentIndex)
	}
}

func (m *MinHeap) siftDown(index int) {
	// there isn't a single child, so stop sifting down
	if index*2+1 >= len(m.nodes) {
		return
	}

	childIndex := index*2 + 1
	childNode := m.nodes[index*2+1]

	// see if the second child node is less than the first, and if so swap it out
	if index*2+2 < len(m.nodes) && m.nodes[index*2+2].lessThan(childNode) {
		childNode = m.nodes[index*2+2]
		childIndex = index*2 + 2
	}

	// if the childNode is less than the parent, then we need to swap them out
	if childNode.lessThan(m.nodes[index]) {
		temp := m.nodes[index]
		m.nodes[index] = m.nodes[childIndex]
		m.nodes[childIndex] = temp
		m.siftDown(childIndex)
	}
}

func (m *MinHeap) Pop() (minHeapNode, error) {
	if len(m.nodes) == 0 {
		return nil, fmt.Errorf("MinHeap is empty")
	}

	node := m.nodes[0]
	m.nodes[0] = m.nodes[len(m.nodes)-1]
	m.nodes = m.nodes[:len(m.nodes)-1]
	m.siftDown(0)

	return node, nil
}

func (m MinHeap) Fetch() (minHeapNode, error) {
	if len(m.nodes) == 0 {
		return nil, fmt.Errorf("MinHeap is empty")
	}

	return m.nodes[0], nil
}

func (m MinHeap) Size() int {
	return len(m.nodes)
}

type Node interface{}

type ObjectNode struct {
	name string
}

func (o ObjectNode) String() string {
	return fmt.Sprintf("%s", o.name)
}

type edge struct {
	nodes  [2]Node
	weight int
}

func (e edge) String() string {
	return fmt.Sprintf("%s<-%d->%s", e.nodes[0], e.weight, e.nodes[1])
}

func (e edge) lessThan(m minHeapNode) bool {
	edge, ok := m.(*edge)
	if !ok {
		log.Fatal("Programming error. Invalid type assertion of minHeapNode type to *edge type")
	}

	return e.weight < edge.weight
}

// A graph is an object which represents an undirected graph which has many
// edges connecting different nodes.
type Graph struct {
	// edges are stored as a map of minHeaps, for each
	edges map[Node]*MinHeap
}

func NewGraph() *Graph {
	return &Graph{
		edges: make(map[Node]*MinHeap, 0),
	}
}

func (g *Graph) Insert(nodeA Node, nodeB Node, weight int) {
	// look up and insert the first Node; creating the edge in the process
	insert := func(node Node, edge *edge) {
		nodeEdges, ok := g.edges[node]
		if !ok {
			g.edges[node] = NewMinHeap([]minHeapNode{})
			nodeEdges = g.edges[node]
		}

		nodeEdges.Insert(edge)
	}

	edge := &edge{
		nodes:  [2]Node{nodeA, nodeB},
		weight: weight,
	}
	insert(nodeA, edge)
	insert(nodeB, edge)
}

func (g Graph) Get(node Node) (*MinHeap, error) {
	nodeEdges, ok := g.edges[node]
	if !ok {
		return nil, fmt.Errorf("Node not found")
	}

	return nodeEdges, nil
}

func (g *Graph) IterNodes(cb func(Node) bool) {
	for node, _ := range g.edges {
		if shouldContinue := cb(node); !shouldContinue {
			break
		}
	}
}

func (g *Graph) IterEdges(cb func(*edge) bool) {
	edges := make(map[*edge]interface{})
	for _, nodeEdges := range g.edges {
		poppedEdges := make([]*edge, 0)
		for {
			heapObj, err := nodeEdges.Pop()
			if err != nil {
				break
			}

			edge := heapObj.(*edge)
			poppedEdges = append(poppedEdges, edge)
			if _, ok := edges[edge]; !ok {
				edges[edge] = struct{}{}
				cb(edge)
			}
		}

		for _, edge := range poppedEdges {
			nodeEdges.Insert(edge)
		}
	}
}

// A tree type is a type which represents a set of vertices which are all only
// connected in one manner. This particular tree type represents a complete
// tree which is non cyclic and where each node only has a maximum of two
// edges.
type Tree struct {
	edges map[Node]map[Node]*edge
}

func NewTree() *Tree {
	return &Tree{
		edges: make(map[Node]map[Node]*edge, 0),
	}
}

func (t *Tree) Insert(nodeA, nodeB Node, nodeEdge *edge) error {
	if _, nodeAExists := t.edges[nodeA]; nodeAExists {
		if _, nodeBExists := t.edges[nodeB]; nodeBExists {
			return fmt.Errorf("both nodes already exist in the tree")
		}
	}

	insert := func(node Node, connectedNode Node, e *edge) (func(), error) {
		nodeEdges, ok := t.edges[node]
		if !ok {
			t.edges[node] = make(map[Node]*edge, 2)
			nodeEdges = t.edges[node]
		}

		if _, alreadyExists := nodeEdges[connectedNode]; alreadyExists {
			return nil, fmt.Errorf("Can not insert same node into the tree twice")
		}

		if len(nodeEdges) > 1 {
			return nil, fmt.Errorf("Each node can only have two edges in tree")
		}

		commit := func() {
			nodeEdges[connectedNode] = e
		}

		return commit, nil
	}

	commitA, err := insert(nodeA, nodeB, nodeEdge)
	if err != nil {
		return err
	}

	commitB, err := insert(nodeB, nodeA, nodeEdge)
	if err != nil {
		return err
	}

	// if both insertions are permitted, then go ahead and commit each
	// change to the the tree. This is done to simplify the implementation
	// so as to allow us to verify that this is a valid operation for
	// _both_ nodes before we commit either.
	commitA()
	commitB()

	return nil
}

func (t Tree) IterEdges(cb func(edge *edge) bool) {
	currentNode := func() Node {
		for node, edges := range t.edges {
			if len(edges) == 1 {
				return node
			}
		}

		return nil
	}()

	if currentNode == nil {
		return
	}

	otherNode := func(node Node, edge *edge) Node {
		if node == edge.nodes[0] {
			return edge.nodes[1]
		}

		return edge.nodes[0]
	}

	traversedEdges := make(map[*edge]interface{}, 0)
	for {
		edges := t.edges[currentNode]
		found := false
		for _, edge := range edges {
			if _, ok := traversedEdges[edge]; ok {
				continue
			}
			currentNode = otherNode(currentNode, edge)
			traversedEdges[edge] = struct{}{}
			cb(edge)
			found = true
		}

		if !found {
			break
		}
	}
}

func FindMST(g *Graph) (*Tree, error) {
	// this is an implementation of Prim's algorithm, which iterates
	// through the graph and builds a tree, node by node by choosing the
	// next node with the smallest edge that doesn't already exist in the
	// tree.

	// requiredNodes returns a map of all of the nodes that are required to
	// build the MST.
	requiredNodes := func() map[Node]interface{} {
		nodes := make(map[Node]interface{}, 0)
		iter := func(node Node) bool {
			nodes[node] = struct{}{}
			return true
		}

		g.IterNodes(iter)
		return nodes
	}()

	// otherNode returns the other node from an Edge; this is useful so we
	// can parse the underlying datatype and compose things together as we
	// have done thus far.
	otherEdgeNode := func(node Node, e *edge) Node {
		if e.nodes[0] == node {
			return e.nodes[1]
		}

		return e.nodes[0]
	}

	tree := NewTree()
	nodes := make([]Node, 0, len(requiredNodes))

	current := func() Node {
		for key, _ := range requiredNodes {
			return key
		}
		// should never be called unless requiredNodes is empty
		return nil
	}()

	for {
		if len(nodes) == len(requiredNodes)-1 {
			break
		}

		nodeEdges, err := g.Get(current)
		if err != nil {
			return nil, err
		}

		poppedEdges := make([]*edge, 0, nodeEdges.Size())

		for {
			minHeapItem, err := nodeEdges.Pop()
			// if there is nothing to pop from the tree, then that means the algorithm was too greedy and is failing!
			if err != nil {
				return nil, fmt.Errorf("Unable to find a Node to continue building the tree")
			}

			minWeightEdge, ok := minHeapItem.(*edge)
			if !ok {
				log.Fatal("Programming error. Type assertion was attempted which is impossible")
			}

			poppedEdges = append(poppedEdges, minWeightEdge)
			newNode := otherEdgeNode(current, minWeightEdge)

			// if the insertion was successful, then we updated the tree successfully and can exit
			if err := tree.Insert(current, newNode, minWeightEdge); err == nil {
				// update the current iterator to the new node
				current = newNode
				nodes = append(nodes, newNode)
				break
			}
		}

		// we put all of the popped edges back into the heap,
		// so this doesn't arbitrarily change the underlying
		// heap. Practically speaking, this is a bit of a code
		// smell so this needs a little refactoring.
		for _, poppedEdge := range poppedEdges {
			nodeEdges.Insert(poppedEdge)
		}
	}

	return tree, nil
}

func main() {
	log.Fatal("This program only exposes tests and benchmarks. Please run ./prims.sh instead")

}
