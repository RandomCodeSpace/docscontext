package community

import (
	"math/rand"
)

// Graph represents an undirected weighted graph for community detection.
type Graph struct {
	Nodes []string
	nodeIndex map[string]int
	Edges     []Edge
	adjMatrix [][]float64 // dense for simplicity
	totalWeight float64
}

// Edge in the graph.
type Edge struct {
	Source, Target int
	Weight         float64
}

// NewGraph builds a graph from node IDs and weighted edges.
func NewGraph(nodes []string, edges [][3]any) *Graph {
	g := &Graph{
		Nodes:     nodes,
		nodeIndex: make(map[string]int, len(nodes)),
	}
	for i, n := range nodes {
		g.nodeIndex[n] = i
	}
	n := len(nodes)
	g.adjMatrix = make([][]float64, n)
	for i := range g.adjMatrix {
		g.adjMatrix[i] = make([]float64, n)
	}
	for _, e := range edges {
		src, _ := e[0].(string)
		tgt, _ := e[1].(string)
		w, _ := e[2].(float64)
		if w == 0 {
			w = 1.0
		}
		si, ok1 := g.nodeIndex[src]
		ti, ok2 := g.nodeIndex[tgt]
		if !ok1 || !ok2 {
			continue
		}
		g.adjMatrix[si][ti] += w
		g.adjMatrix[ti][si] += w
		g.totalWeight += w
		g.Edges = append(g.Edges, Edge{si, ti, w})
	}
	return g
}

// NodeIndex returns the index for a node ID.
func (g *Graph) NodeIndex(id string) (int, bool) {
	idx, ok := g.nodeIndex[id]
	return idx, ok
}

// Louvain runs the Louvain community detection algorithm.
// Returns a map from node index → community ID (integer).
func Louvain(g *Graph, maxIter int) []int {
	n := len(g.Nodes)
	if n == 0 {
		return nil
	}

	// Initialize: each node in its own community
	comm := make([]int, n)
	for i := range comm {
		comm[i] = i
	}

	if g.totalWeight == 0 {
		return comm
	}

	improved := true
	for iter := 0; iter < maxIter && improved; iter++ {
		improved = false
		// Random order
		order := rand.Perm(n)
		for _, i := range order {
			bestComm := comm[i]
			bestGain := 0.0

			// Neighbor communities
			neighborComms := map[int]float64{}
			for j := 0; j < n; j++ {
				if g.adjMatrix[i][j] > 0 {
					neighborComms[comm[j]] += g.adjMatrix[i][j]
				}
			}

			// Current community weight (excluding i)
			ki := g.nodeDegree(i)

			// Remove i from current community
			oldComm := comm[i]
			comm[i] = -1

			for c, w := range neighborComms {
				// Modularity gain (simplified)
				sigmaC := g.communityDegree(comm, c)
				gain := w - (ki*sigmaC)/(2*g.totalWeight)
				if gain > bestGain {
					bestGain = gain
					bestComm = c
				}
			}

			if bestComm != oldComm {
				improved = true
			}
			comm[i] = bestComm
		}
	}

	// Renumber communities 0..k-1
	renumber := map[int]int{}
	next := 0
	result := make([]int, n)
	for i, c := range comm {
		if _, ok := renumber[c]; !ok {
			renumber[c] = next
			next++
		}
		result[i] = renumber[c]
	}
	return result
}

func (g *Graph) nodeDegree(i int) float64 {
	var d float64
	for j := range g.adjMatrix[i] {
		d += g.adjMatrix[i][j]
	}
	return d
}

func (g *Graph) communityDegree(comm []int, c int) float64 {
	var d float64
	for i, ci := range comm {
		if ci == c {
			d += g.nodeDegree(i)
		}
	}
	return d
}

// HierarchicalLouvain runs Louvain at multiple levels.
// Returns a slice of levels, each level is a map nodeID → communityLabel.
func HierarchicalLouvain(g *Graph, maxLevels, maxIter int) [][]int {
	var levels [][]int
	current := Louvain(g, maxIter)
	levels = append(levels, current)

	for level := 1; level < maxLevels; level++ {
		// Count communities at current level
		commSet := map[int]bool{}
		for _, c := range current {
			commSet[c] = true
		}
		if len(commSet) <= 1 {
			break // Can't go higher
		}

		// Build super-graph where nodes = communities
		superNodes := make([]string, 0, len(commSet))
		superNodeIdx := map[int]int{}
		for c := range commSet {
			superNodeIdx[c] = len(superNodes)
			superNodes = append(superNodes, "")
		}

		superAdj := make([][]float64, len(superNodes))
		for i := range superAdj {
			superAdj[i] = make([]float64, len(superNodes))
		}

		for _, e := range g.Edges {
			ci := current[e.Source]
			cj := current[e.Target]
			if ci != cj {
				si := superNodeIdx[ci]
				sj := superNodeIdx[cj]
				superAdj[si][sj] += e.Weight
				superAdj[sj][si] += e.Weight
			}
		}

		superEdges := [][3]any{}
		for i := range superAdj {
			for j := i + 1; j < len(superAdj); j++ {
				if superAdj[i][j] > 0 {
					superEdges = append(superEdges, [3]any{superNodes[i], superNodes[j], superAdj[i][j]})
				}
			}
		}

		superGraph := NewGraph(superNodes, superEdges)
		superComm := Louvain(superGraph, maxIter)

		// Map back to original nodes
		nextLevel := make([]int, len(g.Nodes))
		for i, c := range current {
			si := superNodeIdx[c]
			nextLevel[i] = superComm[si]
		}

		// Check if we actually merged
		nextSet := map[int]bool{}
		for _, c := range nextLevel {
			nextSet[c] = true
		}
		if len(nextSet) >= len(commSet) {
			break
		}

		current = nextLevel
		levels = append(levels, current)
	}

	return levels
}
