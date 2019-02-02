package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	// "io/ioutil"
	// "path/filepath"
	// rel "reflect"
)

func main() {
	fmt.Println("runing")
	writeFile()
	readFile()
}
func writeFile() {
	filePath := "a.txt"
	fout, err := os.Create(filePath)
	defer fout.Close()
	if err != nil {
		fmt.Println(filePath, err)
		return
	}
	for i := 0; i < 10; i++ {
		fout.WriteString("Jst a test ! \r\n")
		fout.Write([]byte("Jst a test!\r\n"))
	}
}

func readFile() {
	userFile := "a.txt"
	fin, err := os.Open(userFile)
	defer fin.Close()
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	buf := make([]byte, 1024)
	for {
		n, _ := fin.Read(buf)
		if 0 == n {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}

func bufioTest() {
	fi, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	w := bufio.NewWriter(fo)

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		if n2, err := w.Write(buf[:n]); err != nil {
			panic(err)
		} else if n2 != n {
			panic("error in writing")
		}
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}
}

//
// func ioutilTest(){
// 	b, err := ioutil.ReadFile("input.txt")
//   if err != nil { panic(err) }
//
//   err = ioutil.WriteFile("output.txt", b, 0644)
//   if err != nil { panic(err) }
// }
//
// func getFilelist(path string) {
//         err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
//                 if ( f == nil ) {return err}
//                 if f.IsDir() {return nil}
//                 println(path)
//                 return nil
//         })
//         if err != nil {
//                 fmt.Printf("filepath.Walk() returned %v\n", err)
//         }
// }
