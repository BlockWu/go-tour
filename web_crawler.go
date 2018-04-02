package main

import (
	"fmt"
	"sync"
)

type Safe struct {
	aurls []string
	mux   sync.Mutex
}

var num int = 0

type Fetcher interface {
	// Fetch ���� URL �� body ���ݣ����ҽ������ҳ�����ҵ��� URL �ŵ�һ�� slice �С�
	Fetch(url string) (body string, urls []string, err error)
}

func gothrough(url string, aurls []string) bool {
	for _, v := range aurls {
		//fmt.Println(v + "hshsaihihs")
		if url == v {
			return false
		}
	}
	return true
}

var crawl = Safe{make([]string, 10), sync.Mutex{}}

//ȫ�ֱ���
var c1 chan string

var quit chan bool

// Crawl ʹ�� fetcher ��ĳ�� URL ��ʼ�ݹ����ȡҳ�棬ֱ���ﵽ�����ȡ�
func Crawl(url string, depth int, fetcher Fetcher) {
	//fmt.Println(url)
	// TODO: ���е�ץȡ URL��
	// TODO: ���ظ�ץȡҳ�档
	// ���沢û��ʵ���������������
	if depth <= 0 {
		quit <- true
		return
	}
	crawl.mux.Lock()
	if gothrough(url, crawl.aurls) == false {
		crawl.mux.Unlock()
		quit <- true
		return
	}
	crawl.aurls = append(crawl.aurls, url)
	crawl.mux.Unlock()

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		quit <- true
		return
	}

	c1 <- url

	for _, u := range urls {
		go Crawl(u, depth-1, fetcher)
	}
	for i := 0; i < len(urls); i++ {
		num++
		<-quit
	}
	quit <- true
	return

}

func main() {
	c1 = make(chan string)
	quit = make(chan bool)

	go Crawl("https://golang.org/", 4, fetcher)
	for {
		select {
		case url := <-c1:
			fmt.Println("found: ", url)
		case <-quit:
			fmt.Println(num)
			return
		}
	}
}

// fakeFetcher �Ƿ������ɽ���� Fetcher��
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher ������� fakeFetcher��
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
