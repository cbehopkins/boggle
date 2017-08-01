package boggle

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

type DictMap struct {
	// is (the path that got us here) a word in itself?
	isword bool

	// Additionally the other words that use this as a root
	// what is the next rune in that word
	currentRunes map[rune]*DictMap
}

func NewDictMap(input []string) *DictMap {
	itm := new(DictMap)
	itm.currentRunes = make(map[rune]*DictMap)
	if len(input) > 0 {
		itm.Populate(input)
	}
	return itm
}
func dumpIndent(level int) string {
	var ret_str string

	for i := 0; i < level; i++ {
		ret_str += " "
	}
	return ret_str
}
func IndentString(level int, inTxt string) string {
	var ret_str string
	nl := ""
	for _, tmpStr := range strings.Split(inTxt, "\n") {
		ret_str += nl + dumpIndent(1) + tmpStr
		nl = "\n"
	}
	return ret_str
}
func (dic DictMap) String() string {
	var ret_str string
	if dic.isword {
		ret_str += "Word"
	}
	for key, value := range dic.currentRunes {
		ret_str += "\n" + string(key) + IndentString(1, value.String())
	}
	return ret_str
}
func (dic *DictMap) Populate(input []string) {
	for _, word := range input {
		//log.Println("Populating Word", input)
		dic.Add(word)
	}
}
func Readln(r *bufio.Reader) (string, error) {
  var (
      isPrefix bool  = true
          err      error = nil
              line, ln []byte
                )
                  for isPrefix && err == nil {
                      line, isPrefix, err = r.ReadLine()
                          ln = append(ln, line...)
                            }
                              return string(ln), err
                              }

func (dic *DictMap) PopulateFile(filename string) {
	f, err := os.Open(filename)
	if err == os.ErrNotExist {
		return
	} else if err != nil {
		if os.IsNotExist(err) {
			return
		} else {
			fmt.Printf("error opening file: %T\n", err)
			os.Exit(1)
			return
		}
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
}
func (dic *DictMap) Add(word string) {
	if len(word) == 0 {
		dic.isword = true
		return
	}

	// Get the first rune in the string
	// and record how many bytes we used for this
	r, length := utf8.DecodeRuneInString(word)
	// remove the appropriate number of bytes from the string
	new_word := word[length:]

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
	dependantDict.Add(new_word)
	dic.currentRunes[r] = dependantDict
}

func (dic DictMap) Exists(inTxt string) bool {
	isWord, _ := dic.partialExists(inTxt)
	return isWord
}

func (dic DictMap) partialExists(inTxt string) (isword, partial bool) {
	if len(inTxt) == 0 {
		return dic.isword, true
	}

	r, length := utf8.DecodeRuneInString(inTxt)
	new_word := inTxt[length:]
	dependantDict, ok := dic.currentRunes[r]

	if !ok {
		return false, false
	} else {
		return dependantDict.partialExists(new_word)
	}
}
