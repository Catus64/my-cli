package main

import (
	"fmt"
	githashread "gocmd/testfiles/GitHashRead"

	//gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
)

func main() {
	required := true
	repo, _ := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)

	//fmt.Println(repo.GitDir)

	//66726975a39cdb5babb3a02e8761ecd37e1c7c49 blob

	//064a74c46df49ea6afb685b1004433026c81c152 commit

	//b769d70e94410966bcdb314ed77f90d43bc41980 tree

	//be0fc874cd0cf5c8764f5ffac306f254fdba477f big tree

	//obj := githashread.Object_Read(*repo, "b672c3cd024a603265df11f330597f38beb2bbf6")
	//fmt.Println(obj.Get_Format())

	sha, _ := githashread.Hash_Object("settings", "blob", *repo)

	fmt.Printf("SHA: %x \n", sha)

	/*
		temp := make(map[string][]byte)
		kvlm := gitobj.KvlmDict{
			Dict: temp,
		}
		s := gitobj.Kvlm_Parse(obj.Deserialize(), 0, kvlm)
		//fmt.Println(string(s.Dict["parent"]))

		gitobj.Kvlm_Serialize(s)
	*/

	//fmt.Println(obj.Get_Format())
	//gitobj.Tree_Parse_One(obj.Deserialize(), 0)
}
