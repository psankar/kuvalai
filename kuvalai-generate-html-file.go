package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type slide struct {
	Contents    string
	Code        string
	SlideNumber int
}

func main() {

	/* Parse the input file and tokenize to get a slide array */
	file, err := os.Open("kuvalai-example.md")
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()

	outfile, err := os.Create("kuvalai-example.html")
	if err != nil {
		panic(err)
		return
	}
	defer outfile.Close()

	output := bufio.NewWriter(outfile)
	fmt.Fprintln(output, `<html>

<head>
    <meta charset="UTF-8">
    <title>Kuvalai</title>
    <link rel="stylesheet" type="text/css" href="css/bootstrap.min.css">
    <link rel="stylesheet" type="text/css" href="css/kuvalai.css">
    <script type="text/javascript" src="js/kuvalai.js"></script>
    <script type="text/javascript" src="js/jquery.min.js"></script>
    <script type="text/javascript" src="js/bootstrap.min.js"></script>
</head>

<body>
    <div>
        <div id="m" class="carousel slide">
            <!-- Carousel indicators -->
            <ol class="carousel-indicators">
                <li data-target="#m" data-slide-to="0" class="active"></li>
                <li data-target="#m" data-slide-to="1"></li>
                <li data-target="#m" data-slide-to="2"></li>
            </ol>
            <!-- Carousel items -->
            <div class="carousel-inner">`)

	newSlide := slide{"", "", 1}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "---") {
			fmt.Println("End of a slide reached", newSlide.Contents)
			/* Write slide into the file as html */

			var itemdiv string
			if newSlide.SlideNumber == 1 {
				itemdiv = `                <div class="active item">`
			} else {
				itemdiv = `                <div class="item">`
			}

			var tmpl *template.Template
			tmpl, err = template.New("test").Parse(itemdiv + `                    <div class="row" style="margin:1px">
                        <div class="col-xs-8">
                            <div class="row">
			    <textarea id="kuvCode{{.SlideNumber}}" class="code-area">{{.Code}}</textarea>
                            </div>
                            <div class="row">
			    <textarea id="kuvOutput{{.SlideNumber}}" class="code-output" readonly>Press Execute to run the above program</textarea>
                            </div>
                            <button class="btn btn-primary" style="position:relative;top:-50px;right:-90%" onclick="executeCode({{.SlideNumber}})">Execute</button>
                        </div>
                        <div class="column">
{{.Contents}}
                        </div>
                    </div>
                </div>`)
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
        </div>
    </div>
</body>

</html>`)
	output.Flush()
}
