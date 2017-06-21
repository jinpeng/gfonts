package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	PORT        = ":7011"
	CSS_HOST    = "fonts.sample.com"
	FONT_HOST   = "fonts.sample.com"
	FONT_FOLDER = "./static"
)

var cssContents map[string]string = make(map[string]string)
var fontContents map[string][]byte = make(map[string][]byte)

func downloadFromUrl(urlString string) {
	u, err := url.Parse(urlString)
	if err != nil {
		log.Fatal(err)
		return
	}
	tokens := strings.Split(u.Path, "/")
	fileName := tokens[len(tokens)-1]
	folder := FONT_FOLDER + strings.Join(tokens[:len(tokens)-1], "/")
	filePath := folder + "/" + fileName
	// Check if file exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, os.FileMode(0777))
		}
	} else {
		return
	}
	fmt.Println("Downloading", urlString, "to", filePath)

	output, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error while creating", filePath, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(urlString)
	if err != nil {
		fmt.Println("Error while downloading", urlString, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", urlString, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

func fetch(url string) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return content, nil
}

func fetchFont(url string) {
	downloadFromUrl(url)
}

func fetchCSS(url string) (string, error) {
	if css, ok := cssContents[url]; ok {
		return css, nil
	}
	cssbytes, err := fetch(url)
	if err != nil {
		return "", err
	} else {
		css := string(cssbytes)
		re := regexp.MustCompile("url([^)]*)")
		for _, fontUrl := range re.FindAllString(css, -1) {
			fontUrl = fontUrl[4:]
			fetchFont(fontUrl)
		}
		reHost := regexp.MustCompile("fonts\\.gstatic\\.com")
		css = reHost.ReplaceAllString(css, FONT_HOST)
		cssContents[url] = css
	}
	return string(cssbytes), err
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://fonts.googleapis.com/css?" + r.URL.Query().Encode()
	css, ok := cssContents[url]
	if !ok {
		_, err := fetchCSS(url)
		if err != nil {
			fmt.Fprintf(w, "Failed to retrieve CSS: %s", url)
			return
		}
	}
	css, ok = cssContents[url]
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	fmt.Fprint(w, css)
}

func addDefaultHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}

func main() {
	http.HandleFunc("/css", cssHandler)
	http.Handle("/", addDefaultHeaders(http.FileServer(http.Dir(FONT_FOLDER))))
	http.ListenAndServe(PORT, nil)
}
