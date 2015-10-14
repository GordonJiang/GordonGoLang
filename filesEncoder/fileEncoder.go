package main

import (
	"path/filepath"
	"log"
	"os"
	"flag"
	"crypto/sha256"
	"io/ioutil"
	"sync/atomic"
	"sync"
	"runtime"
)

var rootPath = flag.String("path", "/home/Gordon/gocode", "root path to enumerate files")
var threads = flag.Int("threads", 100, "number of thread pool to encode the file")
var outputRoot = flag.String("out", "/home/Gordon/out", "root output path for encoded files")
var filePathCh chan string
var filesEncoded uint32 = 0

func init(){
	// initialize the logger
	log.SetOutput(os.Stdout)
}

func main(){
	flag.Parse()

	filePathCh = make(chan string, (*threads) * 2)
	//finish := make(chan bool, *threads)
	var waitgroup sync.WaitGroup

	// create thread pool
	initThreadPool(filePathCh, &waitgroup)

	// kick off the thread enumerating files under directory rootPath
	enumerateFiles(*rootPath, filePathCh)

	// wait until all files are completed
	log.Println("Wait until all thread exit")
	waitgroup.Wait()
	log.Printf("Totally processed %d files", atomic.LoadUint32(&filesEncoded))
}

func initThreadPool(ch <-chan string, waitgroup *sync.WaitGroup){
	log.Printf("Creating thread pool with %d threads", *threads)
	for i:=0; i<*threads; i++{
		waitgroup.Add(1)
		go encoder(i,ch,waitgroup)
	}
}

func encoder(id int, ch <-chan string, wg *sync.WaitGroup){
	defer wg.Done()
	counter := 0
	for fullPath := range ch{
		data, err := ioutil.ReadFile(fullPath);
		if  err!=nil{
			log.Fatalf("Failed reading file:%q. Error:%q", fullPath, err)
			continue
		}
		
		hashedData := sha256.Sum256([]byte(data))
		_, file := filepath.Split(fullPath)
		outPath := filepath.Join(*outputRoot, file)
		if err = ioutil.WriteFile(outPath, hashedData[0:],0644); err!=nil{
			log.Fatal(err)
		}
		atomic.AddUint32(&filesEncoded,1)
		runtime.Gosched()

		counter++
		//log.Printf("encoder [%d] completed encoding file %q. output is %q.", id, fullPath, outPath)
	}
	log.Printf("exiting encorder [%d]. This encoder totally encoded %d files", id, counter)
}

func enumerateFiles(root string, ch chan string) (err error){
	defer close(filePathCh)
	
	// first version - using slow filePath.Walk to enumerate all files/folders
	err = filepath.Walk(root, func (path string, info os.FileInfo, er error) error{
		if er !=nil{
			log.Println("pathWalker error: ",er)
			return nil
		}
		if !info.Mode().IsRegular(){
			//log.Printf("found directory %q. skip encoding", path)
			return nil
		}
		if path == "/usr/local/man"{
			log.Printf("found /usr/local/man. Its size:%d. Its Mode:%d",info.Size(), uint32(info.Mode()))
		}
		filePathCh<-path
		//log.Printf("found file %q. Send to encoder", path)
		return nil
	})
	if err!=nil {
		log.Fatal("Error in enumerating files.", err)
	}
	return
}
