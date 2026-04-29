![CI](https://github.com/Koha90/shopcore/actions/workflows/ci.yml/badge.svg)

# shopcore

[English](README.md) · [Русский](README.ru.md)

**shopcore** это e-commerce платформа с несколькими интерфейсами и общим ядром продаж, каталога и bot runtime.

Это не просто менеджер ботов. Проект развивается как платформа, где Telegram, TUI и будущий Web используют общую доменную логику и application services.

---

## Текущий фокус

- Telegram как основной клиентский канал продаж
- TUI как текущая рабочая панель оператора/админа
- Web как будущая админка и витрина
- общий каталог, заказы и flow для разных интерфейсов
- runtime-конфигурация и база данных отдельно для каждого бота

---

## Интерфейсы

### Telegram

Telegram сейчас основной runtime adapter.

Реализовано:

- клиентская навигация по каталогу
- reply и inline клавиатуры
- сценарии старта через `/start`
- schema-driven навигация по каталогу
- оформление заказа
- сохранение заказа
- уведомления о заказах в админский чат
- действия админа по статусу заказа
- уведомления админу об обычных сообщениях клиента
- ответы клиенту из карточки сообщения
- ответы клиенту из карточки заказа
- текстовые ответы админа
- ответы админа с фото и optional caption
- отдельный flow service для каждого бота
- catalog provider для каждого бота через `database_id`

### TUI

TUI сейчас используется как операторская/админская панель.

Реализовано:

- список ботов
- статус runtime
- запуск/остановка/рестарт
- редактирование токена
- редактирование конфигурации
- отображение и редактирование `StartScenario`
- синхронизация runtime spec через `manager.UpdateSpec(...)`

### Web

Планируется.

Ожидаемое направление:

- редактирование каталога
- управление заказами
- инструменты поддержки клиентов
- web admin/storefront поверх общих сервисов

---

## Архитектура

shopcore следует clean-ish ports/adapters подходу.

Основные принципы:

- flow не зависит от транспорта
- зависимости явные
- application services небольшие
- read/write стороны разделяются там, где это полезно
- SQL не попадает в Telegram runtime
- Telegram-специфика не попадает в flow
- runtime wiring отдельный для каждого бота
- тесты добавляются вместе с логикой
- двигаемся маленькими шагами, без больших переписываний

Основные пакеты:

```text
internal/
  app/
    bootstrap/        bootstrap приложения и runtime wiring
    pgapp/            wiring Postgres pool/config
    runtime/
      telegram/       Telegram runtime adapter
      demo/           demo runtime части
    seed/             demo seed data
    tuiapp/           wiring TUI-приложения

  botconfig/          домен/сервис конфигурации ботов
  catalog/
    postgres/         Postgres catalog adapter
    service/          catalog write-side application services
  flow/               transport-agnostic navigation/input flow
  manager/            lifecycle manager ботов
  order/
    postgres/         Postgres order adapter
    service/          order application service
  tui/                terminal UI
  transport/          HTTP transport части
```

---

## Быстрый старт

```bash
make up
make migrate
make run-tui
```

Запуск тестов:

```bash
go test ./...
```

---

## Flow

`internal/flow` не зависит от конкретного транспорта.

Он отвечает за:

- start scenarios
- screen/view models
- reply и inline keyboard models
- состояние session
- back через history
- pending text/photo input
- навигацию по каталогу
- transport-neutral effects

Start scenarios:

```text
reply_welcome
inline_catalog
```

Навигация по каталогу строится по схеме:

```text
city -> category -> district -> product -> variant
```

Generic catalog actions and screens:

```text
catalog:select:<level>:<id>
catalog:screen:<path>
```

Back navigation использует `Session.History`. Он не прыгает в hardcoded root.

---

## Каталог

Каталог хранится в Postgres и загружается в `flow.Catalog` через `CatalogProvider`.

Текущие таблицы:

```text
cities
catalog_categories
catalog_city_categories
catalog_districts
catalog_products
catalog_variants
```

Runtime wiring:

```text
bot config
↓
database_id
↓
database profile
↓
Postgres pool
↓
catalog provider
↓
per-bot flow service
↓
Telegram runtime
```

Важные правила каталога:

- продукт без вариантов пропускается
- район без товаров пропускается
- категория без валидной ветки пропускается
- город без children пропускается

---

## Заказы

Order flow сейчас поддерживает:

- выбор catalog variant
- подтверждение заказа
- сохранение заказа
- уведомление в настроенный admin chat
- действия админа:
  - взять в работу
  - закрыть
- обновление admin card после callback actions
- ответ клиенту из карточки заказа

---

## Сообщения клиентов и ответы админа

Обычный клиентский текст обрабатывается Telegram runtime:

1. Если session имеет pending input, текст уходит в `flow.HandleText`.
2. Если pending input нет и отправитель это клиент, текст отправляется в настроенный admin chat.
3. Admin users не эхоятся в уведомления о сообщениях клиента.

Admin reply flow:

```text
сообщение клиента или карточка заказа
↓
админ нажимает reply action
↓
бот открывает pending input в личке админа
↓
админ отправляет текст или фото
↓
flow возвращает transport-neutral effect
↓
Telegram runtime отправляет ответ клиенту
```

Поддерживаются ответы:

- обычный текст
- фото с optional caption

Для фото используется transport-owned media token. Flow считает token opaque-значением и только передаёт его через `EffectSendPhoto`.

---

## Runtime

Поведение runtime:

- каждый бот имеет свой runtime spec
- каждый бот имеет свой `flow.Service`
- несколько ботов могут использовать одну базу
- несколько ботов могут использовать разные базы
- изменения config синхронизируются через `manager.UpdateSpec(...)`
- изменение token/config/start scenario не требует полного рестарта системы

Telegram metadata:

- Telegram bot id
- Telegram username
- Telegram bot display name

---

## Seed data

Seed создаёт:

- реальный database profile
- demo bot
- demo catalog data

Demo catalog включает:

- Москву
- Санкт-Петербург
- категории цветов и подарков
- районы
- товары
- варианты

Если token пустой, demo bot создаётся disabled.

---

## Что сейчас работает

- Telegram bot runtime
- `/start`
- reply welcome scenario
- inline catalog scenario
- загрузка каталога из Postgres
- schema-driven catalog traversal
- back navigation через history
- per-bot catalog wiring по `database_id`
- управление ботами из TUI
- runtime spec sync
- Telegram metadata sync через `getMe`
- persisted order flow
- admin order workflow
- уведомления о сообщениях клиента
- ответы админа текстом
- ответы админа с фото

---

## Статус разработки

Базовая проверка:

```bash
go test ./...
```

Ожидаемый результат: все тесты проходят.

Покрываемые области:

- `internal/flow`
- catalog Postgres loader/build/repository
- catalog service write use cases
- Telegram runtime contracts
- bot configuration
- order service and storage
- manager behavior

---

## Ближайший roadmap

### Catalog admin

- новые write use cases
- редактирование category/city/district/product/variant
- управление каталогом из TUI
- Telegram admin bot/catalog tools
- позже Web admin catalog editor

### Customer communication

- сохранение customer inquiries
- сохранение admin replies
- delivery status
- retry support
- message history
- расширенная работа с media

### Orders

- более богатые order cards
- комментарии к заказу
- история заказа
- улучшенный admin workflow
- customer-visible order status

### Web

- admin dashboard
- редактирование каталога
- управление заказами
- customer support view

---

## Принципы разработки

- не ломать уже работающее
- делать маленькие проверяемые изменения
- не заниматься premature overengineering
- держать зависимости явными
- оборачивать важные ошибки
- документировать важные package/type/function
- добавлять тесты вместе с логикой
- держать flow transport-agnostic
- не тащить SQL в runtime
- относиться к shopcore как к commerce platform, а не временному bot script
