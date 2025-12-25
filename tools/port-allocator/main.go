package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var baseDir string
	flag.StringVar(&baseDir, "base-dir", ".", "Base directory for the project")
	flag.Parse()

	allocator := NewPortAllocator(baseDir)
	if err := allocator.AllocatePorts(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Port allocation completed successfully")
	os.Exit(0)
}

