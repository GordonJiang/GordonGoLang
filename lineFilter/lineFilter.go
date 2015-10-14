package main

import(
	"fmt"
	"bufio"
	"strings"
	"os"
)

func main(){
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan(){
		uc := strings.ToUpper(scanner.Text())
		if uc == "QUIT" || uc == "EXIT"{
			break
		}
		fmt.Printf("%q\n",uc)
	}
	
	if err := scanner.Err(); err != nil{
		fmt.Println("error:",err)
		os.Exit(1)
	}
	fmt.Println("Exit")
}
