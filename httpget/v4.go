// Ассинхронный вариант с контекстом и пулом соединений в poolHTTPReq
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type respStCWP struct {
	lenBody int64
	err     error
}

const poolHTTPReq = 2
const byteInMegabytev4 = 1024 * 1024

func main() {
	urlsList1 := []string{
		"https://youtube.com",
		"https://ya.ru",
		"https://reddit.com",
		"https://google.com",
		"https://mail.ru",
		"https://amazon.com",
		"https://instagram.com",
		"https://wikipedia.org",
		"https://linkedin.com",
		"https://netflix.com",
	}
	urlsList2 := append(urlsList1, "https://111.321", "https://999.000")

	{
		t1 := time.Now()
		byteSum, err := requestSumAsyncWithCtxAndPool(urlsList1)
		fmt.Printf("Сумма страниц в Мб=%.2f, ошибка - %v \n", (float64(byteSum) / byteInMegabytev4), err)
		fmt.Printf("Время выполнение запросов %.2f сек. \n", time.Now().Sub(t1).Seconds())
	}
	fmt.Println("++++++++")
	{
		t1 := time.Now()
		byteSum, err := requestSumAsyncWithCtxAndPool(urlsList2)
		fmt.Printf("Сумма страниц в Мб=%.2f, ошибка - %v \n", (float64(byteSum) / byteInMegabytev4), err)
		fmt.Printf("Время выполнение запросов %.2f сек. \n", time.Now().Sub(t1).Seconds())
	}
}

func requestSumAsyncWithCtxAndPool(urls []string) (int64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	ansCh := make(chan respStCWP, len(urls))
	semaphore := make(chan struct{}, poolHTTPReq)

	for _, url := range urls {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(u string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
			if err != nil {
				ansCh <- respStCWP{lenBody: 0, err: err}
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				ansCh <- respStCWP{lenBody: 0, err: err}
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				ansCh <- respStCWP{lenBody: 0, err: err}
				return
			}

			ansCh <- respStCWP{lenBody: int64(len(body)), err: nil}
		}(url)
	}

	go func() {
		wg.Wait()
		close(ansCh)
		close(semaphore)
	}()

	var sum int64
	var err error
	for bodyLen := range ansCh {
		sum += bodyLen.lenBody
		if bodyLen.err != nil && !errors.Is(bodyLen.err, context.Canceled) {
			if err != nil {
				err = fmt.Errorf("Ошибка %v у сайта %v;%v", bodyLen.err, bodyLen.lenBody, err)
			} else {
				err = fmt.Errorf("Ошибка %v у сайта %v", bodyLen.err, bodyLen.lenBody)
			}
			cancel()
		}
	}
	return sum, err
}
