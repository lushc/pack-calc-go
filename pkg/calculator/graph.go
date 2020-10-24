package calculator

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
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
	graph := quantityGraph{
		packSizeCount:         len(sizes),
		candidates:            make(map[int]quantityNode),
		WeightedDirectedGraph: multi.NewWeightedDirectedGraph(),
	}
	root := quantityNode{quantity}
	graph.AddNode(root)

	// generate permutations by recursively subtracting packs, reducing the available packs each iteration
	for len(sizes) > 0 {
		graph.subtractPacks(root, sizes)
		sizes = sizes[:len(sizes)-1]
	}

	// TODO: find the shortest path to the quantity closest to zero, counting pack sizes

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
		for _, line := range graph.WeightedLinesOf(g.WeightedLines(n.ID(), nextNode.ID())) {
			if line.Weight() == weight {
				continue
			}
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
