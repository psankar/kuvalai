package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Message struct {
	Output string
}

type Code struct {
	Code string
}

func compilerService(w http.ResponseWriter, req *http.Request) {

	/* Get the source code */
	err := req.ParseForm()
	if err != nil {
		panic(err)
		return
	}
	var c Code
	err = json.Unmarshal([]byte(req.Form["a"][0]), &c)
	if err != nil {
		panic(err)
		return
	}

	/* create a temporary file with the passed sourcecode as the content */
	var tmpFile *os.File
	tmpFile, err = ioutil.TempFile("", "")
	if err != nil {
		panic(err)
		return
	}
	tmpFileName := tmpFile.Name()
	_, err = tmpFile.WriteString(c.Code)
	if err != nil {
		panic(err)
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
		fmt.Println("Compilation error")
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

	var result string
	result = out.String()

	/* Format the output in case of error, replace the filename to sample.cpp */
	if err != nil {
		result = strings.Replace(result, tmpFileName, "sample.cpp", -1)
	}

	m := Message{result}

	var b []byte
	b, err = json.Marshal(m)
	if err != nil {
		panic(err)
		return
	}
	w.Write(b)

	return
}

type slide struct {
	Contents    string
	Code        string
	SlideNumber int
}

func kuvalaiServer(w http.ResponseWriter, req *http.Request) {

	/* Parse the input file and tokenize to get a slide array */
	file, err := os.Open("kuvalai-example.md")
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()

	output := bufio.NewWriter(w)
	fmt.Fprintln(output, `<html>

<head>
    <meta charset="UTF-8">
    <title>Kuvalai</title>
    <link rel="stylesheet" type="text/css" href="resources/css/bootstrap.min.css">
    <link rel="stylesheet" type="text/css" href="resources/css/kuvalai.css">
    <script type="text/javascript" src="resources/js/kuvalai.js"></script>
    <script type="text/javascript" src="resources/js/jquery.min.js"></script>
    <script type="text/javascript" src="resources/js/bootstrap.min.js"></script>
</head>

<body>
    <div>
        <div id="m" class="carousel slide">
            <!-- Carousel items -->
            <div class="carousel-inner">`)

	newSlide := slide{"", "", 1}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "---") {
			/* Write slide into the file as html */

			var tmplString string
			tmplString = "\n" + "<!-- Slide {{.SlideNumber}} -->" + "\n"
			if newSlide.SlideNumber == 1 {
				tmplString += `                <div class="active item">` + "\n"
			} else {
				tmplString += `                <div class="item">` + "\n"
			}

			tmplString += `                    <div class="row" style="margin:1px">` + "\n"
			if len(newSlide.Code) > 0 {
				tmplString += `                        <div class="col-xs-8">
                            <div class="row">
			    <textarea id="kuvCode{{.SlideNumber}}" class="code-area">{{.Code}}</textarea>
                            </div>
                            <div class="row">
			    <textarea id="kuvOutput{{.SlideNumber}}" class="code-output" readonly>Press Execute to run the above program</textarea>
                            </div>
                            <button class="btn btn-primary" style="position:relative;top:-50px;right:-90%" onclick="executeCode({{.SlideNumber}})">Execute</button>
                        </div>` + "\n"
			}
			tmplString += `<div class="column">
{{.Contents}}
                        </div>
                    </div>
                </div>` + "\n"

			var tmpl *template.Template
			tmpl, err = template.New("test").Parse(tmplString)
			if err != nil {
				panic(err)
				return
			}
			err = tmpl.Execute(output, newSlide)
			if err != nil {
				panic(err)
				return
			}
			fmt.Fprint(output, "\n")

			newSlide.Contents = ""
			newSlide.Code = ""
			newSlide.SlideNumber++

		} else if strings.HasPrefix(line, ".code") {
			/* TODO Read the file pointed by line and get its
			 * contents here */
			filePath := strings.TrimPrefix(line, ".code ")
			var code []byte
			code, err = ioutil.ReadFile(filePath)
			if err != nil {
				panic(err)
				return
			}
			newSlide.Code = string(code)

		} else if strings.HasPrefix(line, "#") {
			newSlide.Contents += "<h1>" + strings.TrimPrefix(line, "#") + "</h1>\n"
		} else {
			newSlide.Contents += line
			newSlide.Contents += "</br>\n"
		}
	}

	fmt.Fprint(output, `            </div>
            <!-- Carousel nav -->
            <button class="previous" href="#m" data-slide="prev">Prev</button>
            <button class="next" href="#m" data-slide="next">Next</button>
            <!-- Carousel indicators -->
            <ol class="carousel-indicators">
                <li data-target="#m" data-slide-to="0" class="active"></li>` + "\n")

	/* The -1 is needed below as li carousel-indicators begin at 0 not 1 */
	for i := 1; i < newSlide.SlideNumber-1; i++ {
		fmt.Fprintf(output, "                <li data-target=\"#m\" data-slide-to=\"%d\"></li>\n", i)
	}

	fmt.Fprint(output, `            </ol>
	</div>
    </div>
</body>

</html>`)
	output.Flush()
}

func main() {
	http.HandleFunc("/", kuvalaiServer)
	http.HandleFunc("/Compile", compilerService)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe("localhost:8080", nil)
}
