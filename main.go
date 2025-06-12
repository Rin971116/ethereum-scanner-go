package main

import (
	"fmt"
	"my-go-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//curl http://localhost:8080/hello
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/greet/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, fmt.Sprintf("Hello, %s!", name))
	})

	router.Run(":8080") // listen and serve on

	printBrandName("123")

}

func changeName(person *models.Person, newName string) {
	person.Name = newName
}

var changeNameFunc func(*models.Person, string) = func(person *models.Person, newName string) {
	person.Name = newName
}

var a int = 10

var b string = "hello"

var c bool = true

var d float64 = 3.14

var e []int = []int{1, 2, 3}

var f map[string]int = map[string]int{"one": 1, "two": 2}

var g chan int = make(chan int)

var h struct {
	Name string
	Age  int
} = struct {
	Name string
	Age  int
}{Name: "Alice", Age: 30}

var j func(int) int = func(x int) int {
	return x * x
}

type BrandType string

var (
	BrandTypeApple     BrandType = "Apple"
	BrandTypeSamsung   BrandType = "Samsung"
	BrandTypeGoogle    BrandType = "Google"
	BrandTypeMicrosoft BrandType = "Microsoft"
	BrandTypeSony      BrandType = "Sony"
)

func printBrandName(brand BrandType) {
	switch brand {
	case BrandTypeApple:
		fmt.Println("Brand is Apple")
	case BrandTypeSamsung:
		fmt.Println("Brand is Samsung")
	case BrandTypeGoogle:
		fmt.Println("Brand is Google")
	case BrandTypeMicrosoft:
		fmt.Println("Brand is Microsoft")
	case BrandTypeSony:
		fmt.Println("Brand is Sony")
	default:
		fmt.Println("Unknown brand")
	}
}

type Person struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Birthday string `json:"birthday"`
}

type Person2 struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Birthday string `json:"birthday"`
}

var p struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Birthday string `json:"birthday"`
} = struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Birthday string `json:"birthday"`
}{Name: "Alice",
	Age:      30,
	Birthday: "1993-01-01",
}
