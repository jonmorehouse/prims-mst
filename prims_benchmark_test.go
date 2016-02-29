package main

import (
	"encoding/base64"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func buildGraph(totalNodes int) *Graph {
	nodes := make([]string, 0, totalNodes)
	graph := NewGraph()

	for i := 0; i < totalNodes; i++ {
		// get ascii encoded version of number so we can encode it
		node := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(i)))
		node = strings.Replace(node, "=", "", -1)

		nodes = append(nodes, node)
	}

	for index, node := range nodes {
		for siblingIndex, siblingNode := range nodes {
			// generate a score between 10 and 1000
			if siblingIndex == index {
				continue
			}
			score := rand.Intn(1000-10) + 10
			graph.Insert(node, siblingNode, score)
		}
	}

	return graph
}

func benchmark(b *testing.B, graphSize int) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		graph := buildGraph(graphSize)
		b.StartTimer()

		_, err := FindMST(graph)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}
}

func BenchmarkPrims10(b *testing.B) { benchmark(b, 10) }

func BenchmarkPrims100(b *testing.B) { benchmark(b, 100) }

func BenchmarkPrims1000(b *testing.B) { benchmark(b, 1000) }

func BenchmarkPrims10000(b *testing.B) { benchmark(b, 10000) }
