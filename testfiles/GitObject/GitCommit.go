package GitObjLib

import (
	"bytes"
	"fmt"
	"strings"
)

type GitCommit struct {
	GitObjectData
	KvlmDict
	format []byte
}

type KvlmDict struct {
	Dict map[string][]byte
}

func (c *GitCommit) Serialize() *[]byte {
	ret := Kvlm_Serialize(c.KvlmDict)
	return &ret
}

func (commit *GitCommit) Deserialize() []byte {
	if commit.data == nil {
		return Kvlm_Serialize(commit.KvlmDict)
	}
	// clear dict before re-parsing to prevent doubling
	temp := make(map[string][]byte)
	kvlm := KvlmDict{Dict: temp}
	start := 0
	commit.KvlmDict = Kvlm_Parse(commit.data, start, kvlm)
	return commit.data
}

func (c *GitCommit) Get_Format() string {
	return string(c.format)
}

func Kvlm_Parse(data []byte, start int, dict KvlmDict) KvlmDict {

	space := bytes.IndexByte(data[start:], ' ') + start
	newline := bytes.IndexByte(data[start:], '\n') + start

	// base case: no more key-value pairs
	if space-start < 0 || newline < space {
		dict.Dict["data"] = data[start:]
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
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)

	_, ok := dict.Dict[string(key)]
	if ok {
		// key exists — append with null separator (e.g. multiple parents)
		dict.Dict[string(key)] = append(dict.Dict[string(key)], 0x00)
		dict.Dict[string(key)] = append(dict.Dict[string(key)], valueCopy...)
	} else {
		// first time seeing this key
		dict.Dict[string(key)] = valueCopy
	}

	return Kvlm_Parse(data, end+1, dict)
}

func Kvlm_Serialize(kvlm KvlmDict) []byte {

	var ret []byte

	keys := []string{"tree", "parent", "author", "committer"}
	for _, key := range keys {
		value, ok := kvlm.Dict[key]
		if !ok {
			continue // skip missing optional fields like parent
		}
		ret = append(ret, []byte(key)...)
		ret = append(ret, ' ')
		ret = append(ret, value...)
		ret = append(ret, '\n')
	}
	ret = append(ret, '\n')
	ret = append(ret, kvlm.Dict["data"]...)

	fmt.Println(string(ret))

	return ret
}

func GetKvlmValues(dict map[string][]byte, key string) []string {
	val, ok := dict[key]
	if !ok {
		return nil
	}
	// split on null byte separator
	parts := bytes.Split(val, []byte{0x00})
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(string(p))
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
