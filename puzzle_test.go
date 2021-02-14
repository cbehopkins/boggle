package boggle

import (
	"log"
	"testing"

	"github.com/cbehopkins/wordlist"
)

func TestBas0(t *testing.T) {
	testPuzzle0 := NewPuzzle(2)
	testPuzzle1 := NewPuzzle(2)
	grid := [][]rune{
		{'a', 'b'},
		{'c', 'd'},
	}
	testPuzzle1.Grid = grid
	testPuzzle1.Copy(testPuzzle0)

}
func TestBuildDict0(t *testing.T) {
	testDict := []string{
		"be",
		"bad",
	}

	dictMap := NewDictMap(testDict)
	testStr := "bad"
	if !dictMap.Exists(testStr) {
		log.Fatalf("%s not found in %s", testStr, dictMap.String())
	}

	testStr = "b"
	if dictMap.Exists(testStr) {
		log.Fatalf("%s is found in %s", testStr, dictMap.String())
	}
	_, partEx := dictMap.partialExists(testStr)
	if !partEx {
		log.Fatalf("%s is not partially found in %s", testStr, dictMap.String())
	}
}
func TestVisit0(t *testing.T) {
	testDict := []string{
		"bad",
		"cab",
		"baddy",
		"cabby",
		"dab",
	}

	dictMap := NewDictMap(testDict)
	pz := NewPuzzle(2)

	pz.SetDict(dictMap)
	nilFunc := func(string) {}
	pz.StartWorker(nilFunc)

	grid := [][]rune{
		{'a', 'b'},
		{'c', 'd'},
	}
	pz.Grid = grid

	if (pz.getRune(Coord{0, 0}) != 'a') {
		log.Fatal("not a")
	}
	if (pz.getRune(Coord{1, 0}) != 'b') {
		log.Fatal("not b")
	}
	if (pz.getRune(Coord{0, 1}) != 'c') {
		log.Fatal("not c")
	}
	if (pz.getRune(Coord{1, 1}) != 'd') {
		log.Fatal("not d")
	}

	// So that we can self test
	// we do the functions of (pz *Puzzle) StartWorker here
	pz.initWorker()
	var foundWord string
	wordRxd := pz.rxWord(&foundWord)

	err := pz.visit("", Coord{0, 0})
	if err != nil {
		log.Fatal("Visit problems", err)
	}
	err = pz.visit("", Coord{0, 1})
	if err != nil {
		log.Fatal("Visit problems", err)
	}
	<-wordRxd

	if foundWord != "cab" {
		log.Fatalf("Got wrong word, expected cab, got %s\n", foundWord)
	}
	// Having visited, let's make sure the worker can finish
	close(pz.workerComplete) // Close manually as we haven't started a worker
	pz.Shutdown()
}
func TestFileRead(t *testing.T) {
	dic := NewDictMap([]string{}) // TBD fix this
	dic.PopulateFile("wordlist.txt")
	dic.Wait()
	if !dic.Exists("expression") {
		log.Fatal("Expression NE")
	}
	if dic.Exists("expr") {
		log.Fatal("Expr Ex")
	}
	pz := NewPuzzle(2)
	pz.SetDict(dic)
	visitFunc := func(wrd string) {
		log.Println("Found Word", wrd)
	}
	pz.StartWorker(visitFunc)

	grid := [][]rune{
		{'a', 'b'},
		{'c', 'd'},
	}
	pz.Grid = grid

	pz.RunWalk()

	// Having visited, let's make sure the worker can finish
	pz.Shutdown()
}
func TestVisit1(t *testing.T) {
	testDict := []string{
		"bad",
		"cab",
		"baddy",
		"cabby",
		"dab",
	}
	expectedWords := []string{
		"bad",
		"cab",
		"dab",
	}
	expectedMap := make(map[string]struct{})

	for _, wrd := range expectedWords {
		expectedMap[wrd] = struct{}{}
	}
	dictMap := NewDictMap(testDict)
	pz := NewPuzzle(2)

	pz.SetDict(dictMap)

	grid := [][]rune{
		{'a', 'b'},
		{'c', 'd'},
	}
	pz.Grid = grid

	// So that we can self test
	// we do the functions of (pz *Puzzle) StartWorker here
	wrkFunc := func(wrd string) {
		_, ok := expectedMap[wrd]
		if ok {
			delete(expectedMap, wrd)
		} else {
			log.Fatal("Received word not in expected", wrd)
		}
	}
	pz.StartWorker(wrkFunc)
	pz.RunWalk()

	// Having visited, let's make sure the worker can finish
	pz.Shutdown()
	var tooMany bool
	for key := range expectedMap {
		log.Println("Word unfound", key)
		tooMany = true
	}
	if tooMany {
		log.Fatal("Oh Well!")
	}
}

func TestAsPerJs(t *testing.T) {

	data, err := wordlist.Asset("data/wordlist.txt")
	if err != nil {
		t.Fatal("Error", err)
	}
	dic := NewDictMap([]string{})
	dic.PopulateFromBa(data)

	var ra [][]rune
	ra = [][]rune{
		{'b', 'a', 'b', 'd'},
		{'w', 's', 'y', 't'},
		{'a', 'e', 'd', 'g'},
		{'q', 'u', 'p', 'o'},
	}
	sortedResult := NewPuzzleSolve(ra, dic)

	log.Println(sortedResult)
	if len(sortedResult[0]) != 6 {
		t.Fatal("Length of word incorrect:")
	}
}

func TestIndent(t *testing.T) {
	log.Println("\"" + IndentString(1, "bob") + "\"")
	log.Println("\"" + IndentString(1, "fred\nsteve") + "\"")
}
