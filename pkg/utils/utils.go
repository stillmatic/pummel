package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadData(inpFile string) ([]map[string]interface{}, error) {
	//	read csv and parse into list of maps - each row is a map
	f, err := os.Open(inpFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	headers := make([]string, 0)
	records := make([]map[string]interface{}, 0)
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if len(headers) == 0 {
			for _, h := range strings.Split(line, ",") {
				headers = append(headers, h)
			}
		} else {
			record := make(map[string]interface{})
			for i, v := range strings.Split(line, ",") {
				record[headers[i]], err = strconv.ParseFloat(v, 64)
				if err != nil {
					fmt.Println("Error", err, "parsing float", v)
					record[headers[i]] = v
				}
			}
			records = append(records, record)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
