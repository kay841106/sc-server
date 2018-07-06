package main

import (
	"fmt"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	var A int
	d1 := []byte("123\ngo\n")
	aa, err := ioutil.ReadFile("kk")

	fmt.Println(A)
	err = ioutil.WriteFile("kk", d1, 0644)
	check(err)
	fmt.Println(string(aa))
	// 	file, err := os.Open("file")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		file, _ := os.Create("file")
	// 	}

	// 	defer file.Close()

	// 	n3, err := file.WriteString("writes\n")
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println(file)
}
