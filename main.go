package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var URL = "http://localhost/rand.txt"
var numberOfRoutings int = 10


func main() {
	res, _ := http.Head(URL)
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0])
	len_sub := length / numberOfRoutings
	diff := length % numberOfRoutings
	body := make([]string, 11)
	for i := 0; i < numberOfRoutings ; i++ {
		wg.Add(1)

		min := len_sub * i // Min range
		max := len_sub * (i + 1) // Max range

		if (i == numberOfRoutings - 1) {
			max += diff // Add the remaining bytes in the last request
		}

		go func(min int, max int, i int) {
			client := &http.Client {}
			req, _ := http.NewRequest("GET", "http://localhost/rand.txt", nil)
			range_header := "bytes=" + strconv.Itoa(min) +"-" + strconv.Itoa(max-1) // Add the data for the Range header of the form "bytes=0-100"
			req.Header.Add("Range", range_header)
			resp,_ := client.Do(req)
			defer resp.Body.Close()
			reader, _ := ioutil.ReadAll(resp.Body)
			body[i] = string(reader)
			ioutil.WriteFile(strconv.Itoa(i), []byte(string(body[i])), 0x777) // Write to the file i as a byte array
			wg.Done()
			//          ioutil.WriteFile("new_oct.png", []byte(string(body)), 0x777)
		}(min, max, i)
	}
	wg.Wait()
}

