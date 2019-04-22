package main

import (
	"fmt"
	"os"

	"github.com/lunixbochs/struc"
)

type points struct {
	Len   int16   `struc:"big,int16,sizeof=Point"`
	Point []point `struc:"[]points,big"`
}

type point struct {
	Dimension int8    `struc:"big,int8,sizeof=Values"`
	Values    []int16 `struc:"[]int16,big"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "USAGE: go run main.go <file> <r|w>")
		return
	}
	fmt.Println("A awesome program for reading/writing n-dimensional points!!!")
	filename := os.Args[1]
	mode := os.Args[2]

	switch mode {
	case "r", "R":
		readMode(filename)
	case "w", "W":
		writeMode(filename)
	default:
		fmt.Fprintln(os.Stderr, "Unrecognized file mode")
		return
	}
}

func readMode(filename string) {
	fmt.Printf("READ mode for the file %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	data := points{}
	struc.Unpack(file, &data)
	for i := int16(0); i < data.Len; i++ {
		fmt.Printf("#%d (%dD):", i, data.Point[i].Dimension)
		for j := int8(0); j < data.Point[i].Dimension; j++ {
			fmt.Printf(" %d", data.Point[i].Values[j])
		}
		fmt.Println()
	}
}
func writeMode(filename string) {
	var choose int = 1
	data := points{
		Len:   0,
		Point: make([]point, 0),
	}
	fmt.Printf("WRITE mode for the file %s\n", filename)

	for choose != 0 {
		fmt.Print("MENU:\t0) write!\n\t1) add point\n\t2) view points\n>")
		if n, _ := fmt.Scanf("%d", &choose); n <= 0 {
			choose = -1
		}
		switch choose {
		case 1:
			var dim int8
			fmt.Print("Point dimension: ")
			fmt.Scanf("%d", &dim)
			data.Len++
			data.Point = append(data.Point, point{
				Dimension: dim,
				Values:    make([]int16, dim),
			})
			fmt.Print("All coordinates: ")
			for i := int8(0); i < dim; i++ {
				fmt.Scanf("%d", &data.Point[data.Len-1].Values[i])
			}
		case 2:
			for i := int16(0); i < data.Len; i++ {
				fmt.Printf("#%d (%dD):", i, data.Point[i].Dimension)
				for j := int8(0); j < data.Point[i].Dimension; j++ {
					fmt.Printf(" %d", data.Point[i].Values[j])
				}
				fmt.Println()
			}
		}
	}

	//WRITE!
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = struc.Pack(file, &data)
	if err != nil {
		panic(err)
	}
}
