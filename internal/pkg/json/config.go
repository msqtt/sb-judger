package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const langConfPath = "configs/lang.json"

type LangConfMap map[string]*LanguageConfig

// GetLangConfs reads lang.json file from named confPath then return arrays of language conf struct
// and an error, if any.
// if you leave named confPath empty, it will be default path.
func GetLangConfs(confPath string) (LangConfMap, error) {
	if confPath == "" {
		confPath = langConfPath
	}
	var res LangConfMap
	err := getObjectFromJson(confPath, &res)
	return res, err
}

func getObjectFromJson(path string, res any) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file from path %s: %w", path, err)
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	return nil
}
