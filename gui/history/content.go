package history

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Commit struct {
	Hash    string
	Message string
	Author  string
	Date    string
	Parents []string // hashes of parent commits (0, 1, or 2 entries)
}

type commitNode struct {
	commit   Commit
	cx, cy   float32 // centre position
	isAlt    bool    // alternate parent?
	selected bool
}

type HistoryCanvas struct {
	widget.BaseWidget
	nodes      []commitNode
	onSelect   func(Commit)
	canvasSize fyne.Size
}

type historyRenderer struct {
	hc      *HistoryCanvas
	objects []fyne.CanvasObject
}

func (r *historyRenderer) Layout(_ fyne.Size) {}
func (r *historyRenderer) MinSize() fyne.Size { return r.hc.canvasSize }

func (r *historyRenderer) Refresh() {
	r.objects = r.buildObjects()
	canvas.Refresh(r.hc)
}

func (r *historyRenderer) Destroy() {}

func (r *historyRenderer) Objects() []fyne.CanvasObject {
	if len(r.objects) == 0 {
		r.objects = r.buildObjects()
	}
	return r.objects
}

func (r *historyRenderer) buildObjects() []fyne.CanvasObject {
	var objs []fyne.CanvasObject
	hc := r.hc

	// index nodes by hash for line drawing
	nodeByHash := map[string]*commitNode{}
	for i := range hc.nodes {
		nodeByHash[hc.nodes[i].commit.Hash] = &hc.nodes[i]
	}

	// draw connector lines first (behind bubbles)
	for _, n := range hc.nodes {
		for pi, parentHash := range n.commit.Parents {
			if pn, ok := nodeByHash[parentHash]; ok {
				line := canvas.NewLine(colorLine)
				line.StrokeWidth = lineThickness
				// use Move/Resize trick: fyne.Line uses Position1/Position2
				line.Position1 = fyne.NewPos(n.cx, n.cy)
				line.Position2 = fyne.NewPos(pn.cx, pn.cy)
				_ = pi
				objs = append(objs, line)
			}
		}
	}

	// draw bubbles and labels
	for _, n := range hc.nodes {
		col := colorMain
		if n.isAlt {
			col = colorAlt
		}

		// outer ring when selected
		if n.selected {
			ring := canvas.NewCircle(colorSelect)
			ring.Move(fyne.NewPos(n.cx-bubbleRadius-4, n.cy-bubbleRadius-4))
			ring.Resize(fyne.NewSize(bubbleDiam+8, bubbleDiam+8))
			objs = append(objs, ring)
		}

		circle := canvas.NewCircle(col)
		circle.Move(fyne.NewPos(n.cx-bubbleRadius, n.cy-bubbleRadius))
		circle.Resize(fyne.NewSize(bubbleDiam, bubbleDiam))
		objs = append(objs, circle)

		// first 5 chars of hash
		label := canvas.NewText(n.commit.Hash[:5], colorText)
		label.TextSize = 11
		label.TextStyle = fyne.TextStyle{Bold: true}
		label.Alignment = fyne.TextAlignCenter
		// centre text roughly inside circle
		label.Move(fyne.NewPos(n.cx-bubbleRadius+20, n.cy-8))
		objs = append(objs, label)

		// small tag below bubble
		tag := canvas.NewText(n.commit.Date, color.RGBA{R: 180, G: 180, B: 180, A: 255})
		tag.TextSize = 9
		tag.Move(fyne.NewPos(n.cx-bubbleRadius+4, n.cy+bubbleRadius+4))
		objs = append(objs, tag)
	}

	return objs
}

var sampleCommits = []Commit{
	{Hash: "a1b2c3d4e5", Message: "Fix login bug", Author: "Alice", Date: "2024-05-01", Parents: []string{"b2c3d4e5f6"}},
	{Hash: "b2c3d4e5f6", Message: "Merge feature branch", Author: "Bob", Date: "2024-04-30", Parents: []string{"c3d4e5f6g7", "f1e2d3c4b5"}},
	{Hash: "c3d4e5f6g7", Message: "Add dashboard UI", Author: "Alice", Date: "2024-04-28", Parents: []string{"d4e5f6g7h8"}},
	{Hash: "d4e5f6g7h8", Message: "Refactor auth module", Author: "Charlie", Date: "2024-04-25", Parents: []string{"e5f6g7h8i9"}},
	{Hash: "e5f6g7h8i9", Message: "Initial commit", Author: "Alice", Date: "2024-04-20", Parents: []string{}},
	// alternate parent of the merge commit – shown but not expanded
	{Hash: "f1e2d3c4b5", Message: "Feature: dark mode", Author: "Bob", Date: "2024-04-27", Parents: []string{"d4e5f6g7h8"}},
}

func findCommit(hash string, commits []Commit) (Commit, bool) {
	for _, c := range commits {
		if c.Hash == hash {
			return c, true
		}
	}
	return Commit{}, false
}

func buildChain(head string, commits []Commit, maxDepth int) []Commit {
	var chain []Commit
	seen := map[string]bool{}
	current := head

	for depth := 0; depth < maxDepth; depth++ {
		c, ok := findCommit(current, commits)
		if !ok || seen[current] {
			break
		}
		seen[current] = true
		chain = append(chain, c)

		if len(c.Parents) == 0 {
			break
		}
		// follow main (first) parent
		current = c.Parents[0]
	}
	return chain
}

const (
	bubbleRadius  float32 = 28
	bubbleDiam    float32 = bubbleRadius * 2
	colSpacing    float32 = 110 // horizontal distance between bubble centres
	mainRowY      float32 = 80  // Y centre of the main commit row
	altRowY       float32 = 180 // Y centre of alternate parent row
	lineThickness float32 = 2
	canvasH       float32 = 280
)

var (
	colorMain   = color.RGBA{R: 80, G: 140, B: 255, A: 255}  // blue – main chain
	colorAlt    = color.RGBA{R: 255, G: 160, B: 60, A: 255}  // orange – merge alt
	colorLine   = color.RGBA{R: 180, G: 180, B: 180, A: 255} // grey connector
	colorSelect = color.RGBA{R: 255, G: 220, B: 60, A: 255}  // yellow highlight
	colorText   = color.White
)

func newHistoryCanvas(chain []Commit, commits []Commit, onSelect func(Commit)) *HistoryCanvas {
	h := &HistoryCanvas{onSelect: onSelect}

	totalCols := len(chain)
	canvasW := float32(totalCols)*colSpacing + colSpacing
	h.canvasSize = fyne.NewSize(canvasW, canvasH)

	// main chain nodes (left to right, newest first)
	for i, c := range chain {
		cx := canvasW - float32(i+1)*colSpacing // newest on the left
		h.nodes = append(h.nodes, commitNode{commit: c, cx: cx, cy: mainRowY})

		// if this is a merge commit, add the alt parent as a leaf below
		if len(c.Parents) == 2 {
			altHash := c.Parents[1]
			if alt, ok := findCommit(altHash, commits); ok {
				altCX := cx - colSpacing*0.5 // offset slightly left
				h.nodes = append(h.nodes, commitNode{commit: alt, cx: altCX, cy: altRowY, isAlt: true})
			}
		}
	}

	h.ExtendBaseWidget(h)
	return h
}

func (h *HistoryCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &historyRenderer{hc: h}
}

func (h *HistoryCanvas) MinSize() fyne.Size { return h.canvasSize }

func (h *HistoryCanvas) Tapped(ev *fyne.PointEvent) {
	for i, n := range h.nodes {
		dx := ev.Position.X - n.cx
		dy := ev.Position.Y - n.cy
		if dx*dx+dy*dy <= bubbleRadius*bubbleRadius {
			// deselect all, select tapped
			for j := range h.nodes {
				h.nodes[j].selected = false
			}
			h.nodes[i].selected = true
			h.Refresh()
			if h.onSelect != nil && !n.isAlt {
				h.onSelect(h.nodes[i].commit)
			}
			return
		}
	}
}

func makeDetailPanel() (*fyne.Container, func(Commit)) {
	hashLabel := canvas.NewText("", color.RGBA{R: 120, G: 200, B: 255, A: 255})
	hashLabel.TextSize = 13
	hashLabel.TextStyle = fyne.TextStyle{Bold: true}

	msgLabel := widget.NewLabel("")
	msgLabel.Wrapping = fyne.TextWrapWord

	authorLabel := canvas.NewText("", color.RGBA{R: 180, G: 180, B: 180, A: 255})
	authorLabel.TextSize = 12

	dateLabel := canvas.NewText("", color.RGBA{R: 180, G: 180, B: 180, A: 255})
	dateLabel.TextSize = 12

	parentsLabel := canvas.NewText("", color.RGBA{R: 180, G: 180, B: 180, A: 255})
	parentsLabel.TextSize = 11

	divider := canvas.NewRectangle(color.RGBA{R: 60, G: 60, B: 60, A: 255})
	divider.SetMinSize(fyne.NewSize(2, 0))

	heading := canvas.NewText("Commit Details", color.White)
	heading.TextSize = 16
	heading.TextStyle = fyne.TextStyle{Bold: true}

	placeholder := canvas.NewText("Click a commit bubble to see details", color.RGBA{R: 120, G: 120, B: 120, A: 255})
	placeholder.TextSize = 12

	inner := container.NewVBox(heading, placeholder)

	update := func(c Commit) {
		hashLabel.Text = fmt.Sprintf("Hash: %s", c.Hash)
		hashLabel.Refresh()

		msgLabel.SetText(c.Message)

		authorLabel.Text = fmt.Sprintf("Author: %s", c.Author)
		authorLabel.Refresh()

		dateLabel.Text = fmt.Sprintf("Date: %s", c.Date)
		dateLabel.Refresh()

		parentStr := "none"
		if len(c.Parents) > 0 {
			parentStr = ""
			for i, p := range c.Parents {
				if i > 0 {
					parentStr += ", "
				}
				parentStr += p[:5]
			}
		}
		parentsLabel.Text = fmt.Sprintf("Parents: %s", parentStr)
		parentsLabel.Refresh()

		inner.Objects = []fyne.CanvasObject{
			heading,
			hashLabel,
			msgLabel,
			authorLabel,
			dateLabel,
			parentsLabel,
		}
		inner.Refresh()
	}

	panel := container.NewBorder(nil, nil, divider, nil,
		container.NewPadded(inner),
	)

	return panel, update
}

func HistoryPageContent(pathName string) fyne.CanvasObject {
	title := canvas.NewText("Commit History", color.White)
	title.TextSize = 40
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := canvas.NewText("Visual history of your repository", color.Gray{Y: 150})
	subtitle.TextSize = 15

	header := container.NewVBox(title, subtitle)

	// build chain (5 levels deep)
	chain := buildChain(sampleCommits[0].Hash, sampleCommits, 5)

	detailPanel, updateDetail := makeDetailPanel()

	hCanvas := newHistoryCanvas(chain, sampleCommits, func(c Commit) {
		updateDetail(c)
	})

	// wrap the drawing canvas in a scroll so it's navigable on small windows
	graphScroll := container.NewHScroll(hCanvas)
	graphScroll.SetMinSize(fyne.NewSize(0, canvasH+40))

	graphCard := container.NewVBox(
		canvas.NewText("Graph", color.RGBA{R: 208, G: 200, B: 200, A: 255}),
		graphScroll,
	)

	// left: graph + detail stacked; right: detail panel
	content := container.NewBorder(
		header,
		nil,
		nil,
		container.NewPadded(detailPanel),
		container.NewVBox(graphCard),
	)

	return content
}
