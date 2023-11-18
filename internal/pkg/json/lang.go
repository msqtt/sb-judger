package json

type Compile struct {
	Cmd   string   `json:"cmd"`
	Flags []string `json:"flags"`
}

type LanguageConfig struct {
	Pids    uint32 `json:"pids"`
	Src     string `json:"src"`
	Out     string `json:"out"`
	Compile `json:"compile"`
	Run     []string `json:"run"`
	TrimMsg []string `json:"trimMsg"`
}
