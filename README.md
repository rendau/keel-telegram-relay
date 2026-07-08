# keel-telegram-relay

Маленький сервис-мостик: принимает [keel](https://keel.sh) webhook-уведомления и пересылает
их в Telegram. Нужен потому, что у keel нет нативного Telegram-нотификатора — только
generic webhook, а Telegram Bot API требует свой формат (`chat_id` + `text`).

## Env

| Переменная         | Обязательна | По умолчанию              | Описание                                          |
|--------------------|-------------|---------------------------|---------------------------------------------------|
| `TELEGRAM_TOKEN`   | да          | —                         | токен бота (`@BotFather`)                          |
| `TELEGRAM_CHAT_ID` | да          | —                         | id чата/канала/группы для уведомлений             |
| `HTTP_PORT`        | нет         | `80`                      | порт HTTP-сервера                                 |
| `TELEGRAM_API_URL` | нет         | `https://api.telegram.org`| переопределение endpoint (например, через прокси) |
| `DEBUG`            | нет         | `false`                   | text-логи вместо json                             |

## HTTP

- `POST /` — endpoint для keel webhook. Принимает JSON keel-события, шлёт в Telegram.
- `GET /healthcheck` — healthcheck.

## Подключение keel

В values keel-чарта (`keel/keel`) включить webhook на адрес этого сервиса:

```yaml
notificationLevel: info
webhook:
  enabled: true
  endpoint: "http://keel-telegram-relay.default.svc.cluster.local/"
```

Keel шлёт события вида:

```json
{
  "name": "update deployment",
  "message": "Successfully updated deployment default/oms (ghcr.io/mechta-market/oms:1.2.3)",
  "type": "deployment update",
  "level": "success",
  "createdAt": "2026-07-08T10:00:00Z"
}
```

## Локальный запуск

```sh
cp .env.example .env   # заполнить TELEGRAM_TOKEN / TELEGRAM_CHAT_ID
make
./cmd/build/svc
```

## CI/CD

`.github/workflows/deploy.yml` на push в `master` собирает бинарь, пакует в образ и
пушит `ghcr.io/rendau/keel-telegram-relay:latest`.
