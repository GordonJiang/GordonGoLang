package main

import (
	"fmt"
	"github.com/GordonJiang/GordonGoLang/stringutil"
	"flag"
	"html/template"
	"os"
	"runtime/pprof"
	"math"
	"time"
	"regexp"
)

var (
address = flag.String("addr",":8080","http service address")
times = flag.Int("times", 1000, "loop times")
cpuprofile = flag.String("cpuprofile", "", "file to save performance profile output")
templateStr = "<head/> {{if .}} section in if. Name: {{.Name}}, Value: {{.Value}}  {{else}} section in else {{end}}"
ch = make(chan int)
)

type ValuePair struct{
	Name string
	Value int
}

func init(){
	fmt.Println("In init Function")
}

func main(){

	// flag testing code
	flag.Parse()
	fmt.Printf("address variable value: %q Type is: %T", (*address), address)
	fmt.Println("")

	// performance profiling
	if *cpuprofile != ""{
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

	// template testing code
	valPair := ValuePair{"Windows Intune", 1000000}
	templ, err := template.New("test template").Parse(templateStr); 
	if err != nil{
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
	for i,c := range str {
		fmt.Printf("%q, index:%d",c,i)
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
}

func Produce(){
	for n:=0;;n++{
		ch<-n
	}
}

func DoWork(n int){
	total := float64(n + 100)
	for i:=0;i<n;i++{
		num := total + float64(<-ch) + float64(i)
		num = math.Sqrt(num * total)
	}
}

func tickers(){
	ticker := time.NewTicker(time.Millisecond * 500)
	go func (){
		for t := range ticker.C{
			fmt.Println("Tick at ", t)
		}
	}()
	time.Sleep(time.Second * 5)
	ticker.Stop()
	fmt.Println("Exit tickers function")
}

func regex(){
	r := regexp.MustCompile("a(x*)b(y|z)c")
	fmt.Printf("%q\n", r.FindStringSubmatch("-axxxbyc-"))
	fmt.Printf("%q\n", r.FindStringSubmatch("-abzc-"))
}
