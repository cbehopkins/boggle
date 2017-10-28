package boggle

import (
	"errors"
	"log"
)

// Puzzle erm??? Not sure. Any ideas?
type Puzzle struct {
	Grid           [][]rune
	Visited        [][]bool
	newWordChan    chan string
	workerComplete chan struct{}
	dict           *DictMap
}

// StartWorker Start the worker funciton
// Needs a callback for what to do with each word you find
func (pz *Puzzle) StartWorker(sw func(string)) {
	pz.initWorker()
	pz.startWorker(sw)
}
func (pz *Puzzle) initWorker() {
	pz.newWordChan = make(chan string)
	pz.workerComplete = make(chan struct{})
}
func (pz *Puzzle) startWorker(sw func(string)) {
	go pz.newWordWorker(sw)
}
// SetDict set the dicitonary we wish to use
func (pz *Puzzle) SetDict(dct *DictMap) {
	pz.dict = dct
}
// NewPuzzle return a new puzzle of the specified size
func NewPuzzle(size int) *Puzzle {
	itm := new(Puzzle)
	itm.Visited = make([][]bool, size)
	for i := 0; i < size; i++ {
		itm.Visited[i] = make([]bool, size)
	}
	return itm
}
// Len reports on the size of the puzzle
func (pz Puzzle) Len() int {
	return len(pz.Grid)
}
// Copy one puzzle into anoteher destination one
func (pz Puzzle) Copy(dst *Puzzle) {
	dst.newWordChan = pz.newWordChan
	dst.dict = pz.dict
	pzLen := pz.Len()

	if dst.Len() != pzLen {
		dst.Grid = make([][]rune, pzLen)
		dst.Visited = make([][]bool, pzLen)
	}

	for i := 0; i < pzLen; i++ {
		row := make([]rune, pzLen)
		rowV := make([]bool, pzLen)
		for j := 0; j < pzLen; j++ {
			row[j] = pz.Grid[i][j]
			rowV[j] = pz.Visited[i][j]
		}
		dst.Grid[i] = row
		dst.Visited[i] = rowV
	}
}
func (pz Puzzle) newWordWorker(sw func(string)) {
	for word := range pz.newWordChan {
		sw(word)
	}
	close(pz.workerComplete)
}
func (pz Puzzle) rxWord(wrdPnt *string) (completeChan chan struct{}) {
	completeChan = make(chan struct{})
	go func() {
		word := <-pz.newWordChan
		log.Println("**** Received word", word)
		*wrdPnt = word
		close(completeChan)
	}()

	return
}
// Shutdown the generation
func (pz Puzzle) Shutdown() {
	log.Println("Shutdown Called")
	close(pz.newWordChan)
	log.Println("Shutdown sent")
	<-pz.workerComplete
	log.Println("Shutdown Complete")
}

// ErrVisited reports that wee have visited here before
var ErrVisited = errors.New("Error, we have visited this before")
// Coord is a struct of the coord in use
type Coord struct {
	xC int
	yC int
}

func (crd Coord) decode() (xC, yC int) {
	return crd.xC, crd.yC
}
func (crd *Coord) setX(xC int) {
	crd.xC = xC
}
func (crd *Coord) setY(yC int) {
	crd.yC = yC
}
func (crd Coord) getCoords(size int) []Coord {
	retArray := make([]Coord, 0)
	Xc, Yc := crd.decode()

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			candidateX := Xc + i
			candidateY := Yc + j
			if candidateX < 0 {
				continue
			}
			if candidateY < 0 {
				continue
			}
			if candidateX >= size {
				continue
			}
			if candidateY >= size {
				continue
			}
			newCrd := Coord{xC: candidateX, yC: candidateY}
			retArray = append(retArray, newCrd)
		}
	}
	return retArray
}

func (pz Puzzle) visit(runningWord string, vC Coord) error {
	if pz.visitedTrue(vC) {
		return ErrVisited
	}

	pz.setVisited(vC)
	defer pz.clearVisited(vC)

	//log.Println("Visiting Coordinate, with run", vC, runningWord)
	// To visit a coordinagte we:
	// Look and see if we currently have a word
	run := pz.getRune(vC)
	newWord := runningWord + string(run)
	var isWord, partial bool
	isWord, partial = pz.dict.partialExists(newWord)
	if isWord {
		// send it out on the results channel
		pz.newWord(newWord)
	}
	if partial {
		// Walk from this coord
		//log.Println("Walking from:", vC)
		pz.Walk(newWord, vC)
	}
	return nil
}
func (pz *Puzzle) setVisited(vC Coord) {
	Xc, Yc := vC.decode()
	pz.Visited[Yc][Xc] = true
}
func (pz *Puzzle) clearVisited(vC Coord) {
	Xc, Yc := vC.decode()
	pz.Visited[Yc][Xc] = false
}

func (pz Puzzle) visitedTrue(vC Coord) bool {
	Xc, Yc := vC.decode()
	return pz.Visited[Yc][Xc]
}
func (pz Puzzle) getRune(crd Coord) rune {
	Xc, Yc := crd.decode()
	return pz.Grid[Yc][Xc]
}
// newWord States that a new word has been found
func (pz Puzzle) newWord(inTxt string) {
	pz.newWordChan <- inTxt
}

// Walk the puzzle
// starting fromt he currel (partially complete) puzzle
// and try everything at every location
func (pz Puzzle) Walk(runningWord string, startC Coord) error {
	pzCopy := new(Puzzle)
	pz.Copy(pzCopy)

	// Only ever work on a copy of the data
	// so that we can safely modify it
	if !pzCopy.visitedTrue(startC) {
		log.Fatalf("Que?%v", pzCopy)
	}
	// Calculate each co-ordinate we can visit
	for _, crd := range startC.getCoords(pz.Len()) {
		if pzCopy.visitedTrue(crd) {
			continue
		}
		// visit it
		err := pzCopy.visit(runningWord, crd)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// RunWalk Run a walk through the puzzle
// visit each coord in turn
func (pz Puzzle) RunWalk() {
	pzLen := pz.Len()
	for i := 0; i < pzLen; i++ {
		for j := 0; j < pzLen; j++ {
			pz.visit("", Coord{i, j})
		}
	}
}
