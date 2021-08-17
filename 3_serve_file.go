package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var FileDoesNotExist = errors.New("NotFound")

func main() {

	//Command To Execute -  go build 3_serve_file.go && ./3_serve_file

	http.HandleFunc("/", handleIndex)

	//NOTE :: Graceful shutdown not there
	http.ListenAndServe(":5000", nil)
}

func handleIndex(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("Content-Type", "text/html")

	//If not root path then go to handle file
	if r.URL.Path != "/" {
		handleFile(rw, r)
		return
	}

	rw.WriteHeader(http.StatusOK)
	q := r.URL.Query()
	message := q.Get("message")

	//If query parameter has data then write that data
	if message != "" {
		rw.Write([]byte(message))
		return
	}

	rw.Write([]byte("Hello World !!!"))
}

func handleFile(rw http.ResponseWriter, r *http.Request) {

	fd, err := getFileDetails(r.URL.Path[1:])

	if err == FileDoesNotExist {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(FileDoesNotExist.Error()))
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Some error occured"))
		return
	}

	rw.Write([]byte(fd))

}

func getFileDetails(f string) ([]byte, error) {

	//NOTE :: We can cache the file in memory / use templates instead of getting it every time for faster load

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Println(err)
	}

	fp := dir + "/views/" + f //Absolute File Path

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return []byte{}, FileDoesNotExist
	}

	fileDetails, err := os.ReadFile(fp)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return fileDetails, nil

}
