package berlin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputFile *os.File
	err        error
	outdir     string
)

//Initialise the outputdir
func initOutput() {

	// make and clear output directory
	fname := os.Args[0]
	f2name := strings.Split(fname, "/") // TODO: use path.Split?
	outdir = fmt.Sprint(f2name[len(f2name)-1], ".out")
	os.Mkdir(outdir, 0775)
	dir, err3 := os.Open(outdir)
	files, _ := dir.Readdir(1)
	// clean output dir, copied from mumax
	if len(files) != 0 {
		filepath.Walk(outdir, func(path string, i os.FileInfo, err error) error {
			if path != outdir {
				os.RemoveAll(path)
			}
			return nil
		})
	}

	if err3 != nil {
		panic(err3)
	}

	outputFile, err = os.Create(outdir + "/table.txt")
	check(err)

}

//checks the error
func check(e error) {
	if e != nil {
		panic(e)
	}
}
