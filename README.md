![CI](https://github.com/Koha90/github.com/koha90/shopcore/actions/workflows/ci.yml/badge.svg)

# github.com/koha90/shopcore

E-commerce платформа с мульти-интерфейсами:
- Telegram (primary)
- TUI
- Web (future)

## Архитектура

Clean Architecture + Ports/Adapters

## Быстрый старт

```bash
make up
make migrate
make run-tui
```

## Структура проекта

internal/
manager # lifecycle ботов
botconfig # бизнес-логика конфигов
app/
appkit # сборка приложения (composition root)
bootstrap # запуск ботов
seed # dev данные
pgapp # postgres config + pool
Entry points

cmd/
migrate
tui
Текущее состояние

    Postgres ✔

    миграции ✔

    seed ✔

    bootstrap ✔

    TUI ✔

