package main

import (
	"path/filepath"
	"log"
	"os"
	"flag"
	"crypto/sha256"
	"io/ioutil"
)

var rootPath = flag.String("path", "/home/Gordon/gocode", "root path to enumerate files")
var threads = flag.Int("threads", 100, "number of thread pool to encode the file")
var outputRoot = flag.String("out", "/home/Gordon/out", "root output path for encoded files")
var filePathCh chan string

func init(){
	// initialize the logger
	log.SetOutput(os.Stdout)
}

func main(){
	flag.Parse()

	filePathCh = make(chan string, (*threads) * 2)
	finish := make(chan bool, *threads)

	// create thread pool
	initThreadPool(filePathCh, finish)

	// kick off the thread enumerating files under directory rootPath
	enumerateFiles(*rootPath, filePathCh)

	// wait until all files are completed
	log.Println("Wait until all thread exit")
	for t:=0; t<*threads; t++{
		<-finish
	}
}

func initThreadPool(ch <-chan string, finish chan bool){
	log.Printf("Creating thread pool with %d threads", threads)
	for i:=0; i<*threads; i++{
		go encoder(i,ch,finish)
	}
}

func encoder(id int, ch <-chan string, finish chan bool){
	//log.Println("Started encoder with id:", id)
	defer log.Printf("exiting encorder [%d]", id)
	for fullPath := range ch{
		data, err := ioutil.ReadFile(fullPath);
		if  err!=nil{
			log.Fatal(err)
			continue
		}
		
		hashedData := sha256.Sum256([]byte(data))
		_, file := filepath.Split(fullPath)
		outPath := filepath.Join(*outputRoot, file)
		if err = ioutil.WriteFile(outPath, hashedData[0:],0644); err!=nil{
			log.Fatal(err)
		}
		log.Printf("encoder [%d] completed encoding file %q. output is %q.", id, fullPath, outPath)
	}
	finish<-true
}

func enumerateFiles(root string, ch chan string) (err error){
	defer close(filePathCh)
	
	// first version - using slow filePath.Walk to enumerate all files/folders
	
	return filepath.Walk(root, func (path string, info os.FileInfo, er error) error{
		if er !=nil{
			log.Println("pathWalker error: ",er)
			return nil
		}
		if info.IsDir(){
			//log.Printf("found directory %q. skip encoding", path)
			return nil
		}
		filePathCh<-path
		//log.Printf("found file %q. Send to encoder", path)
		return nil
	})
}
