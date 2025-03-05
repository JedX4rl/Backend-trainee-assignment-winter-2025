# Merch store

## Описание

Сервис для внутреннего магазина мерча с возможностью:
- Автоматической регистрацией новых пользователей
- Просмотра истории операций
- Покупки товаров за монеты
- Передачи монет между пользователями

## Запуск проекта

1. Клонировать репозиторий:
```bash
git clone https://github.com/JedX4rl/Backend-trainee-assignment-winter-2025.git
cd Backend-trainee-assignment-winter-2025
```
2. Запустить сервисы:

```bash
docker-compose build
docker-compose up
```  

Сервис будет доступен по адресу:
```
http://localhost:8080
```  

## Нагрузочное тестирование
<img width="803" alt="image" src="https://github.com/user-attachments/assets/021f816f-1005-4314-a715-15a32cc260d0" />

Код тестирования:
```
func main() {

	rate := vegeta.Rate{Freq: 1000, Per: time.Second} 
	duration := 10 * time.Second                    
	authBody := `{
    "toUser": "aboba2",
    "amount": 1
}`
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/api/buy/pen",
		Body:   []byte(authBody),
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFib2JhMiIsImlkIjoiMiIsImV4cCI6MTczOTc0MTA1NX0.19YqRZQUN6age8m3QHyeZd3RzA3DKVOe-nXr-rJNXns"},
		},
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99-й перцентиль: %s\n", metrics.Latencies.P99)
	fmt.Printf("Средняя задержка: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Процент успешных запросов: %.2f%%\n", metrics.Success*100)
}
```
