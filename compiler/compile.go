package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	goFilePath := strings.Replace(os.Args[1]+".go", ".sx", "", -1)
	shexFileRaw, _ := os.ReadFile("shex.go")
	sourceFile, _ := os.ReadFile(os.Args[1])
	shexFile := strings.Replace(string(shexFileRaw), "rawSource, _ := os.ReadFile(os.Args[1])", "", -1)
	shexFile = strings.Replace(shexFile, "\"os\"", "", -1)
	shexFile = strings.Replace(shexFile, "string(rawSource)", "`"+string(sourceFile)+"`", -1)
	os.WriteFile(goFilePath, []byte(shexFile), 0644)

	cmd := exec.Command("go", "build", goFilePath)
	cmd.Run()

	os.Remove(goFilePath)
}
