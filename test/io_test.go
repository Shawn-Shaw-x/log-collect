package test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"sync"
	"testing"
)

func Test_BIO(t *testing.T) {
	client := resty.New()

	// BIO GET
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"search": "golang",
			"page":   "1",
		}).
		SetHeader("Accept", "application/json").
		Get("https://httpbin.org/get")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Response Body:", resp.String())
}

// goroutine模拟NIO
func TestNIO(t *testing.T) {
	var wg sync.WaitGroup
	ch := make(chan string, 3)

	urls := []string{
		"https://httpbin.org/get?name=go1",
		"https://httpbin.org/get?name=go2",
		"https://httpbin.org/get?name=go3",
	}

	// NIO GET
	for _, url := range urls {
		wg.Add(1)
		go func(*sync.WaitGroup, chan string, string) {
			defer wg.Done()

			client := resty.New()
			resp, err := client.R().Get(url)
			if err != nil {
				ch <- fmt.Sprintf("Request to %s failed: %v", url, err)
				return
			}
			ch <- fmt.Sprintf("Response from %s: %d", url, resp.StatusCode())
		}(&wg, ch, url)
	}

	// 等待所有请求完成
	wg.Wait()
	close(ch)

	// 打印结果
	for msg := range ch {
		fmt.Println(msg)
	}
}

// multiple
func TestSimulateNIO(t *testing.T) {
	type Job struct {
		url string
		idx int
	}

	// 任务队列
	jobChan := make(chan Job, 10)
	// 结果集
	resultChan := make(chan string, 10)

	// 将5个请求绑定到2个（少量）线程上，实现多路复用
	workerCount := 2
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for job := range jobChan {
				client := resty.New()
				resp, err := client.R().Get(job.url)
				if err != nil {
					resultChan <- fmt.Sprintf("[Worker %d] Request %d failed: %v", workerID, job.idx, err)
					continue
				}
				resultChan <- fmt.Sprintf("[Worker %d] Response %d: %d", workerID, job.idx, resp.StatusCode())
			}
		}(i)
	}

	// 模拟5个请求
	urls := []string{
		"https://httpbin.org/get?name=go1",
		"https://httpbin.org/get?name=go2",
		"https://httpbin.org/get?name=go3",
		"https://httpbin.org/get?name=go4",
		"https://httpbin.org/get?name=go5",
	}

	// 模拟 Selector：统一调度任务
	go func() {
		for i, url := range urls {
			jobChan <- Job{url: url, idx: i + 1}
		}
		close(jobChan)
	}()

	// 收集结果
	for i := 0; i < len(urls); i++ {
		fmt.Println(<-resultChan)
	}
}

func asyncRequest(url string, callback func(string)) {
	go func() {
		client := resty.New()
		resp, err := client.R().Get(url)
		if err != nil {
			callback(fmt.Sprintf("Request to %s failed: %v", url, err))
			return
		}
		callback(fmt.Sprintf("Response from %s: %d", url, resp.StatusCode()))
	}()
}

func TestSimulateAIO(t *testing.T) {
	urls := []string{
		"https://httpbin.org/get?name=go1",
		"https://httpbin.org/get?name=go2",
		"https://httpbin.org/get?name=go3",
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		asyncRequest(url, func(result string) {
			fmt.Println(result)
			wg.Done()
		})
	}

	wg.Wait()
}
