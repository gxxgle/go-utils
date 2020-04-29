package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gxxgle/go-utils/json"
	"github.com/gxxgle/go-utils/path"
)

// Init config form json file
func Init(config interface{}, paths ...string) error {
	var file string

	if len(paths) > 0 {
		file = paths[0]
	}

	if len(file) == 0 {
		dir, _ := filepath.Abs(path.CurrentDir())
		bin := path.CurrentFilename()
		file = fmt.Sprintf("%s/%s.json", dir, bin)
	}

	fh, err := os.Open(file)
	if err != nil {
		return err
	}

	defer fh.Close()
	bs, err := ioutil.ReadAll(fh)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, config)
}
