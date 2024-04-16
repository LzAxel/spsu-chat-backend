# REST API для чата

## Стартануть проект 

```bash
# Запускаем БД в докере
$ docker-compose up -d

# Запускаем приложение
$ make run
```

## Документация к API

### Авторизация

#### Регистрация

```properties
POST /auth/sign-up
```

**Тело запроса**
```json
{
    "username": "",
    "display_name": "",
    "password": ""
}
```

**Тело ответа**
```properties
Статус 201 - при успешной регистрации
```
---
#### Вход

```properties
POST /auth/sign-in
```

**Тело запроса**
```json
{
    "username": "",
    "password": ""
}
```

**Тело ответа**
```json
{
    "access_token": "",
    "refresh_token": ""    
}
```
---
#### Обновление токенов

```properties
POST /auth/refresh
```

**Тело запроса**
```json
{
    "refresh_token": ""
}
```
**Тело ответа**
```json
{
    "access_token": "",
    "refresh_token": ""    
}
```

### Операции с юзерами

#### Получить всех юзеров **(Для админов!)**

```properties
GET /users
```

**Тело ответа**
```json
{
    "users": [
        {
            "id": 0,
            "username": "",
            "display_name": "",
            "type": 0,
            "created_at": "",
        },
        ...
    ],
    "pagination": {
        "page": 1,
        "page_limit": 20,
        "page_count": 1,
        "total": 1
    }
}
```
---
#### Получить своего юзера

```properties
GET /users/self
```

**Тело ответа**
```json
{
    "users": {
        "id": 0,
        "username": "",
        "display_name": "",
        "type": 0,
        "created_at": ""
    }
}
```
---
#### Получить юзера по ID

```properties
GET /users/:id
```
Вместо `:id` в запросе подставить числовой ID юзера

**Тело ответа**
```json
{
    "user": {
        "id": 0,
        "username": "",
        "display_name": "",
        "type": 0,
        "created_at": ""
    }
}
```