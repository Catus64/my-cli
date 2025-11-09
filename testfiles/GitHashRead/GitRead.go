package GitHashRead

//
//Cross compatibality with windows is yet to be implemented
//

import (
	"bytes"
	"compress/zlib"
	"fmt"
	gitobj "gocmd/testfiles/GitObject"
	gitpath "gocmd/testfiles/Gitrepostruct"
	"io"
	"os"
)

// reading blob data to print them out
func Read_Blob(path string) *[]byte {

	//check whether file exist
	_, err := os.Stat(path)
	if err != nil {
		panic(fmt.Sprintf("%v does not exist/error", path))
	}

	fmt.Println("File exist")

	//read file
	filedata, err := os.ReadFile(path)
	if err != nil {
		panic("File cannot be read")
	}

	//decompressing via zlib
	r, err := zlib.NewReader(bytes.NewReader(filedata))
	if err != nil {
		panic("not zlib format")
	}
	defer r.Close()

	//reading decompressed data
	decompressed, err := io.ReadAll(r)
	if err != nil {
		panic("read fail")
	}

	return &decompressed
}

func Object_Read(repo gitpath.GitRepository, sha string) gitobj.GitObject {

	//splitting strings for path check
	first := sha[:2]
	rest := sha[2:]

	//getting obj_path
	obj_path := gitpath.Repo_Path(repo, "objects", first, rest)
	fmt.Println("Reading: ", obj_path)

	//reading blob with zlib
	decompressed := Read_Blob(obj_path)

	fmt.Println("Decompressed: ", decompressed)

	//make obj based on format
	obj := gitobj.MakeGitObj(*decompressed)
	//fmt.Println(string(obj.Deserialize()))

	return obj
}
