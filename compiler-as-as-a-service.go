package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

type Message struct {
	Output string
}

type Code struct {
	Code string
}

func compilerService(w http.ResponseWriter, req *http.Request) {

	fmt.Println("New request came")

	/* Get the source code */
	err := req.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
		fmt.Println(err)
		return
	}
	var c Code
	err = json.Unmarshal([]byte(req.Form["a"][0]), &c)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return
	}
	fmt.Println(c.Code)

	/* create a temporary file with the passed sourcecode as the content */
	var tmpFile *os.File
	tmpFile, err = ioutil.TempFile("", "")
	if err != nil {
		fmt.Println("Error creating temporary file")
		fmt.Println(err)
		return
	}
	tmpFileName := tmpFile.Name()
	fmt.Println("Temporary file created is: " + tmpFileName)
	_, err = tmpFile.WriteString(c.Code)
	if err != nil {
		fmt.Println("Error writing to the temporary file")
		fmt.Println(err)
		return
	}

	/* compile+execute the source code and get the output */
	/* there is an assumption here that the tmpFileName.out will be a new
	 * file name */
	cmd := exec.Command("g++", "-x", "c++", "-o", tmpFileName+".out", tmpFileName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running the compiler command")
		fmt.Println(err)
	} else {
		cmd = exec.Command(tmpFileName + ".out")
		cmd.Stdout = &out
		cmd.Stderr = &out
		_ = cmd.Run()
	}

	/* Remove the temporary files created */
	_ = os.Remove(tmpFileName)
	_ = os.Remove(tmpFileName + ".out")

	/* Return the output/error of the cpp program to the caller */
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	m := Message{out.String()}

	var b []byte
	b, err = json.Marshal(m)
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return
	}
	w.Write(b)

	return
}

func main() {
	http.HandleFunc("/", compilerService)
	http.ListenAndServe("localhost:8080", nil)
	fmt.Println("Hello World")
}
