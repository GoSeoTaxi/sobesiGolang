// Банальный синхронный вариант

package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const byteInMegabyte = 1024 * 1024

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
		byteSum, err := requesSumm(urlsList1)
		fmt.Printf("Сумма страниц в Мб=%.2f, ошибка - %v \n", (float64(byteSum) / byteInMegabyte), err)
		fmt.Printf("Время выполнение запросов %.2f сек. \n", time.Now().Sub(t1).Seconds())
	}
	fmt.Println("++++++++")
	{
		t1 := time.Now()
		byteSum, err := requesSumm(urlsList2)
		fmt.Printf("Сумма страниц в Мб=%.2f, ошибка - %v \n", (float64(byteSum) / byteInMegabyte), err)
		fmt.Printf("Время выполнение запросов %.2f сек. \n", time.Now().Sub(t1).Seconds())
	}
}

func requesSumm(urlsSlv []string) (int64, error) {

	var sum int64

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, v := range urlsSlv {
		resp, err := client.Get(v)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		sum += int64(len(body))

	}
	return sum, nil
}
