# Prim's algorithm for finding MSTs

This is an attempt at implementing and benchmarking a small program which finds the minimum spanning tree for a dense, undirected graph.

**Prim's** algorithm is an approach at building a Minimum Spanning Tree which builds a Tree, one node at a time. A **tree* has the characteristics that each vertex has only two edges (and as such, there is only a single route between two nodes in the given tree). This algorithm chooses the least weighted edge whenever adding a new node into the tree. A MinimumHeap is used to optimize choosing this loop and as such each iteration during the build cycle of the tree will be able to have a logN look up time while finding the minimal edge to add to the tree.

_This_ particular implementation of **Prim's** algorithm is an appraoch at an idiomatic approach in **golang** and as such decision decisions such as naming, interfaces and structure have been made with community best practices in mind.


## Setup

A go runtime environment is bootstrapped and accessible in the included `Vagrant` virtual machine. If not familiar with Vagrant, please refer to the installation [directions](https://www.vagrantup.com/docs/installation/).

To download provision and login to a virtual machine with the Go runtime:

```bash
$ vagrant up 
$ vagrant ssh
```

Now, change directories and run the `prim.sh` script

## Implementation Details

There are 3 concrete types in this implementation, **MinHeap**, **Tree** and **Graph**. 

**MinHeap** is an abstract datatype which exposes an interface for storing, sorting and fetching objects adhering to the **Node** and **minHeapNode** interfaces. This implementation delegates comparison operations to the `node` type being passed in and as such attempts to centralize any forced type assertions to the most practical and relevant place (within the details of the node itself).

**Graph** is an abstract Graph type which builds and maintains a graph backed by an **Adjacency List**. Each node in the graph's adjacency list is backed by a MinHeap which is allows for fast lookup time of the "smallest" edges connecting that node to other nodes in the graph.

Finally **Tree** is a type which is a low level type for building and maintaining a set of connected Nodes and Edges. A **Tree** is characterized by having only a single connection between any two nodes. The **prim's** algorithm builds and returns a **tree** which is considered the minimal spanning tree for a given graph.


```bash
$ cd /opt/prims-mst
$ ./prims.sh
```

Alternatively, if you have a `go` compiler installed locally you don't need to install `vagrant` and can just run `./prims.sh` locally (assuming you are in a bash-friendly environment).

## Next Steps

I'd like to explore some additional optimizations around **Prim's** algorithm. Most notably, optimizing the "greediness" and how we choose minimally weighted edges.

First, **Prim's** is an inherently greedy algorithm. It would be interesting to choose a starting node in an non-randomized fashion. For instance, some basic heuristics could be applied to pick the most reasonable starting node. We could choose the node with the _fewest_ number of edges (to allow greater flexibility later on while building the tree) or we could choose a node with the lowest overall weight.

Secondly, it would be interesting to try and delay picking the most optimal node until further processing has been done. For instance, let's say we are building a tree of **10** nodes. It would be interesting if we picked, say **4** starting nodes and built out 4 subtrees in parallel for the first few steps. This would allow us to make a more educated and _less greedy_ choice. After say, the 3rd or 4th step we could simply continue moving forward with a single tree, the tree with the over lowest weight thus far.

Finally, a **Fibonacci Heap** could be used to optimize the "selection" process of the next edge to add to the tree. Currently, a **Binary Heap** is used for selecting and maintaining a list of lowest edges on a node by node basis. Extracting a single node from said heap is an O(logn) operation. A **Fibonacci Heap** can perform the same operation in with constant tight bound, meaning that at upper bounds performance will scale optimally. This would be particularly optimal in densely populated graphs with many edges between each node.


