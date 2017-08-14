package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/cbehopkins/boggle"
)

func formatRequest(size int) (string, *regexp.Regexp) {
	ret_str := "Please give input in form:{"
	reg_expr_txt := "^{?"
	comma := ""
	reg_comma := ""

	for i := 0; i < size; i++ {
		ret_str += comma + "[a-z]"
		reg_expr_txt += reg_comma + "([a-z])"
		comma = ","
		reg_comma = "[ ,]?"
	}
	ret_str += "}"
	reg_expr_txt += `}?`
	//fmt.Println("Regex is",reg_expr_txt)
	var ret_expr = regexp.MustCompile(reg_expr_txt)
	return ret_str, ret_expr
}
func getLine(size int, reader *bufio.Reader) []string {
	req_txt, inRe := formatRequest(size)
	var found bool
	ret_array := make([]string, size)
	for !found {
		fmt.Println(req_txt)
		text, _ := reader.ReadString('\n')
		ln := inRe.FindStringSubmatch(text)
		if len(ln) == (size + 1) {
			found = true
			for i := 0; i < size; i++ {
				ret_array[i] = ln[i+1]
			}
		} else {
			fmt.Println("*********** Error in Input******")
		}
	}
	return ret_array

}
func toRuneArray(inArr [][]string) (ret_arr [][]rune) {
	pzLen := len(inArr)
	ret_arr = make([][]rune, pzLen)
	for i := 0; i < pzLen; i++ {
		row := make([]rune, pzLen)
		for j := 0; j < pzLen; j++ {
			r, _ := utf8.DecodeRuneInString(inArr[i][j])
			row[j] = r
		}
		ret_arr[i] = row
	}
	return ret_arr
}

func main() {
	//dic := boggle.NewDictMap([]string{})
	//wg := dic.PopulateFile("wordlist.txt")
	dic := boggle.TopDict()

	reader := bufio.NewReader(os.Stdin)
	var loneNumber = regexp.MustCompile(`^[0-9]+`)

	var size int
	for size == 0 {
		fmt.Print("Enter number: ")
		text, _ := reader.ReadString('\n')
		ln := loneNumber.FindString(text)
		if ln != "" {
			i, err := strconv.ParseInt(ln, 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			size = int(i)
		}
	}
	fmt.Println("Size set to:", size)
	puzzleText := make([][]string, size)
	for i := 0; i < size; i++ {
		fmt.Printf("Line %d:", i)
		new_line := getLine(size, reader)
		puzzleText[i] = new_line
	}
	fmt.Println(puzzleText)
	grid := toRuneArray(puzzleText)

	//	dic.Wait()
	//	pz := boggle.NewPuzzle(size)
	//	pz.SetDict(dic)
	//	pz.Grid = grid
	pz := dic.NewPuzzle(size, grid)

	wrds_found := make(map[string]struct{})
	wrkFunc := func(wrd string) {
		//fmt.Println("Found Word", wrd)
		wrds_found[wrd] = struct{}{}
	}
	pz.StartWorker(wrkFunc)
	pz.RunWalk()
	pz.Shutdown()

	for wrd, _ := range wrds_found {
		fmt.Println("Found Word", wrd)
	}
}
