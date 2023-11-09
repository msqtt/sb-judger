package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLangConfs(t *testing.T) {
	b, err := ioutil.ReadFile("./lang_test.json")
	require.NoError(t, err)
	require.NotEmpty(t, b)

	lc, err := GetLangConfs("./lang_test.json")
	require.NoError(t, err)
	require.NotEmpty(t, lc)

	b2, err2 := json.Marshal(lc)
	require.NoError(t, err2)

	require.NotEmpty(t, b2)
	var b1 bytes.Buffer
	err = json.Compact(&b1, b)
	require.NoError(t, err)

	require.ElementsMatch(t, b1.Bytes(), b2)
}
