package main

import (
	githashread "gocmd/testfiles/GitHashRead"

	//gitlog "gocmd/testfiles/GitLog"

	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
)

func main() {
	required := true
	repo, _ := gitpath.Repo_find(gitpath.Get_Os_Dir(), required)

	//fmt.Println(repo.GitDir)

	//66726975a39cdb5babb3a02e8761ecd37e1c7c49 blob

	//064a74c46df49ea6afb685b1004433026c81c152 commit

	//2f4fa26a5a15f2e6f1ea994b3aaebcf664bcf0e5 big tree

	//f04239944489fd19e785c14a39ee71fc5d0fc66f small tree

	//weird tree b769d70e94410966bcdb314ed77f90d43bc41980

	obj := githashread.Object_Read(*repo, "2f4fa26a5a15f2e6f1ea994b3aaebcf664bcf0e5")
	//fmt.Println(obj.Get_Format())

	//sha, _ := githashread.Hash_Object("settings", "blob", *repo)

	//fmt.Printf("SHA: %x \n", sha)

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
	leafs := gitobj.Tree_Parse(obj.Deserialize())
	//leaf := leafs[2]
	//out := leaf.Tree_Leaf_Sort_Key()
	gitobj.Tree_Serialize(leafs)
	//fmt.Println(string(out))

	//log.Log(*repo, obj)
	//gitlog.Log(*repo)
}
