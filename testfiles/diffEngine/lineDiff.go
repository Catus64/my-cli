package diffengine

func ComputeLineDiff(oldLines []string, newLines []string) (oldKept []bool, newKept []bool) {
	n := len(oldLines)
	m := len(newLines)

	// dp[i][j] = length of the LCS between oldLines[:i] and newLines[:j]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	// Fill the table bottom-up (row 0 / col 0 are already correctly 0 - "empty" base case)
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if oldLines[i-1] == newLines[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	oldKept = make([]bool, n)
	newKept = make([]bool, m)

	// Backtrack from the bottom-right corner to figure out
	// which specific lines were part of the LCS.
	i, j := n, m
	for i > 0 && j > 0 {
		switch {
		case oldLines[i-1] == newLines[j-1]:
			// this line exists in both — it's part of the LCS
			oldKept[i-1] = true
			newKept[j-1] = true
			i--
			j--
		case dp[i-1][j] >= dp[i][j-1]:
			// moving up was at least as good — this old line was removed
			i--
		default:
			// moving left was strictly better — this new line was added
			j--
		}
	}

	return oldKept, newKept
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
