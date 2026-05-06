package scanner

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

func ScanWithDeps(pathFlag string, reader io.Reader, writer io.Writer) (string, error) {
	var input string
 
	if pathFlag == "" {
		fmt.Fprint(writer, "Enter directory path: ")
		bufReader := bufio.NewReader(reader)
		in, err := bufReader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		input = in
	}
 
	dir, err := GetDirectory(pathFlag, input)
	if err != nil {
		return "", err
	}
 
	fmt.Fprintf(writer, "Working on directory: %s\n", dir)
	return dir, nil
}
 
// Scan is the original function that calls ScanWithDeps with real dependencies
func Scan() (string, error) {
	pathFlag := flag.String("path", "", "Directory to organize")
	flag.Parse()
 
	return ScanWithDeps(*pathFlag, os.Stdin, os.Stdout)
}