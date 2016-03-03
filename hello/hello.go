package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/GordonJiang/GordonGoLang/stringutil"
	//	"sync"
)

var (
	address     = flag.String("addr", ":8080", "http service address")
	times       = flag.Int("times", 1000, "loop times")
	cpuprofile  = flag.String("cpuprofile", "", "file to save performance profile output")
	templateStr = "<head/> {{if .}} section in if. Name: {{.Name}}, Value: {{.Value}}  {{else}} section in else {{end}}"
	ch          = make(chan int)
)

type ValuePair struct {
	Name  string
	Value int
}

func init() {
	fmt.Println("In init Function")
}

type s1 struct {
	name  string
	tp    string
	files []string
}

type s2 struct {
	name  string
	tp    string
	files []string
}

type Privilege int32

const (
	Role_Reader Privilege = 0
	Role_Admin  Privilege = 1
	Role_Dev    Privilege = 2
)

type Person struct {
	name  string
	email string
}

type Admin struct {
	admPrivileges []Privilege
	*Person
}

// Test if two sorted string slice are equal.
func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func changeParam(s string) {
	fmt.Println("In changeParam: ", s)
	s = "New value"
}

func main() {
	fis, err := ioutil.ReadDir("/home/gjiang")
	if err != nil {
		os.Exit(2)
	}
	for _, fi := range fis {
		fmt.Printf("%s. %v\n", fi.Name(), fi.IsDir())
	}

	return

	notify := make(chan int)
	go func() {
		for {
			time.Sleep(time.Second * 20)
			notify <- 1
		}
	}()

	notify2 := make(chan int)
	go func() {
		for {
			time.Sleep(time.Second * 20)
			notify2 <- 1
		}
	}()

	timeout := time.After(time.Second * 15)
	for {
		select {
		case r := <-notify:
			fmt.Println("Notified by job1", r)
			return
		case r := <-notify2:
			fmt.Println("Notified by job2", r)
			return
		case <-timeout:
			fmt.Println("Timed out!", time.Now())
			os.Exit(1)
		}
	}

	return

	ch := make(chan int)

	go func(out chan<- int) {
		for i := 0; i < 100; i++ {
			out <- rand.Int()
		}
		close(out)
	}(ch)

	// goroutine receives data and keep printing out
	go func(in <-chan int) {
		for x := range in {
			fmt.Println("Received:", x)
		}
	}(ch)

	time.Sleep(10 * time.Second)

	first := func() *http.Response {
		ch := make(chan *http.Response, 3)
		go func() { resp, _ := http.Get("http://www.microsoft.com"); ch <- resp }()
		go func() { resp, _ := http.Get("http://www.google.com"); ch <- resp }()
		go func() { resp, _ := http.Get("http://www.facebook.com"); ch <- resp }()
		return <-ch
	}()

	fmt.Println("The first caught web is: ", first.Request)
	return

	// go routine receives data and keep printing out
	go func() {
		for {
			fmt.Println("Received:", <-ch)
		}
	}()
	time.Sleep(10 * time.Second)

	src, dst := "/home/gjiang/tmp", "/home/gjiang/output"

	var buf bytes.Buffer
	cont, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println("Failed to load file:", src)
		os.Exit(-1)
	}
	buf.WriteString("Hash: ")
	h := sha256.Sum256(cont)

	for _, b := range h {
		buf.WriteString(fmt.Sprintf("%X", b))
	}

	ioutil.WriteFile(dst, buf.Bytes(), 0777)

	return

	dir := "abc/qqq/ppp/xyz"
	dirs := filepath.SplitList(dir)
	fmt.Println("Dirs: ", dirs)

	ss := strings.Split("", `/`)
	fmt.Println("size: ", len(ss))
	fmt.Println(ss[0])

	path := filepath.Join("", ss[0])
	fmt.Println(path)
	path = filepath.Join(path, "owners")
	fmt.Println(path)

	// try deepequal
	sa1 := []string{"abc", "xyz", "cc"}
	sa2 := []string{"xyz", "cc"}
	fmt.Println("string array1 and array 2: ", StringSliceEqual(sa1, sa2))

	// try string splitN
	res := make([]string, 2)
	fmt.Println("len:", len(res))
	res = strings.SplitN("repo/fsdfs/few/rewrw.ext", "/", 2)
	fmt.Println("len:", len(res))
	fmt.Println("Split result: ", res[0], res[1])

	// runtime/debug callstack
	func1()

	// test string, rune memory
	s := "abcde"
	r := []rune(s)
	fmt.Printf("r's address: %d , s's address: %x \n", &r, &s)
	r[0] = '事'
	fmt.Println(string(r))
	fmt.Println(s)

	// test empty slice
	var arr []int
	fmt.Println("array length:", arr)

	// flag testing code
	flag.Parse()
	fmt.Printf("address variable value: %q Type is: %T", (*address), address)
	fmt.Println("")

	// performance profiling
	if *cpuprofile != "" {
		fmt.Println("Creating CPU profile")
		pprofile, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Println("Failed to create CPU profile file:", *cpuprofile)
			return
		}
		pprof.StartCPUProfile(pprofile)
		defer pprof.StopCPUProfile()
		go Produce()
		DoWork(*times)
	}

	// time format
	t := time.Now().UTC()
	fmt.Println("current date default output:", t)
	fmt.Println("Current date is:", t.Format("2006-01-02 15:04:05"))

	// template testing code
	valPair := ValuePair{"Windows Intune", 1000000}
	templ, err := template.New("test template").Parse(templateStr)
	if err != nil {
		fmt.Println("Failed to create template with template str: %q", templateStr)
		return
	}
	templ.Execute(os.Stdout, valPair)
	fmt.Println("")

	// string rune testing code
	str := "中国china"
	fmt.Println(stringutil.Reverse(str))
	//runes := make([]rune,2)
	//index:=0
	for i, c := range str {
		fmt.Printf("%q, index:%d", c, i)
		//runes[1-index] = c
		//index++
	}
	/*
		fmt.Printf("%q", runes)
		fmt.Println("with bytes")

		bytes := []byte(str)
		fmt.Printf("len:%d, value:%q",len(bytes), bytes)
		for i,j:=0,len(bytes)-1; i<j; i,j = i+1, j-1{
			bytes[i],bytes[j] = bytes[j], bytes[i]
		}
		fmt.Printf("%q",bytes)


		fmt.Println("using rune")
		r := []rune(str)
		fmt.Println(r)
	*/

	// regex
	regex()

	// ticker, timer functions
	//tickers()

	// environment variables
	fmt.Println("GoPath variable value:", os.Getenv("GOPATH"))
	fmt.Println("All environment variables", os.Environ())

	// json
	data := map[string]interface{}{
		"one": 1,
		"two": 2,
	}
	encoded, _ := json.Marshal(data)
	var decoded map[string]interface{}
	json.Unmarshal(encoded, &decoded)
	fmt.Printf("json encoded: %q", encoded)
	fmt.Println("decoded value:", decoded)
	fmt.Println("data:", data)
	fmt.Println("data == decoded:", reflect.DeepEqual(data, decoded))

	// sync pool
	//val := 0
	//p := sync.Pool
	/*
		{	New:func()interface{}{
			return val
		}
		}
	*/
}

func func1() {
	fmt.Println("In func1:")
	fmt.Println(string(debug.Stack()))
}

func Produce() {
	for n := 0; ; n++ {
		ch <- n
	}
}

func DoWork(n int) {
	total := float64(n + 100)
	for i := 0; i < n; i++ {
		num := total + float64(<-ch) + float64(i)
		num = math.Sqrt(num * total)
	}
}

func tickers() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at ", t)
		}
	}()
	time.Sleep(time.Second * 5)
	ticker.Stop()
	fmt.Println("Exit tickers function")
}

func regex() {
	r := regexp.MustCompile("a(x*)b(y|z)c")
	fmt.Printf("%q\n", r.FindStringSubmatch("-axxxbyc-"))
	fmt.Printf("%q\n", r.FindStringSubmatch("-abzc-"))
}
