# Netology

## Coding Skills

Вам нужно написать небольшой сервис, который мы назовём Slowly.

### Основная функциональность

Сервис принимает запросы по протоколу HTTP (если сделаете поддержку HTTP2 - будет плюсом).

Запрос должен выглядеть следующим образом:

```
POST /api/slow
Content-Type: application/json

{
    "timeout": 3000 // время в мс
}
```

Ответ приходит через указанное количество времени (поле timeout):
```
HTTP 200 OK
Content-Type: application/json

{
    "status": "ok"
}
```

На все другие запросы должен приходить ответ HTTP 404 Not Found

### Middleware

Поскольку пользователи начали злоупотреблять сервисом, было решено написать Middleware, которое "отрубает" обработчики, ожидающие больше 5 секунд. Реализуйте подобное Middleware: для всех подобных ошибок, должен возвращаться ответ:
```
HTTP 400 Bad Request
Content-Type: application/json

{
    "error": "timeout too long"
}
```

### Результаты

Напишите автотесты на Happy Path (на код 200) и Sad Path (на код 400). Организация и используемые инструменты - на ваше усмотрение
Выложите весь код на Github (при необходимости комментарии можете написать в README.md)
Подключите какой-нибудь облачный CI, который будет запускать тесты - Github Actions, Travis CI или любой другой на ваше усмотрение
Пришлите нам ссылку на ваш репо

Важно: вы сами решаете как организовать код, как покрывать его тестами (с учётом требований) и т.д. - представьте, что вы готовите пример, чтобы объяснить человеку, как (с вашей точки зрения) "делать правильно".
