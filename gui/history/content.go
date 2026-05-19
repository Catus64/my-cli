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

type commitNode struct {
	commit   Commit	// the commit data
	cx, cy   float32 // position of the bubble
	isAlternate bool // true = yellow, false = blue
	selected bool // true = selected
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
		for _, parentHash := range n.commit.Parents {
			if pn, ok := byHash[parentHash]; ok {
				
				dy := float64(n.cy - pn.cy)
				
				// Detect if this line skips a row (causing a collision)
				isLongDrop := math.Abs(dy) > float64(rowGapY*1.2)

				var arriveX, arriveY float64

				if isLongDrop {
					// Draw a C-Curve to gracefully route AROUND the bubble in the middle
					segments := 20
					
					// Control point: push the curve out to the right by 90 pixels
					cpX := (float64(pn.cx) + float64(n.cx))/2 + 90 
					cpY := (float64(pn.cy) + float64(n.cy))/2
					
					prevX := float64(pn.cx)
					prevY := float64(pn.cy)
					
					for i := 1; i <= segments; i++ {
						t := float64(i) / float64(segments)
						inv := 1.0 - t
						
						// Quadratic Bezier formula
						currX := inv*inv*float64(pn.cx) + 2*inv*t*cpX + t*t*float64(n.cx)
						currY := inv*inv*float64(pn.cy) + 2*inv*t*cpY + t*t*float64(n.cy)
						
						seg := canvas.NewLine(bubbleLine)
						seg.StrokeWidth = lineWidth
						seg.Position1 = fyne.NewPos(float32(prevX), float32(prevY))
						seg.Position2 = fyne.NewPos(float32(currX), float32(currY))
						objs = append(objs, seg)
						
						// Save the second-to-last point to calculate the perfect arrow angle
						if i == segments-1 {
							arriveX = prevX
							arriveY = prevY
						}
						prevX = currX
						prevY = currY
					}
				} else {
					// Draw normal straight line for adjacent commits
					line := canvas.NewLine(bubbleLine)
					line.StrokeWidth = lineWidth
					line.Position1 = fyne.NewPos(n.cx, n.cy)
					line.Position2 = fyne.NewPos(pn.cx, pn.cy)
					objs = append(objs, line)
					
					arriveX = float64(pn.cx)
					arriveY = float64(pn.cy)
				}

				// Calculate angles for the arrow head
				// We use arriveX/arriveY so the arrow perfectly matches the end of the curve!
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
	for idx, n := range hc.nodes {
		col := bubbleMain
		if n.isAlternate {
			col = bubbleAlt
		}

		// Add a label at top of current commit
		if idx == 0 {
			tag := canvas.NewText("Latest Save", color.RGBA{R: 100, G: 220, B: 100, A: 255})
			tag.TextSize = 11
			tag.TextStyle = fyne.TextStyle{Bold: true}
			tag.Alignment = fyne.TextAlignCenter
			tag.Resize(fyne.NewSize(100, 15))
			tag.Move(fyne.NewPos(n.cx-50, n.cy-bubbleRadius-20))
			objs = append(objs, tag)
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
		labelText := n.commit.SHA
		if len(labelText) > 5 {
			labelText = labelText[:5]
		}
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

	// collect all remaining commits not in main chain (mean branch commits)
	var altChain []string
	for _, c := range commits {
		if !seen[c.SHA] {
			altChain = append(altChain, c.SHA)
			seen[c.SHA] = true
		}
	}

	totalRows := (len(mainChain) + bubblesPerRow - 1) / bubblesPerRow

	var nodes []commitNode

	// place main chain nodes
	for i, sha := range mainChain {
		c := byHash[sha]
		row := i / bubblesPerRow
		col := i % bubblesPerRow

		var cx float32
        if row%2 == 0 {
            // Even rows go left to right 
            cx = float32(col+1) * bubbleGap
        } else {
            // Odd rows go right to left 
            cx = float32(bubblesPerRow-col) * bubbleGap
        }
		cy := mainY + float32(row)*rowGapY
		nodes = append(nodes, commitNode{commit: c, cx: cx, cy: cy})
	}

	// place alt commits relative to their merge point
	placedAlt := map[string]bool{}

	// Calculate a safe Y starting point BELOW the entire main chain
	lastMainRowY := mainY + float32(totalRows-1)*rowGapY
	baseAltY := lastMainRowY + rowGapY
	
	// Track how many separate branches we draw so they don't overlap each other
	branchIndex := 0

	originalLen := len(nodes)
	for i := originalLen - 1; i >= 0; i-- {
		node := nodes[i]
		
		if len(node.commit.Parents) < 2 {
			continue
		}

		altSHA := node.commit.Parents[1]
		
		// Assign this specific branch its own unique Y level
		altY := baseAltY + float32(branchIndex)*100

		altOffset := 0
		drewBranch := false

		for {
			if altSHA == "" || placedAlt[altSHA] {
				break
			}
			altC, ok := byHash[altSHA]
			if !ok {
				break
			}
			placedAlt[altSHA] = true
			drewBranch = true

			// Place going RIGHT from merge point
			nodes = append(nodes, commitNode{
				commit: altC,
				cx:     node.cx + float32(altOffset)*bubbleGap, // ← right not left
				cy:     altY,
				isAlternate:  true,
			})

			if len(altC.Parents) == 0 {
				break
			}
			
			// Stop if parent is already in main chain
			if _, inMain := func() (*commitNode, bool) {
				for i := range nodes {
					if nodes[i].commit.SHA == altC.Parents[0] && !nodes[i].isAlternate {
						return &nodes[i], true
					}
				}
				return nil, false
			}(); inMain {
				break
			}

			altSHA = altC.Parents[0]
			altOffset++
		}
		
		// If we drew a branch, increment the index to push the NEXT branch lower
		if drewBranch {
			branchIndex++
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

func HistoryPageContent(repoPath string, app fyne.App) fyne.CanvasObject {
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

	widthMargin := canvas.NewRectangle(color.Transparent)
	widthMargin.SetMinSize(fyne.NewSize(30, 0))

	heightMargin := canvas.NewRectangle(color.Transparent)
	heightMargin.SetMinSize(fyne.NewSize(0, 20))

	header := container.NewVBox(heightMargin, title, subtitle, heightMargin)
	paddedHeader := container.NewBorder(nil, nil, widthMargin, widthMargin, header)

	// Use NewBorder to let graphScroll fill all remaining space
	return container.NewBorder(
		container.NewPadded(paddedHeader),
		nil,
		nil,
		nil,
		container.NewPadded(graphScroll), // fills all available space
	)
}