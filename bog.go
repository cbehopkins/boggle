package boggle

func TopDict() (dic *DictMap) {
	dic = NewDictMap([]string{})
	dic.PopulateFile("wordlist.txt")
	return dic
}
