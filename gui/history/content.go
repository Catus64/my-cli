package history

import (
	"image/color"
	"math"
	"strings"

	gitCurrent "gocmd/testfiles/GitCurrent"
	githashread "gocmd/testfiles/GitHashRead"
	gitlog "gocmd/testfiles/GitLog"
	gitobj "gocmd/testfiles/GitObject"
	gitsave "gocmd/testfiles/GitSave"
	gitpath "gocmd/testfiles/Gitrepostruct"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Commit struct {
	SHA         string    // Full commit hash
	ShortSHA    string    // Short Form commit hash
	Author      string    // Commit author
	Date        string    // Commit date
	Message     string    // Commit message
	TreeSHA     string    // File tree snapshot
	Parents     []string  // parent commit SHA (can be 2 if merge)
	VersionNum  int		  // Version number 
	VersionName string	  // Version name 
	HasVersion  bool	  // Whether this commit is tagged as a version
}

const (
	bubbleRadius   float32 = 28 // radius of commit bubbles
	bubbleDiameter   float32 = bubbleRadius * 2 // diameter of commit bubbles
	bubbleGap    float32 = 110	// left and right gap between bubbles
	mainY     float32 = 80	// Y position of the first row
	altY      float32 = 180	// Y position of the first row of alternate branches 
	lineWidth     float32 = 2	// line thickness
	canvasHeight   float32 = 280			
	bubblesPerRow int     = 15	// Max bubbles per row
	rowGapY       float32 = 200	// Up and Down gap between rows
)

var (
	bubbleMain   = color.RGBA{R: 80, G: 140, B: 255, A: 255} // main bubbles
	bubbleAlt    = color.RGBA{R: 255, G: 160, B: 60, A: 255} // branch bubbles
	bubbleLine   = color.RGBA{R: 180, G: 180, B: 180, A: 255} // connector lines
	bubbleSelect = color.RGBA{R: 255, G: 220, B: 60, A: 255} // yellow ring for selected bubble
	bubbleText   = color.White // text color for bubbles
)

type mergeInfo struct {
    label string
    isAlt bool // true = orange (branch to branch), false = green (branch to main)
}

type commitNode struct {
	commit   Commit	// the commit data
	cx, cy   float32 // position of the bubble
	isAlternate bool // true = yellow, false = blue
	selected bool // true = selected
	isHead bool // true if this is the HEAD commit (latest)
	mergeLabels []mergeInfo
}

// Bubble Graph
type HistoryCanvas struct {
	widget.BaseWidget
	nodes      []commitNode // all bubble nodes to draw
	canvasSize fyne.Size 
	onSelect   func(Commit)
}

type historyRenderer struct {
	historycanvas *HistoryCanvas
	objects []fyne.CanvasObject
}

func (h *HistoryCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &historyRenderer{historycanvas: h}
}

func (h *HistoryCanvas) MinSize() fyne.Size { return h.canvasSize }

func (h *HistoryCanvas) Tapped(ev *fyne.PointEvent) {
	for i, n := range h.nodes {
		dx := ev.Position.X - n.cx // horizontal distance
		dy := ev.Position.Y - n.cy // vertical distance
		if dx*dx+dy*dy <= bubbleRadius*bubbleRadius {
			// if click is inside the circle
			for j := range h.nodes {
				h.nodes[j].selected = false // deselect all
			}
			h.nodes[i].selected = true // select the clicked one
			h.Refresh() //redraw
			if h.onSelect != nil {
				h.onSelect(h.nodes[i].commit) // callback
			}
			return
		}
	}
}

func (r *historyRenderer) Layout(_ fyne.Size) {}
func (r *historyRenderer) MinSize() fyne.Size { return r.historycanvas.canvasSize }
func (r *historyRenderer) Destroy() {}
func (r *historyRenderer) Refresh() {
	r.objects = r.build()
	canvas.Refresh(r.historycanvas)
}
func (r *historyRenderer) Objects() []fyne.CanvasObject {
	if len(r.objects) == 0 {
		r.objects = r.build()
	}
	return r.objects
}

func (r *historyRenderer) build() []fyne.CanvasObject {
	var objs []fyne.CanvasObject
	hc := r.historycanvas

	// index nodes by hash (SHA) 
	byHash := map[string]*commitNode{}
	for i := range hc.nodes {
		byHash[hc.nodes[i].commit.SHA] = &hc.nodes[i]
	}

	// draw connector lines and arrows first (behind bubbles)
	for _, n := range hc.nodes {
		for parentIndex, parentHash := range n.commit.Parents {
			if parentIndex == 1 {  
				continue // skip merge parent line, show label instead
			}

			if pn, ok := byHash[parentHash]; ok {
				
				// dy := float64(n.cy - pn.cy)
				
				// Detect if this line skips a row (causing a collision)
				var arriveX, arriveY float64

				if n.cy == pn.cy {
					// Same Row: Draw a normal straight horizontal line
					line := canvas.NewLine(bubbleLine)
					line.StrokeWidth = lineWidth
					line.Position1 = fyne.NewPos(pn.cx, pn.cy)
					line.Position2 = fyne.NewPos(n.cx, n.cy)
					objs = append(objs, line)
					
					arriveX = float64(pn.cx)
					arriveY = float64(pn.cy)
				} else {
					// Different Row: Draw an L-shape
					
					if math.Abs(float64(pn.cx-n.cx)) < 1 {
						// Same X, different Y — draw a pure vertical line
						vLine := canvas.NewLine(bubbleLine)
						vLine.StrokeWidth = lineWidth
						vLine.Position1 = fyne.NewPos(pn.cx, pn.cy)
						vLine.Position2 = fyne.NewPos(n.cx, n.cy)
						objs = append(objs, vLine)

						// Arrow comes from directly above — set arriveY slightly above bubble
						arriveX = float64(pn.cx)
						arriveY = float64(n.cy) - float64(bubbleRadius)

					} else {
						// Different X — draw normal L-shape
						vLine := canvas.NewLine(bubbleLine)
						vLine.StrokeWidth = lineWidth
						vLine.Position1 = fyne.NewPos(pn.cx, pn.cy)
						vLine.Position2 = fyne.NewPos(pn.cx, n.cy)
						objs = append(objs, vLine)

						hLine := canvas.NewLine(bubbleLine)
						hLine.StrokeWidth = lineWidth
						hLine.Position1 = fyne.NewPos(pn.cx, n.cy)
						hLine.Position2 = fyne.NewPos(n.cx, n.cy)
						objs = append(objs, hLine)

						arriveX = float64(pn.cx)
						arriveY = float64(n.cy)
					}
				}
				// Calculate angles for the arrow head
				// Use arriveX/arriveY so the arrow perfectly matches the end of the curve!
				arrowDx := float64(n.cx) - arriveX
				arrowDy := float64(n.cy) - arriveY
				angle := math.Atan2(arrowDy, arrowDx)

				// Arrow properties
				arrowLen := float64(14)      
				arrowAngle := math.Pi / 6    

				// Calculate where the tip should stop (edge of the CHILD bubble)
				tipX := float64(n.cx) - float64(bubbleRadius)*math.Cos(angle)
				tipY := float64(n.cy) - float64(bubbleRadius)*math.Sin(angle)

				leftX := tipX - arrowLen*math.Cos(angle-arrowAngle)
				leftY := tipY - arrowLen*math.Sin(angle-arrowAngle)
				rightX := tipX - arrowLen*math.Cos(angle+arrowAngle)
				rightY := tipY - arrowLen*math.Sin(angle+arrowAngle)

				arrow1 := canvas.NewLine(bubbleLine)
				arrow1.StrokeWidth = lineWidth
				arrow1.Position1 = fyne.NewPos(float32(tipX), float32(tipY))
				arrow1.Position2 = fyne.NewPos(float32(leftX), float32(leftY))
				objs = append(objs, arrow1)

				arrow2 := canvas.NewLine(bubbleLine)
				arrow2.StrokeWidth = lineWidth
				arrow2.Position1 = fyne.NewPos(float32(tipX), float32(tipY))
				arrow2.Position2 = fyne.NewPos(float32(rightX), float32(rightY))
				objs = append(objs, arrow2)
			}
		}
	}

	// draw bubbles and labels
	for _, n := range hc.nodes {
		col := bubbleMain
		if n.isAlternate {
			col = bubbleAlt
		}

		// Add a label at top of current commit
		if n.isHead {
			tag := canvas.NewText("Latest Save", color.RGBA{R: 100, G: 220, B: 100, A: 255})
			tag.TextSize = 11
			tag.TextStyle = fyne.TextStyle{Bold: true}
			tag.Alignment = fyne.TextAlignCenter
			tag.Resize(fyne.NewSize(100, 15))
			tag.Move(fyne.NewPos(n.cx-50, n.cy-bubbleRadius-20))
			objs = append(objs, tag)
		}

		// Show merge label above branch bubble
		if len(n.mergeLabels) > 0 {
			total := len(n.mergeLabels)
			for k, entry := range n.mergeLabels {
				labelColor := color.RGBA{R: 100, G: 220, B: 100, A: 255} // green
				if entry.isAlt {
					labelColor = color.RGBA{R: 255, G: 160, B: 60, A: 255} // orange
				}
				ml := canvas.NewText(entry.label, labelColor)
				ml.TextSize = 11
				ml.TextStyle = fyne.TextStyle{Bold: true}
				ml.Alignment = fyne.TextAlignCenter
				ml.Resize(fyne.NewSize(120, 15))
				ml.Move(fyne.NewPos(n.cx-60, n.cy-bubbleRadius-20-float32(total-1-k)*16))
				objs = append(objs, ml)
			}
		}

		// yellow ring when selected
		if n.selected {
			ring := canvas.NewCircle(bubbleSelect)
			ring.Move(fyne.NewPos(n.cx-bubbleRadius-4, n.cy-bubbleRadius-4))
			ring.Resize(fyne.NewSize(bubbleDiameter+8, bubbleDiameter+8))
			objs = append(objs, ring)
		}

		circle := canvas.NewCircle(col)
		circle.Move(fyne.NewPos(n.cx-bubbleRadius, n.cy-bubbleRadius))
		circle.Resize(fyne.NewSize(bubbleDiameter, bubbleDiameter))
		objs = append(objs, circle)

		// first 5 chars of SHA inside bubble
		labelText := n.commit.ShortSHA
		label := canvas.NewText(labelText, bubbleText)
		label.TextSize = 11
		label.TextStyle = fyne.TextStyle{Bold: true}
		label.Alignment = fyne.TextAlignCenter

		// Force the text box to be exactly as wide as the bubble
		label.Resize(fyne.NewSize(bubbleDiameter, 15)) 
		
		// Start the text box at the exact left edge of the bubble, and adjust the Y position slightly to center it vertically.
		label.Move(fyne.NewPos(n.cx-bubbleRadius, n.cy-7))

		objs = append(objs, label)
	}

	return objs
}

// ── load real commits from repo ───────────────────────────────────────────────

func loadCommits(repoPath string) ([]Commit, *gitpath.GitRepository) {
	repo, err := gitpath.Repo_find(repoPath, false)
	if err != nil || repo == nil {
		return nil, nil
	}

	headSHA, err := gitobj.Ref_Resolve(*repo, "HEAD")
	if err != nil || headSHA == nil {
		return nil, repo
	}

	branch, _ := gitCurrent.Get_Active_Branch(*repo)

	var commits []Commit
	seen := map[string]bool{}
	queue := []string{*headSHA} // start from HEAD mean start from current commit

	for len(queue) > 0 {
		sha := queue[0]
		queue = queue[1:]

		if seen[sha] || sha == "" {
			continue
		}
		seen[sha] = true

		obj, err := githashread.Object_Read(*repo, sha)
		if err != nil {
			continue
		}
		gitCommit, ok := obj.(*gitobj.GitCommit)
		if !ok {
			continue
		}
		gitCommit.Deserialize()

		date, author := gitlog.Format_Date_Author(string(gitCommit.KvlmDict.Dict["author"]))
		treeSHA := strings.TrimSpace(string(gitCommit.KvlmDict.Dict["tree"]))
		parents := gitobj.GetKvlmValues(gitCommit.KvlmDict.Dict, "parent")
		message := strings.TrimSpace(string(gitCommit.KvlmDict.Dict["data"]))

		c := Commit{
			SHA:      sha,
			ShortSHA: sha[:5],
			Author:   author,
			Date:     date,
			Message:  message,
			TreeSHA:  treeSHA,
			Parents:  parents,
		}

		if branch != "" {
			entry, err := gitsave.ReadVersionRef(*repo, branch, sha)
			if err == nil {
				c.HasVersion = true
				c.VersionNum = entry.Number
				c.VersionName = entry.Name
			}
		}

		commits = append(commits, c)

		// enqueue ALL parents — not just first
		for _, p := range parents {
			if !seen[p] {
				queue = append(queue, p)
			}
		}
	}

	return commits, repo
}

func getChangedFiles(repo *gitpath.GitRepository, c Commit) []string {
	if repo == nil {
		return nil
	}

	// build map of current commit tree
	current := map[string]string{}
	gitCurrent.TreeToMap(*repo, c.TreeSHA, "", current)

	// if no parents, all files are new
	if len(c.Parents) == 0 {
		var result []string
		for path := range current {
			result = append(result, path)
		}
		return result
	}

	// build map of parent tree
	parentObj, err := githashread.Object_Read(*repo, c.Parents[0])
	if err != nil {
		return nil
	}
	parentCommit, ok := parentObj.(*gitobj.GitCommit)
	if !ok {
		return nil
	}
	parentCommit.Deserialize()
	parentTreeSHA := strings.TrimSpace(string(parentCommit.KvlmDict.Dict["tree"]))

	parent := map[string]string{}
	gitCurrent.TreeToMap(*repo, parentTreeSHA, "", parent)

	var result []string

	// added or modified
	for path, sha := range current {
		if parentSHA, exists := parent[path]; exists {
			if parentSHA != sha {
				result = append(result, path)
			}
		} else {
			result = append(result, path)
		}
	}

	// deleted
	for path := range parent {
		if _, exists := current[path]; !exists {
			result = append(result, path)
		}
	}

	return result
}

func buildNodes(commits []Commit) ([]commitNode, float32, float32) {
	if len(commits) == 0 {
		return nil, 200, canvasHeight
	}

	byHash := map[string]Commit{}
	for _, c := range commits {
		byHash[c.SHA] = c
	}

	// walk main chain (first parent only)
	var mainChain []string
	seen := map[string]bool{}
	current := commits[0].SHA
	for {
		c, ok := byHash[current]
		if !ok || seen[current] {
			break
		}
		seen[current] = true
		mainChain = append(mainChain, current)
		if len(c.Parents) == 0 { // if no parents, mean it is first commit, stop here
			break
		}
		current = c.Parents[0]
	}

	// Reverse main chain to start from oldest commit
	for i, j := 0, len(mainChain)-1; i < j; i, j = i+1, j-1 {
		mainChain[i], mainChain[j] = mainChain[j], mainChain[i]
	}

	// collect all remaining commits not in main chain (mean branch commits)
	var altChain []string
	for _, c := range commits {
		if !seen[c.SHA] {
			altChain = append(altChain, c.SHA)
			seen[c.SHA] = true
		}
	}

	totalRows := 1

	var nodes []commitNode

	// place main chain nodes
	for i, sha := range mainChain {
		c := byHash[sha]
		cx := float32(i+1) * bubbleGap
		cy := mainY
		isLatest := i == len(mainChain)-1   // last node is latest commit
		nodes = append(nodes, commitNode{commit: c, cx: cx, cy: cy, isHead: isLatest})
	}

	// place alt commits relative to their merge point
	// placedAlt := map[string]bool{}

	// Calculate a safe Y starting point BELOW the entire main chain
	lastMainRowY := mainY + float32(totalRows-1)*rowGapY
	baseAltY := lastMainRowY + rowGapY
	
	// Track how many separate branches we draw so they don't overlap each other
	branchIndex := 0

	originalLen := len(nodes)
	for i := 0; i < originalLen; i++ {
		node := nodes[i]
		
		if len(node.commit.Parents) < 2 {
			continue
		}

		altSHA := node.commit.Parents[1]
		
		// 1. Gather all commits in this branch line first to see where it ends up
		var branchCommits []Commit
		currSHA := altSHA
		var hitNode *commitNode

		for {
			if currSHA == "" {
				break
			}

			// Check if this commit is already placed anywhere (main or alt)
			var found *commitNode
			for j := range nodes {
				if nodes[j].commit.SHA == currSHA {
					found = &nodes[j]
					break
				}
			}

			// Stop tracing if we hit a bubble that has already been drawn
			if found != nil {
				hitNode = found
				break 
			}

			altC, ok := byHash[currSHA]
			if !ok {
				break
			}

			branchCommits = append(branchCommits, altC)

			if len(altC.Parents) == 0 {
				break
			}
			currSHA = altC.Parents[0]
		}

		if len(branchCommits) == 0 {
			continue // Nothing new to draw
		}
		
		// 2. Determine if we can safely share the row with the hit node
		needsNewRow := true
		var altY float32
		
		if hitNode != nil && hitNode.isAlternate {
			targetY := hitNode.cy
			
			// Calculate where our new bubbles will physically sit
			rightmostX := node.cx + (bubbleGap / 2)
			leftmostX := rightmostX - float32(len(branchCommits)-1)*bubbleGap
			
			// Calculate the entire X span from our new bubbles all the way back to the hitNode
			spanMin := leftmostX
			if hitNode.cx < spanMin { spanMin = hitNode.cx }
			spanMax := rightmostX
			if hitNode.cx > spanMax { spanMax = hitNode.cx }
			
			// LINE-OF-SIGHT CHECK: 
			// If any existing bubble is sitting inside this horizontal space, our connecting 
			// line would cut right through it! (Meaning this is a parallel branch fork).
			collision := false
			for _, n := range nodes {
				// Only check bubbles on the target row, and ignore the hitNode itself
				if n.cy == targetY && n.commit.SHA != hitNode.commit.SHA {
					if n.cx >= spanMin && n.cx <= spanMax {
						collision = true
						break
					}
				}
			}
			
			// If the visual path is totally clear, it's a direct continuation of the same branch!
			if !collision {
				altY = targetY
				needsNewRow = false
			}
		}

		// 3. If it collided (it's a parallel fork) OR pointed to Main, drop to a brand new row
		if needsNewRow {
			altY = baseAltY + float32(branchIndex)*100
			branchIndex++
		}

		// 4. Find the X starting point (divergence point, not merge point)
		var divergeX float32 = node.cx // fallback

		if hitNode != nil && !needsNewRow {
			// Continuing from an existing alt bubble - start from there
			divergeX = hitNode.cx
		} else {
			// New row - anchor to where the branch diverged from main chain
			oldest := branchCommits[len(branchCommits)-1]
			if len(oldest.Parents) > 0 {
				for k := range nodes {
					if nodes[k].commit.SHA == oldest.Parents[0] {
						divergeX = nodes[k].cx
						break
					}
				}
			}
			if hitNode != nil && hitNode.isAlternate && needsNewRow {
				// Align new row with rightmost bubble in the row above
				for k := range nodes {
					if nodes[k].cy == hitNode.cy && nodes[k].cx > divergeX {
						divergeX = nodes[k].cx
					}
				}
			}
		}

		// 5. Place bubbles oldest → newest going RIGHT from divergence point
		for j := len(branchCommits) - 1; j >= 0; j-- {
			altC := branchCommits[j]
			posFromLeft := len(branchCommits) - 1 - j // 0=oldest, increases toward newest

			var cx float32
			if hitNode != nil && hitNode.isAlternate && needsNewRow {
				// Align under the rightmost bubble of the row above
				cx = divergeX + float32(posFromLeft)*bubbleGap
			} else if !needsNewRow {
				// Continuing from existing alt bubble — full gap
				cx = divergeX + float32(posFromLeft+1)*bubbleGap
			} else {
				// New row from main chain divergence — shift left slightly
				cx = divergeX + float32(posFromLeft+1)*bubbleGap - (bubbleGap / 2)
			}

			if j == 0 {
				nodes = append(nodes, commitNode{
					commit:      altC,
					cx:          cx,
					cy:          altY,
					isAlternate: true,
					mergeLabels: []mergeInfo{{"Merge→" + node.commit.ShortSHA, false}}, // green
				})
				continue
			}

			// if j not equal to 0 nodes, no label:
			nodes = append(nodes, commitNode{commit: altC, cx: cx, cy: altY, isAlternate: true})
		}
	}

	// Branch-to-branch merge labels
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]

		if !node.isAlternate || len(node.commit.Parents) < 2 {
			continue
		}
		mergeParentSHA := node.commit.Parents[1]
		for j := range nodes {
			if nodes[j].commit.SHA == mergeParentSHA && nodes[j].isAlternate {
				nodes[j].mergeLabels = append(nodes[j].mergeLabels,
					mergeInfo{"Merge→" + node.commit.ShortSHA, true}) // orange
				break
			}
		}
	}

	// Dynamically calculate height based on actual node positions
	var maxX, maxY float32
	for _, n := range nodes {
		if n.cx > maxX { maxX = n.cx }
		if n.cy > maxY { maxY = n.cy }
	}
	canvasW := maxX + bubbleRadius + 60
	totalH  := maxY + bubbleRadius + 60

	return nodes, canvasW, totalH
}

func HistoryPageContent(repoPath string, app fyne.App, window fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Save History", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("View past changes and snapshots", color.Gray{Y: 150})
	subtitle.TextSize = 15

	// load real commits
	commits, repo := loadCommits(repoPath)

	if len(commits) == 0 {
		empty := canvas.NewText("No commits found in this repository.", color.Gray{Y: 150})
		empty.TextSize = 16
		return container.NewVBox(
			canvas.NewRectangle(color.Transparent),
			title,
			subtitle,
			canvas.NewRectangle(color.Transparent),
			empty,
		)
	}

	// build graph nodes
	nodes, canvasW, canvasH := buildNodes(commits)

	// graph canvas widget
	hCanvas := &HistoryCanvas{
		nodes:      nodes,
		canvasSize: fyne.NewSize(canvasW, canvasH),
		onSelect: func(c Commit) {
			// Get changed files
			files := getChangedFiles(repo, c)
			// Open new window
			ShowSaveHistoryWindow(app, c, files)
		},
	}
	hCanvas.ExtendBaseWidget(hCanvas)

	graphScroll := container.NewScroll(hCanvas)

	// Latest Save label
	latestSHA := ""
	for _, n := range nodes {
		if n.isHead {
			latestSHA = n.commit.ShortSHA
			break
		}
	}
	latestLabel := canvas.NewText("Latest Save: "+latestSHA, color.White)
	latestLabel.TextSize = 22
	latestLabel.TextStyle = fyne.TextStyle{Bold: true}

	underline := canvas.NewRectangle(color.White)
	underline.SetMinSize(fyne.NewSize(200, 2))

	latestLabelBox := container.NewVBox(latestLabel, underline)

	// Search bar
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search Bubble")
	searchSized := container.NewGridWrap(fyne.NewSize(200, searchEntry.MinSize().Height), searchEntry)

	searchBtn := widget.NewButton("Search", func() {
		query := strings.TrimSpace(searchEntry.Text)
		if query == "" {
			dialog.ShowInformation("Empty Search", "Please enter a save key to search.", window)
			return
		}
		// deselect all first
		for j := range hCanvas.nodes {
			hCanvas.nodes[j].selected = false
		}

		// find matching node
		found := false
		for j, n := range hCanvas.nodes {
			if strings.HasPrefix(n.commit.SHA, query) || strings.HasPrefix(strings.ToLower(n.commit.SHA), strings.ToLower(query)) {
				hCanvas.nodes[j].selected = true
				// scroll to the bubble
				graphScroll.ScrollToTop()
				graphScroll.Offset = fyne.NewPos(n.cx-200, n.cy-100)
				graphScroll.Refresh()
				found = true
				break
			}
		}

		if !found {
			dialog.ShowInformation("Not Found", "No matching save found.", window)
		}
		hCanvas.Refresh()
	})
	searchBtn.Importance = widget.HighImportance

	searchRow := container.NewHBox(searchSized, searchBtn)

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	graphMargin := canvas.NewRectangle(color.Transparent)
	graphMargin.SetMinSize(fyne.NewSize(0, 40))

	header := container.NewVBox(heightMargin, title, subtitle, heightMargin)
	paddedHeader := container.NewBorder(nil, nil, widthMargin, widthMargin, header)

	topBar := container.NewBorder(nil, nil, latestLabelBox, searchRow, nil)
	paddedTopBar := container.NewBorder(nil, nil, widthMargin, widthMargin, topBar)

	topGraph := container.NewVBox(paddedHeader, paddedTopBar, graphMargin)

	// Use NewBorder to let graphScroll fill all remaining space
	return container.NewBorder(
		topGraph,
		nil,
		nil,
		nil,
		container.NewPadded(graphScroll), // fills all available space
	)
}