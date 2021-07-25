package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//go:embed index.html
var content embed.FS

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		svgURL := r.URL.Query().Get("url")
		r.ParseMultipartForm(1000000)

		uploadedFile, fileHeader, err := r.FormFile("svg")

		if svgURL == "" && (fileHeader == nil || fileHeader.Size == 0) {
			body, _ := content.ReadFile("index.html")
			w.Write(body)
			return
		}

		filename := "/app/" + randomString(7) + ".svg"
		fmt.Println("Filename: " + filename)

		outFilename := "out"

		if svgURL != "" {
			fmt.Println("Fetching " + svgURL)
			re := regexp.MustCompile(`(?m)/([\w|-]+)\.svg$`)
			matches := re.FindStringSubmatch(svgURL)
			if len(matches) == 2 {
				outFilename = matches[1]
			}
			if err := downloadFile(svgURL, filename); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			fmt.Println("Downloaded SVG")
		} else {
			outFilename = strings.ReplaceAll(fileHeader.Filename, ".svg", "")
			file, err := os.Create(filename)
			if err != nil {
				return
			}
			defer file.Close()

			_, err = io.Copy(file, uploadedFile)
			if err != nil {
				return
			}
		}

		dxf, err := convertToDXF(filename)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Println("Converted to DXF")

		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", fmt.Sprint(len(dxf)))
		w.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=%s.dxf", outFilename))
		fmt.Fprintf(w, string(dxf))

		os.Remove(filename)
	})

	fmt.Println(http.ListenAndServe(":"+port, nil))
}

func downloadFile(URL, fileName string) error {
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func convertToDXF(filename string) ([]byte, error) {
	cmd := exec.Command("python", "/usr/share/inkscape/extensions/dxf_outlines.py", filename)
	return cmd.Output()
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
