package schema

// collectNodeMap returns id -> node map
func collectNodeMap(root *DenvclustrRoot) map[string]*Node {
	result := make(map[string]*Node)
	for _, n := range root.Nodes {
		result[string(n.Id)] = n
	}
	return result
}
