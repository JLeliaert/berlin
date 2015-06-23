package berlin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputFile     *os.File
	err            error
	outdir         string
	twrite         float64
	outputinterval float64
)

//Initialise the outputdir
func InitOutput(interval float64) {

	outputinterval = interval
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
	writeheader()

}

//checks the error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//writes the head of the outputfile
func writeheader() {
	header := fmt.Sprintf("#t\t<mz>\tB\tDt\n")
	_, err = outputFile.WriteString(header)
	check(err)
}

//Writes the time,average magnetisation and external field to table
func write() {

	if (T==0. || T-twrite >= outputinterval*0.999999999) && outputinterval != 0 {
		//calculate m_avg
		avg := 0.
		totalmoment := 0.
		for i := range Particles {
			avg += Particles[i].mz * Particles[i].weight * Particles[i].msat * Particles[i].V()
			totalmoment += Particles[i].weight * Particles[i].msat * Particles[i].V()
		}
		avg /= totalmoment

		string := fmt.Sprintf("%e\t%v\t%v\t%v\n", T, avg, B_ext(T), Dt)
		_, err = outputFile.WriteString(string)
		check(err)
		twrite = T
	}
}

func Printparticles() {
	for i := range Particles {
		fmt.Println("Particle#: ", i)
		fmt.Println("u_anis: ", Particles[i].u_anis)
		fmt.Println("ku1: ", Particles[i].ku1)
		fmt.Println("r: ", Particles[i].r)
		fmt.Println("weight: ", Particles[i].weight)
		fmt.Println("msat: ", Particles[i].msat)
		fmt.Println()
	}
}
