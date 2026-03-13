# API Documentation — PZ1

## Переменные окружения

| Сервис | Переменная | Дефолт |
|--------|-----------|--------|
| Auth | AUTH_PORT | 8081 |
| Tasks | TASKS_PORT | 8082 |
| Tasks | AUTH_BASE_URL | http://localhost:8081 |

## Auth Service (порт 8081)

### POST /v1/auth/login
Получение токена.
- **200** — `{"access_token": "demo-token", "token_type": "Bearer"}`
- **401** — неверные учётные данные

### GET /v1/auth/verify
Проверка токена. Заголовок: `Authorization: Bearer <token>`
- **200** — `{"valid": true, "subject": "student"}`
- **401** — токен невалиден

## Tasks Service (порт 8082)

Все запросы требуют заголовка `Authorization: Bearer <token>`.

### POST /v1/tasks — создать задачу (201)
### GET /v1/tasks — список задач (200)
### GET /v1/tasks/{id} — задача по ID (200/404)
### PATCH /v1/tasks/{id} — обновить задачу (200/404)
### DELETE /v1/tasks/{id} — удалить задачу (204/404)