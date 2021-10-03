package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var URL = "https://dl2.soft98.ir/adobe/Blender.2.80.x86.rar?1633239795"
var numberOfRoutings int = 10

type finalResult struct {
	data []byte
	mu   sync.Mutex
}

func (f *finalResult) addBytes(newData []byte) {
	f.mu.Lock()
	f.data = append(f.data, newData...)
	f.mu.Unlock()
}

func main() {
	data := &finalResult{make([]byte, 0), sync.Mutex{}}
	res, _ := http.Head(URL)
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0])
	len_sub := length / numberOfRoutings
	diff := length % numberOfRoutings
	//body := make([]string, numberOfRoutings+1)
	for i := 0; i < numberOfRoutings; i++ {
		wg.Add(1)

		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range

		if i == numberOfRoutings-1 {
			max += diff // Add the remaining bytes in the last request
		}

		go func(min int, max int, i int) {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", URL, nil)
			range_header := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
			req.Header.Add("Range", range_header)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			reader, _ := ioutil.ReadAll(resp.Body)
			data.addBytes([]byte(string(reader)))
			//ioutil.WriteFile(strconv.Itoa(i), []byte(string(body[i])), 0x777)
			wg.Done()
			//          ioutil.WriteFile("new_oct.png", []byte(string(body)), 0x777)
		}(min, max, i)
	}
	wg.Wait()
	ioutil.WriteFile("downloaded", data.data, 0x777)
}
