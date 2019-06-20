package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"log"
	"strings"
	"time"
	"io/ioutil"
	"archive/zip"
)


func main() {
	defaultDir, _ := os.UserHomeDir()
	source_dir := flag.String("source_dir", defaultDir, "Full Path of directory to archive")
	dest_dir := flag.String("dest_dir", "", "Full Path of archive directory")

	flag.Parse()

	fmt.Println("Archiving...")
	source, dest := ValidateDirs(*source_dir, *dest_dir)
	archiveFile := ZipWriter(source)
	source = source + "/" + archiveFile
	dest = dest + archiveFile
	MoveFile(source, dest)
}

func MoveFile(source, dest string) {
    inputFile, err := os.Open(source)
    if err != nil {
        log.Fatal(err)
    }
    outputFile, err := os.Create(dest)
    if err != nil {
        inputFile.Close()
        log.Fatal(err)
    }
    defer outputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        log.Fatal(err)
    }
    err = os.Remove(source)
    if err != nil {
        log.Fatal(err)
    }
}

func ValidateDirs(dir1, dir2 string) (string, string) {
	_, err1 := os.Stat(dir1)
	if err1 != nil {
		log.Fatal(err1)
	}

	_, err2 := os.Stat(dir2)
	if err2 != nil {
		log.Fatal("Please add archive directory!")
	}

	if !strings.HasSuffix(dir2, "/") {
		dir2 = dir2 + "/"
	}
	return dir1, dir2
}

func GenerateFileName(defaultDir string) string {
	sliced := strings.Split(defaultDir, "/")
	dirName := sliced[len(sliced)-1]
	formatedTimestamp := time.Now().Format(time.UnixDate)
	timestamp := strings.ReplaceAll(strings.ReplaceAll(formatedTimestamp, " ", "_"), ":", "-")
	filename := dirName + "_" + timestamp + ".zip"
	return filename
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
    files, err := ioutil.ReadDir(basePath)
    if err != nil {
        fmt.Println(err)
    }

    for _, file := range files {
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

func ZipWriter(dir string) string {
    fileName := GenerateFileName(dir)
    outFile, err := os.Create(dir + "/" + fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer outFile.Close()

    if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

    writer := zip.NewWriter(outFile)
    addFiles(writer, dir, "")

    err = writer.Close()
    if err != nil {
        log.Fatal(err)
    }
    return fileName
}
