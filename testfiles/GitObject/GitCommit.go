package GitObjLib

import (
	"bytes"
	"fmt"
)

type GitCommit struct {
	GitObjectData
	KvlmDict
	format []byte
}

type KvlmDict struct {
	Dict map[string][]byte
}

func (blob *GitCommit) Serialize() *[]byte {
	return nil
}

func (commit *GitCommit) Deserialize() []byte {
	temp := make(map[string][]byte)
	kvlm := KvlmDict{
		Dict: temp,
	}
	start := 0
	commit.KvlmDict = Kvlm_Parse(commit.data, start, kvlm)
	return nil
}

func (blob *GitCommit) Get_Format() string {
	return string(blob.format)
}

func Kvlm_Parse(data []byte, start int, dict KvlmDict) KvlmDict {

	space := bytes.IndexByte(data[start:], ' ') + start
	newline := bytes.IndexByte(data[start:], '\n') + start

	if space < 0 || newline < space {
		//end recursion
		dict.Dict["data"] = data[start+1:]
		return dict
	}

	key := data[start:space]
	end := start

	for {
		end_offset := bytes.IndexByte(data[end+1:], '\n')
		end = end + end_offset + 1
		if data[end+1] != ' ' {
			break
		}
	}

	value := data[space+1 : end]

	_, ok := dict.Dict[string(key)]
	if ok {
		dict.Dict[string(key)] = append(dict.Dict[string(key)], 0x00)
		dict.Dict[string(key)] = append(dict.Dict[string(key)], value...)
	} else {
		dict.Dict[string(key)] = value
	}

	//fmt.Println("key: ", string(key), "|value: ", string(value))
	return Kvlm_Parse(data, end+1, dict)
}

func Kvlm_Serialize(kvlm KvlmDict) []byte {

	var ret []byte

	for key, value := range kvlm.Dict {
		if key == "data" {
			continue
		}
		ret = append(ret, []byte(key)...)
		ret = append(ret, ' ')
		ret = append(ret, value...)
		ret = append(ret, '\n')
	}
	ret = append(ret, kvlm.Dict["data"]...)

	fmt.Println(string(ret))

	return nil
}
