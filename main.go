package main

import (
	"fmt"
	tpl "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	fp "path/filepath"
	"strconv"
	"strings"
	"time"
)

const template = `
<!DOCTYPE html>

<html>
	<head>
		<title>good morning from {{.Name}}</title>
			<style>
				body{
				background-image: url({{.Picture}});
				background-postion: center center;
				background-size: cover;
				background-repeat: no-repeat;
				-moz-background-size: cover;
				-webkit-background-size: cover;
				-o-background-size: cover;
				color: white;
				text-shadow: 3px 3px 8px black;
				}
			</style>
	</head>
	<body>
		<h1>good morning from {{.Name}}</h1>
	</body>
</html>

`

type panorama struct {
	name  string
	photo string
}

type locmap map[int]panorama

var lm locmap

func (lm *locmap) Name() string {
	return (*lm)[lm.offset()].name
}

func (lm *locmap) Picture() string {
	return (*lm)[lm.offset()].photo
}

func (lm *locmap) offset() int {
	//get timezone where it is currently between 8am and 9am
	time := time.Now().In(time.UTC).Hour()
	offset := 8 - time
	if offset < (-11) {
		offset = 12 + (offset % 12)
	}
	return offset
}

func main() {
	lm = loadLocations()
	for k, v := range lm {
		fmt.Println(k, v)
	}

	fmt.Println(lm.offset())
	fmt.Println("Good Morning from", lm.Name())
	http.HandleFunc("/", servefun)

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func servefun(w http.ResponseWriter, r *http.Request) {
	htmltemplate, parseErr := tpl.New("html").Parse(template)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	htmltemplate.Execute(w, &lm)

}

func loadLocations() locmap {
	path, dirErr := fp.Abs(fp.Dir(os.Args[0]))
	if dirErr != nil {
		log.Fatal(dirErr)
	}

	path += string(fp.Separator) + "sltz.txt"
	file, fileErr := os.Open(path)
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	defer file.Close()

	fileContents, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		log.Fatal(readErr)
	}

	lines := strings.Split(string(fileContents), "\n")
	ret := locmap(make(map[int]panorama))

	for linenr, line := range lines {
		//ignore empty lines
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 3 {
			log.Fatal("Error in line " + strconv.Itoa(linenr) +
				", not three elements.")
		}
		offset, parseErr := strconv.Atoi(parts[0])
		if parseErr != nil {
			log.Fatal(parseErr)
		}
		ret[offset] = panorama{parts[1], parts[2]}
	}
	return ret
}
