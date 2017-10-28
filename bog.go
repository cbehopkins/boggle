package boggle

// TopDict wrapper funciton around the boggle functionality
func TopDict() (dic *DictMap) {
	dic = NewDictMap([]string{})
	dic.PopulateFile("wordlist.txt")
	return dic
}
