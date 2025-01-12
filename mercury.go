package main

import (
	"fmt"
	errors "go/scanner"
	positions "go/token"
	"os"
	"sort"

	"github.com/arneph/mercury/logic/simulation"
	"github.com/arneph/mercury/logic/text"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Expected path for input file.")
		os.Exit(1)
		return
	}
	path := os.Args[1]
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not read path: %v\n", err)
		os.Exit(1)
		return
	}
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile(path, fileSet.Base(), len(src))
	file.SetLinesForContent(src)

	system, errs := text.BuildFromFile(file, src)
	if errs.Len() > 0 {
		errs.RemoveMultiples()
		errors.PrintError(os.Stderr, errs)
		return
	}
	testNames := make([]string, 0, len(system.Tests))
	for name := range system.Tests {
		testNames = append(testNames, name)
	}
	sort.Strings(testNames)
	for _, name := range testNames {
		test := system.Tests[name]
		fmt.Printf("test %-20s ", name)
		errs := simulation.RunTest(test, file)
		if errs.Len() == 0 {
			fmt.Println("PASS")
		} else {
			fmt.Println("FAIL")
			errors.PrintError(os.Stderr, errs)
		}
	}
}
