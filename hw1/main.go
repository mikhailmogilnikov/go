package main

import (
	"fmt"
	"runtime"
	"os"
)

func main() {
	printUserName()
	printCLIArgs()
	printGoVersion()
}

func printUserName() {
	username := os.Getenv("USER")
	if username == "" {
		username = "Guest"
	}

	fmt.Println("Hello,", username, "!")
}

func printCLIArgs() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No arguments provided")
	} else {
		fmt.Println("Arguments:", args)
	}
}

func printGoVersion() {
	fmt.Println("Go version:", runtime.Version())
}
// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"flag"
// )

// func main() {
// 	var x int
// 	var y bool
// 	var z string

// 	fmt.Println(x, y, z)

// 	if x > 0 {
// 		fmt.Println("x is positive")
// 	} else {
// 		fmt.Println("x is negative")
// 	}

// 	for i := range 5 {
// 		fmt.Println(i)
// 	}

// 	switch z {
// 	case "hello":
// 		fmt.Println("hello")
// 	case "world":
// 		fmt.Println("world")
// 	default:
// 		fmt.Println("default")
// 	}

// 	test()
// 	fmt.Println(sum(10, 20))

// 	p := Person{Name: "Misha", Age: 22}
// 	p.SayHello()
// 	p.SayGoodbye()

// 	v, r := divide2(10, 3)
// 	fmt.Println(v, r)

// 	res, err := divide(10, 1)
// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		fmt.Println(res)
// 	}

// 	name := os.Getenv("USER")

// 	if name == "" {
// 		name = "Guest"
// 	}

// 	lang := flag.String("lang", "en", "language")
// 	flag.Parse()

// 	args := flag.Args()

// 	fmt.Printf("Hello, %s! Lang=%s, Args=%s\n", name, *lang, args)
// }

// func test() int {
// 	return 10
// }

// func sum(x int, y int) int {
// 	return x + y
// }

// func divide(a, b float64) (float64, error) {
// 	if b == 0 {
// 		return 0, fmt.Errorf("division by zero")
// 	}
// 	return a / b, nil
// }

// type Person struct {
// 	Name string
// 	Age  int
// }

// func (p Person) SayHello() {
// 	fmt.Println("Hello, my name is", p.Name)
// }

// func (p Person) SayGoodbye() {
// 	fmt.Println("Goodbye, my name is", p.Name)
// }

// func divide2(a, b int) (quotient, remainder int) {
// 	quotient = a / b
// 	remainder = a % b
// 	return
// }
