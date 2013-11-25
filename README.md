#Kuvalai குவளை

A presentation software primarily focused on:
* Simplicity - Editable via $EDITOR
* Developer oriented - Ability to show and run code from within the slides
* Independent - Not tied to any webservice (ala prezi, slideshare etc) and run
offline using the compiler on the target machine
* Minimal

#License

All source code is licensed under Creative Commons Zero License.
More information at:    http://creativecommons.org/publicdomain/zero/1.0/
Full license text at:   http://creativecommons.org/publicdomain/zero/1.0/legalcode

The following files are copied from the bootstrap project and please refer to
the bootstrap project page for their license.
* css/bootstrap.min.css
* js/bootstrap.min.js
* js/jquery.min.js

#Usage
+ Install Go from http://golang.org
+ go run compiler-as-a-service.go
+ go run kuvalai-generate-html-file.go ; # Generates kuvalai-example.html from kuvalai-example.md
+ Open kuvalai-example.html in your browser
+ Edit the source in the first slide as you wish and press the Execute button.
The program that you have typed will be executed and the output will be printed
in the textarea. If there is a compilation error, the compiler error too will be
printed the same place.
