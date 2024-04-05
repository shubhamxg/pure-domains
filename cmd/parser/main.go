package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/shubhamxg/pure-domains/internal/site"
	"github.com/shubhamxg/pure-domains/pkg"
	"github.com/sqweek/dialog"
)

var (
	total_lines int
	save        *pkg.SaveFile
)

func main() {
	file_dir, err := dialog.File().Load()
	if err != nil {
		fmt.Println("Failed to load file Error > ", err.Error())
	}

	open_file, err := os.Open(file_dir)
	if err != nil {
		fmt.Println("Failed to open file Error > ", err.Error())
	}
	defer open_file.Close()

	file_scanner := bufio.NewScanner(open_file)
	for file_scanner.Scan() {
		total_lines++
	}

	if file_scanner.Err() != nil {
		fmt.Println("Scanner Error > ", file_scanner.Err().Error())
	} else if total_lines == 0 {
		fmt.Println("Loaded File is Empty")
	}

	_, err = open_file.Seek(0, 0)
	if err != nil {
		fmt.Println("Cannot point pointer to the start of file")
	}

	file_reader := bufio.NewReader(open_file)
	save = pkg.Save()
	for {
		line, _, err := file_reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("Failed to read the line: %s", line)
		}

		if len(line) > 0 {
			// CODE
			dmarc_result := site.Dmarc(string(line))
			save.File(dmarc_result, string(line))

			fmt.Printf("[~] Site > %s | Subdomains > %d\n", line, len(dmarc_result))
		}

	}

}
