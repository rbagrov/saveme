package main

import (
	"fmt"
	"flag"
	"os"
	"log"
	"strings"
	"time"
	"io/ioutil"
	"archive/zip"
)


func main() {
	homeDir, _ := os.UserHomeDir()
	dir := flag.String("dir", homeDir, "Full Path of directory to archive")
	flag.Parse()
	_, err := os.Stat(*dir)
	if err != nil {
		log.Fatal(err)
	}
	ZipWriter(*dir)
}

func GenerateFileName(homeDir string) string {
	sliced := strings.Split(homeDir, "/")
	userName := sliced[len(sliced)-1]
	formatedTimestamp := time.Now().Format(time.UnixDate)
	timestamp := strings.ReplaceAll(strings.ReplaceAll(formatedTimestamp, " ", "_"), ":", "-")
	filename := userName + "_" + timestamp + ".zip"
	return filename
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
    files, err := ioutil.ReadDir(basePath)
    if err != nil {
        fmt.Println(err)
    }

    for _, file := range files {
        fmt.Println(basePath + file.Name())
        if !file.IsDir() && file.Mode().IsRegular() {
            dat, _ := ioutil.ReadFile(basePath + file.Name())

            f, _ := w.Create(baseInZip + file.Name())
            _, err = f.Write(dat)
            if err != nil {
                log.Fatal(err)
            }
        } else if file.IsDir() {
            newBase := basePath + file.Name() + "/"
            addFiles(w, newBase, file.Name() + "/")
        }
    }
}

func ZipWriter(homeDir string) {
    outFile, err := os.Create(GenerateFileName(homeDir))
    if err != nil {
        log.Fatal(err)
    }
    defer outFile.Close()

    writer := zip.NewWriter(outFile)
    homeDir = homeDir + "/"
    addFiles(writer, homeDir, "")

    if err != nil {
        log.Fatal(err)
    }

    err = writer.Close()
    if err != nil {
        log.Fatal(err)
    }
}
