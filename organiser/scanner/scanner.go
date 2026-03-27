package scanner

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func GetDirectory(pathFlag string, input string) (string, error) {
	if pathFlag != "" {
		return pathFlag, nil
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("No directory provided")
	}

	return input, nil
}


func Scan() (string, error) {
	pathFlag := flag.String("path", "", "Directory to organize")
	flag.Parse()

	var input string

	if *pathFlag == "" {
		fmt.Print("Enter directory path: ")
		reader := bufio.NewReader(os.Stdin) 
		in, _ := reader.ReadString('\n') 	
		input = in
	}

	dir, err := GetDirectory(*pathFlag, input)
	if err != nil {
		return "", err  
	}

	fmt.Println("Working on directory:", dir)
	return dir, nil 
}