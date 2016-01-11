# SMSC
Пакет для работы с smsc.ru
# Пример работы
### Отправка сообщения
```go
package main

import(
	"fmt"
	"github.com/mil-ast/smsc"
)

func main(){
    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
        }
    }()

    // Инициализация
	sms, err :=  smsc.New("https://smsc.ru/sys/send.php", "login", "password")
	if err != nil {
	    panic(err)
	}
	// добавление дополнительных параметров
	sms.AddParam("fmt", "3")
	
	// отправка сообщения.
	res, err := sms.Send("POST", []string{"+79876543210"}, "Тело сообщения")
	if err != nil {
		panic(err)
	}
	
	fmt.Println(res)
}
```

### Прочие запросы
```go
// запрос баланса
sms, err :=  smsc.New("http://smsc.ru/sys/balance.php", "login", "password")
if err != nil {
    panic(err)
}

res, err := sms.Request()
if err != nil {
    panic(err)
}
fmt.Println(res)
```