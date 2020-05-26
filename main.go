// test project main.go
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	k := 5
	stdin := bufio.NewReader(os.Stdin)
	text, _ := stdin.ReadString('\n')
	reg := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	urls := reg.FindAllString(text, -1)
	urlsRangeArray := rangeArray(urls, k)
	total := 0

	for i := 0; i < len(urlsRangeArray); i++ {
		fmt.Printf("Running channel with %d urls\n", len(urlsRangeArray[i]))
		ch := make(chan int, len(urlsRangeArray[i]))
		for _, url := range urlsRangeArray[i] {
			go RequestUrls(url, ch)
			count := <-ch
			total += count
			fmt.Printf("Count for %s: %d\n", url, count)
		}
		close(ch)
	}
	fmt.Printf("Total: %d\n", total)
}

func rangeArray(a []string, k int) [][]string {
	r := (len(a) + k - 1) / k
	b := make([][]string, r)
	lo, hi := 0, k
	for i := range b {
		if hi > len(a) {
			hi = len(a)
		}
		b[i] = a[lo:hi:hi]
		lo, hi = hi, hi+k
	}
	return b
}

func RequestUrls(url string, ch chan<- int) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	go ParseBody(url, string(body), ch)
}

func ParseBody(url string, body string, ch chan<- int) {
	searchString := "Go"
	result := strings.Count(body, searchString)
	ch <- result
}
