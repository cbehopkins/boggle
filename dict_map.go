package boggle

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

// DictMap is a heirachical tree
// that represents all supplied words as a tree of runes
// each node can be a complete word
type DictMap struct {
	// is (the path that got us here) a word in itself?
	isword bool

	// Additionally the other words that use this as a root
	// what is the next rune in that word
	currentRunes map[rune]*DictMap
	wg           sync.WaitGroup
}

// NewDictMap return a new dictionary
func NewDictMap(input []string) *DictMap {
	itm := new(DictMap)
	itm.currentRunes = make(map[rune]*DictMap)
	if len(input) > 0 {
		itm.Populate(input)
	}
	return itm
}
func dumpIndent(level int) string {
	var retStr string

	for i := 0; i < level; i++ {
		retStr += " "
	}
	return retStr
}
// IndentString return a new string indented as required
func IndentString(level int, inTxt string) string {
	var retStr string
	nl := ""
	for _, tmpStr := range strings.Split(inTxt, "\n") {
		retStr += nl + dumpIndent(1) + tmpStr
		nl = "\n"
	}
	return retStr
}
// NewPuzzle return a new puzzle struicture of specified size
func (dic *DictMap) NewPuzzle(size int, grid [][]rune) (pz *Puzzle) {
	dic.Wait()
	pz = NewPuzzle(size)
	pz.SetDict(dic)
	pz.Grid = grid
	return pz
}
func (dic DictMap) String() string {
	var retStr string
	if dic.isword {
		retStr += "Word"
	}
	for key, value := range dic.currentRunes {
		retStr += "\n" + string(key) + IndentString(1, value.String())
	}
	return retStr
}
// Populate in batches of strings
func (dic *DictMap) Populate(input []string) {
	for _, word := range input {
		//log.Println("Populating Word", input)
		dic.Add(word)
	}
}
// Readln standard Readln interface
// Read from any reader and pull in a full line
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
func (dic *DictMap) populateFile(filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(filename)
	if err == os.ErrNotExist {
		return
	} else if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("error opening file: %T\n", err)
			return
		}
			fmt.Printf("error opening file: %T\n", err)
			os.Exit(1)
			return
			}
	defer f.Close()
	r := bufio.NewReader(f)
	for s, e := Readln(r); e == nil; s, e = Readln(r) {
		s = strings.TrimSpace(s)
		comment := strings.HasPrefix(s, "//")
		comment = comment || strings.HasPrefix(s, "#")
		if comment {
			continue
		}
		if s == "" {
			continue
		}

		dic.Add(s)
	}
	fmt.Println("Finished reading File")

}
// PopulateFile Populate the dictionary from a file of words
func (dic *DictMap) PopulateFile(filename string) *sync.WaitGroup {

	dic.wg.Add(1)
	go dic.populateFile(filename, &dic.wg)
	return &dic.wg
}
// Add a word to the dictionary
// This builds the rune by rune tree
func (dic *DictMap) Add(word string) {
	if len(word) == 0 {
		dic.isword = true
		return
	}

	// Get the first rune in the string
	// and record how many bytes we used for this
	r, length := utf8.DecodeRuneInString(word)
	// remove the appropriate number of bytes from the string
	newWord := word[length:]

	// Do we already have a rune map for that?
	dependantDict, ok := dic.currentRunes[r]
	if !ok {
		// create as needed
		dependantDict = NewDictMap([]string{})

	} else {
		// safty check
		if dependantDict == nil {
			log.Fatal("Uninitialised Dict")
		}
	}
	// add the remaiins to the
	dependantDict.Add(newWord)
	dic.currentRunes[r] = dependantDict
}
// Wait until all the processing has completed
func (dic *DictMap) Wait() {
	dic.wg.Wait()
}
// Exists returns true is the word exists in the dictionary
func (dic DictMap) Exists(inTxt string) bool {
	isWord, _ := dic.partialExists(inTxt)
	return isWord
}

func (dic DictMap) partialExists(inTxt string) (isword, partial bool) {
	if len(inTxt) == 0 {
		return dic.isword, true
	}

	r, length := utf8.DecodeRuneInString(inTxt)
	newWord := inTxt[length:]
	dependantDict, ok := dic.currentRunes[r]

	if !ok {
		return false, false
	}
		return dependantDict.partialExists(newWord)
	}
