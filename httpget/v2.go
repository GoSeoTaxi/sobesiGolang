// Банальный ассинхронный вариант
package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const byteInMegabytev2 = 1024 * 1024

type respSt struct {
	lenBody int64
	err     error
}

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
		byteSum, err := requesSummAsync(urlsList1)
		fmt.Printf("Сумма страниц в Мб=%.2f, ошибка - %v \n", (float64(byteSum) / byteInMegabytev2), err)
		fmt.Printf("Время выполнение запросов %.2f сек. \n", time.Now().Sub(t1).Seconds())
	}
	fmt.Println("++++++++")
	{
		t1 := time.Now()
		byteSum, err := requesSummAsync(urlsList2)
		fmt.Printf("Сумма страниц в Мб=%.2f, ошибка - %v \n", (float64(byteSum) / byteInMegabytev2), err)
		fmt.Printf("Время выполнение запросов %.2f сек. \n", time.Now().Sub(t1).Seconds())
	}
}

func requesSummAsync(urls []string) (int64, error) {
	var wg sync.WaitGroup
	ansCh := make(chan respSt, len(urls))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			resp, err := client.Get(u)
			if err != nil {
				ansCh <- respSt{
					lenBody: 0,
					err:     err,
				}
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				ansCh <- respSt{
					lenBody: 0,
					err:     err,
				}
				return
			}
			ansCh <- respSt{
				lenBody: int64(len(body)),
				err:     nil,
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(ansCh)
	}()

	var sum int64
	var err error
	for bodyLen := range ansCh {
		sum += bodyLen.lenBody
		if bodyLen.err != nil {
			if err == nil {
				err = fmt.Errorf("Ошибка %v у сайта %v", bodyLen.err)
				continue
			}
			err = fmt.Errorf("Ошибка %v у сайта %v;%v", bodyLen.err, err)
		}
	}
	if err != nil {
		return 0, err
	}

	return sum, err
}
