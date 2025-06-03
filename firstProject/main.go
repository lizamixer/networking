package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	mu     sync.RWMutex
	status string
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info(fmt.Sprintf("Сервер запущен на http://localhost:8080"))

	// Ticker обновляет значение каждую секунду
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker.C {
			mu.Lock()
			status = fmt.Sprintf("Текущее время: %s", t.Format("15:04:05"))
			mu.Unlock()
		}
	}()

	// Отдаём HTML-страницу
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Время</title></head>
		<body>
			<h1 id="time">Загрузка...</h1>
			<script>
				setInterval(() => {
					fetch('/time')
						.then(response => response.text())
						.then(text => {
							document.getElementById('time').innerText = text;
						});
				}, 1000);
			</script>
		</body>
		</html>
		`
		fmt.Fprint(w, html)
	})

	// Отдаём данные
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		fmt.Fprint(w, status)
		mu.RUnlock()
	})

	http.ListenAndServe(":8080", nil)
}
