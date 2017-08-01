package boggle

import (
	"errors"
	"log"
)

type Puzzle struct {
	Grid           [][]rune
	Visited        [][]bool
	newWordChan    chan string
	workerComplete chan struct{}
	dict           *DictMap
}

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
func (pz *Puzzle) SetDict(dct *DictMap) {
	pz.dict = dct
}
func NewPuzzle(size int) *Puzzle {
	itm := new(Puzzle)
	itm.Visited = make([][]bool, size)
	for i := 0; i < size; i++ {
		itm.Visited[i] = make([]bool, size)
	}
	return itm
}
func (pz Puzzle) Len() int {
	return len(pz.Grid)
}
func (src Puzzle) Copy(dst *Puzzle) {
	dst.newWordChan = src.newWordChan
	dst.dict = src.dict
	pz_len := src.Len()

	if dst.Len() != pz_len {
		dst.Grid = make([][]rune, pz_len)
		dst.Visited = make([][]bool, pz_len)
	}

	for i := 0; i < pz_len; i++ {
		row := make([]rune, pz_len)
		rowV := make([]bool, pz_len)
		for j := 0; j < pz_len; j++ {
			row[j] = src.Grid[i][j]
			rowV[j] = src.Visited[i][j]
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
func (pz Puzzle) RxWord(wrdPnt *string) (completeChan chan struct{}) {
	completeChan = make(chan struct{})
	go func() {
		word := <-pz.newWordChan
		log.Println("**** Received word", word)
		*wrdPnt = word
		close(completeChan)
	}()

	return
}
func (pz Puzzle) Shutdown() {
	log.Println("Shutdown Called")
	close(pz.newWordChan)
	log.Println("Shutdown sent")
	<-pz.workerComplete
	log.Println("Shutdown Complete")
}

var ErrVisited = errors.New("Error, we have visited this before")

type Coord struct {
	xC int
	yC int
}

func (crd Coord) Decode() (xC, yC int) {
	return crd.xC, crd.yC
}
func (crd *Coord) SetX(xC int) {
	crd.xC = xC
}
func (crd *Coord) SetY(yC int) {
	crd.yC = yC
}
func (crd Coord) getCoords(size int) []Coord {
	ret_array := make([]Coord, 0)
	Xc, Yc := crd.Decode()

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
			ret_array = append(ret_array, newCrd)
		}
	}
	return ret_array
}

func (pz Puzzle) Visit(runningWord string, vC Coord) error {
	if pz.VisitedTrue(vC) {
		return ErrVisited
	}

	pz.SetVisited(vC)
	defer pz.ClearVisited(vC)

	//log.Println("Visiting Coordinate, with run", vC, runningWord)
	// To visit a coordinagte we:
	// Look and see if we currently have a word
	run := pz.GetRune(vC)
	newWord := runningWord + string(run)
	var isWord, partial bool
	isWord, partial = pz.dict.partialExists(newWord)
	if isWord {
		// send it out on the results channel
		pz.NewWord(newWord)
	}
	if partial {
		// Walk from this coord
		//log.Println("Walking from:", vC)
		pz.Walk(newWord, vC)
	}
	return nil
}
func (pz *Puzzle) SetVisited(vC Coord) {
	Xc, Yc := vC.Decode()
	pz.Visited[Yc][Xc] = true
}
func (pz *Puzzle) ClearVisited(vC Coord) {
	Xc, Yc := vC.Decode()
	pz.Visited[Yc][Xc] = false
}

func (pz Puzzle) VisitedTrue(vC Coord) bool {
	Xc, Yc := vC.Decode()
	return pz.Visited[Yc][Xc]
}
func (pz Puzzle) GetRune(crd Coord) rune {
	Xc, Yc := crd.Decode()
	return pz.Grid[Yc][Xc]
}
func (pz Puzzle) NewWord(inTxt string) {
	pz.newWordChan <- inTxt
}
func (pz Puzzle) Walk(runningWord string, startC Coord) error {
	pzCopy := new(Puzzle)
	pz.Copy(pzCopy)

	// Only ever work on a copy of the data
	// so that we can safely modify it
	if !pzCopy.VisitedTrue(startC) {
		log.Fatalf("Que?%v", pzCopy)
	}
	// Calculate each co-ordinate we can visit
	for _, crd := range startC.getCoords(pz.Len()) {
		if pzCopy.VisitedTrue(crd) {
			continue
		}
		// Visit it
		err := pzCopy.Visit(runningWord, crd)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (pz Puzzle) RunWalk() {
	pz_len := pz.Len()
	for i := 0; i < pz_len; i++ {
		for j := 0; j < pz_len; j++ {
			pz.Visit("", Coord{i, j})
		}
	}
}
