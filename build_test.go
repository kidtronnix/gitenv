package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	pwd, err := ioutil.TempDir("", "env")
	if err != nil {
		t.Fatal(err)
	}

	gitenv := `{
		"links": {
        	".": "src/github.com/ElementAI/gitenv"
	    },
    	"repositories": [
        	{
    	        "url": "https://github.com/pkg/errors.git",
	            "dir": "src/github.com/pkg/errors",
        	    "commit": "3866ebc348c54054262feae422da428fe6cf147d"
	        }
	    ]}`

	path := filepath.Join(pwd, ".gitenv")
	ioutil.WriteFile(path, []byte(gitenv), 0700)

	if err := build(path, []string{}); err != nil {
		t.Fatal(err)
	}

	ok := 0

	filepath.Walk(filepath.Join(pwd, ".env"), func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			ok++
		}

		if strings.HasSuffix(path, "src/github.com/ElementAI/gitenv") && f.Mode()&os.ModeSymlink != 0 {
			ok++
		}

		return err
	})

	if ok != 8 {
		t.Fail()
	}
}
