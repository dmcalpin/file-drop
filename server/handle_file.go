package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/0xAX/notificator"
)

var notify *notificator.Notificator

func init() {
	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "FileDrop",
	})
}

// HandleFile writes a file from a request
// to disk. It will dedupe files as to not overwrite
// existing ones like so: file.txt, file (1).txt
// file (2).txt. Or: myfile, myfile (1), myfile (2)
// if no file type suffix.
//
// The file is written to $HOME/Downloads.
func HandleFile(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST allowed."))
		return
	}

	log.Printf("recieved file")

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Failed to parse form."))
		return
	}

	file, fileHandler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Failed to get the file."))
		return
	}
	defer file.Close()
	fmt.Printf("Upload Info: %v\n", fileHandler.Header)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to find home dir"))
		return
	}

	filePath := homeDir + "/Downloads/" + fileHandler.Filename
	filePath = dedupePath(filePath)

	wFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to open file for writing."))
		return
	}
	defer wFile.Close()

	log.Println("Writing file:", fileHandler.Filename)
	written, err := io.Copy(wFile, file)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error writing file: " + fileHandler.Filename))
		return
	}

	notify.Push("Recieved New File", filePath, "/home/user/icon.png", notificator.UR_NORMAL)

	log.Printf("Wrote %d bytes\n", written)
}

func dedupePath(filePath string) string {

	filePathParts := strings.Split(filePath, "/")
	fileName := filePathParts[len(filePathParts)-1]
	originalFileName := strings.Split(fileName, ".")[0]

	for i := 1; i < 10000; i++ {
		fileName := filePathParts[len(filePathParts)-1]
		_, err := os.Stat(strings.Join(filePathParts, "/"))
		if os.IsNotExist(err) {
			return strings.Join(filePathParts, "/")
		}

		newFileName := ""
		if strings.Contains(fileName, ".") {
			fileNameParts := strings.Split(fileName, ".")
			newFileName = fmt.Sprintf("%s (%d).%s", originalFileName, i, fileNameParts[1])
		} else {
			newFileName = fmt.Sprintf("%s (%d)", originalFileName, i)
		}
		filePathParts[len(filePathParts)-1] = newFileName
	}

	return filePath
}
