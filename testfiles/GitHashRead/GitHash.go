package GitHashRead

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"os"
	"strconv"
)

// Change rawblob byte to raw blob byte with headings
func BuildGitObjectToWrite(obj gitobj.GitObject) []byte {
	size := []byte(strconv.Itoa(len(obj.Deserialize())))

	result := append([]byte(obj.Get_Format()), ' ')
	result = append(result, size...)
	result = append(result, 0x00)
	result = append(result, obj.Deserialize()...)

	return result
}

// writing object to the repo if not yet exist
func Object_Write(obj gitobj.GitObject, repo *gitpath.GitRepository) ([]byte, error) {

	result := BuildGitObjectToWrite(obj)

	sha := sha1.Sum(result)

	hexstring := fmt.Sprintf("%x", sha)

	//fmt.Printf("SHA: %x \n", sha)

	if repo != nil {
		mkdir := true
		location_1 := hexstring[:2]
		location_2 := hexstring[2:]
		path := gitpath.Repo_File(*repo, mkdir, "objects", location_1, location_2)

		_, err := os.Stat(path)
		if err == nil {
			return sha[:], nil
		}
		if os.IsNotExist(err) {
			file, err := os.Create(path)
			if err != nil {
				panic(fmt.Sprintf("%v", err))
			}
			defer file.Close()

			writer := zlib.NewWriter(file)
			defer writer.Close()

			_, err = writer.Write(result)
			if err != nil {
				panic(fmt.Sprintf("%v", err))
			}
			writer.Close()
		}

	}
	//convert byte[n] to slice(byte[])
	return sha[:], nil
}

// wrapper command to call the hashing and writing function
func Hash_Object(file_path string, format string, repo gitpath.GitRepository) ([]byte, error) {

	data, err := os.ReadFile(file_path)
	if err != nil {
		panic(fmt.Sprintf("error opening file %v", file_path))
	}

	obj := gitobj.MakeGitObjWithFormat(data, format)

	sha, _ := Object_Write(obj, &repo)

	return sha, nil
}
