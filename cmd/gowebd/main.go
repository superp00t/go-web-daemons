package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/superp00t/etc"
	"github.com/superp00t/etc/yo"
)

func main() {
	yo.Stringf("o", "output", "output directory", "dist")
	yo.AddSubroutine("gen", []string{"program"}, "compile a Go Webdaemon for Web!", _main)
	yo.Init()
}

func _main(args []string) {
	if len(args) == 0 {
		yo.Fatal("you must include a program to build")
	}

	// parse program name

	program := args[0]
	importPath := program
	s := strings.Split(program, "/")
	head := s[len(s)-1]

	withoutExt := strings.Split(head, ".")
	if len(withoutExt) == 2 && withoutExt[1] == "go" {
		program = withoutExt[0]
	} else {
		program = head
	}

	dist := yo.StringG("o")
	pth := etc.ParseSystemPath(dist)
	if !pth.IsExtant() {
		yo.Ok("Creating directory", pth.Render())
		pth.Mkdir()
	}

	output := etc.NewBuffer()

	goBuildCmd := exec.Command("go", "build", "-o", pth.Concat(program+".wasm").Render(), importPath)
	goBuildCmd.Env = append(
		os.Environ(),
		"GOOS=js",
		"GOARCH=wasm",
	)
	goBuildCmd.Stderr = output
	goBuildCmd.Stdout = output
	if err := goBuildCmd.Run(); err != nil {
		fmt.Println(output.String())
		yo.Fatal(err)
	}

	sharedWorker := fmt.Sprintf(sharedWorkerTpl, program)

	ioutil.WriteFile(pth.Concat("wasm_exec.js").Render(), []byte(loadWasmExec()), 0700)
	ioutil.WriteFile(pth.Concat(fmt.Sprintf("%s-sharedworker.js", program)).Render(), []byte(sharedWorker), 0700)

	yo.Ok("Successfully built @", pth.Render())
}

func loadWasmExec() string {
	if os.Getenv("GOROOT") == "" {
		yo.Fatal("you need a GOROOT to run this command")
	}

	path := etc.ParseSystemPath(os.Getenv("GOROOT")).Concat("misc", "wasm", "wasm_exec.js")

	b, err := ioutil.ReadFile(path.Render())
	if err != nil {
		yo.Fatal(err)
	}

	return string(b)
}

const sharedWorkerTpl = `self.importScripts(
  "wasm_exec.js"
)

global.pendingConnections = [];

self.onconnect = function(e) {
	global.pendingConnections.push(e);
}

const go = new Go();
WebAssembly.instantiateStreaming(
  fetch("%s.wasm"), go.importObject
)
.then((result) => {
  go.run(result.instance);
});`
