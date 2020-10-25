package calculator

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/path"
)

// GraphPackCalculator generates a graph of quantity permutations with the available pack sizes
type GraphPackCalculator struct {
	PackSizes []int
}

// multigraph of quantities, allowing for multiple weights (lines) between two nodes (edge)
type quantityGraph struct {
	packSizeCount int
	candidates    map[int]quantityNode
	*multi.WeightedDirectedGraph
}

type quantityNode struct {
	quantity int
}

const headroomMultiplier int = 50

// Calculate the required number of packs
func (c GraphPackCalculator) Calculate(quantity int) RequiredPacks {
	packs := make(RequiredPacks)

	if quantity <= 0 {
		return packs
	}

	sizes := c.PackSizes
	sort.Ints(sizes)

	// reduce the problem space when the quantity is far greater than the sum of available pack sizes
	if permutationClamp := sum(sizes) * headroomMultiplier; quantity > permutationClamp {
		largestSize := sizes[len(sizes)-1]
		// subtract packs to bring the quantity down to the clamp
		packs[largestSize] = int(math.Floor(float64(quantity-permutationClamp) / float64(largestSize)))
		quantity -= packs[largestSize] * largestSize
	}

	// create a graph with the initial quantity as root node
	qGraph := quantityGraph{
		packSizeCount:         len(sizes),
		candidates:            make(map[int]quantityNode),
		WeightedDirectedGraph: multi.NewWeightedDirectedGraph(),
	}
	rootNode := quantityNode{quantity}
	qGraph.AddNode(rootNode)

	// generate permutations by recursively subtracting packs, with pack sizes cascading in descending order
	// e.g. [5 4 3 2 1] -> [4 3 2 1] -> [3 2 1] -> [2 1] -> [1]
	for i := len(sizes); i >= 1; i-- {
		availableSizes := make([]int, i)
		copy(availableSizes, sizes[:i])
		sort.Sort(sort.Reverse(sort.IntSlice(availableSizes)))
		qGraph.subtractPacks(rootNode, availableSizes)
	}

	// find the shortest path to the quantity closest to zero
	// TODO: aid traversal by removing vertices & edges which don't lead to the candidate
	candidateNode := qGraph.closestCandidate()
	shortest, _ := path.AStar(rootNode, candidateNode, qGraph, nil)
	path, _ := shortest.To(candidateNode.ID())
	pathLength := len(path)

	// count each weighted line that forms the path as a used pack size
	for i, currentNode := range path {
		nextIndex := i + 1
		if nextIndex >= pathLength {
			break
		}

		lines := qGraph.WeightedLines(currentNode.ID(), path[nextIndex].ID())
		lines.Next()
		packs[int(lines.WeightedLine().Weight())]++
	}

	return packs
}

func (g *quantityGraph) subtractPacks(n quantityNode, packSizes []int) {
	for _, size := range packSizes {
		// stop generating permutations if we've found more paths to 0 than available pack sizes
		if nodesToZero := g.To(int64(0)); nodesToZero.Len() >= g.packSizeCount {
			break
		}

		// find or create a node by the subtracted quantity
		nextQuantity := n.quantity - size
		nextNode := quantityNode{nextQuantity}
		if existingNode := g.Node(nextNode.ID()); existingNode == nil {
			g.AddNode(nextNode)
		}

		// maintain unique weights for edges between two quantities to avoid unnecessary recalculations
		weight := float64(size)
		if g.hasWeightedLineFromTo(n, nextNode, weight) {
			continue
		}

		// link the nodes by pack size
		g.SetWeightedLine(g.NewWeightedLine(n, nextNode, weight))

		// track nodes which satisfy the required quantity, stopping at this depth
		if nextQuantity <= 0 {
			g.candidates[nextQuantity] = nextNode
			continue
		}

		// subtract from the next quantity, increasing depth
		g.subtractPacks(nextNode, packSizes)
	}
}

func (g *quantityGraph) hasWeightedLineFromTo(from graph.Node, to graph.Node, weight float64) bool {
	for _, line := range graph.WeightedLinesOf(g.WeightedLines(from.ID(), to.ID())) {
		if line.Weight() == weight {
			return true
		}
	}
	return false
}

func (g quantityGraph) Weight(xid, yid int64) (w float64, ok bool) {
	return path.UniformCost(g)(xid, yid)
}

func (g *quantityGraph) closestCandidate() quantityNode {
	// create a slice of quantities from the map keys
	quantities := make([]int, len(g.candidates))
	i := 0
	for k := range g.candidates {
		quantities[i] = k
		i++
	}

	// reverse sort so the closest candidate is first
	sort.Sort(sort.Reverse(sort.IntSlice(quantities)))

	return g.candidates[quantities[0]]
}

func (n quantityNode) ID() int64 {
	return int64(n.quantity)
}

func sum(array []int) int {
	sum := 0
	for _, v := range array {
		sum += v
	}
	return sum
}
