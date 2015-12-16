package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/GordonJiang/GordonGoLang/stringutil"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"
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

func main() {
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
