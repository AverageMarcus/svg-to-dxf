package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		svgURL := r.URL.Query().Get("url")
		if svgURL == "" {
			return
		}
		fmt.Println("Fetching " + svgURL)

		filename := "/app/" + randomString(7) + ".svg"

		fmt.Println("Filename: " + filename)
		if err := downloadFile(svgURL, filename); err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Println("Downloaded SVG")

		dxf, err := convertToDXF(filename)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Println("Converted to DXF")

		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Length", fmt.Sprint(len(dxf)))
		w.Header().Add("Content-Disposition", "inline; filename=out.dxf")
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
