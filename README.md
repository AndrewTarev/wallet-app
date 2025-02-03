# Wallet-app

## Основной функционал
1. Создание кошелька
2. Изменение баланса
3. Получение баланса

## Структура проекта
```
.
├── cmd
├── docs                       // Документация swagger
├── internal        
│   ├── app
│   │   ├── app_errors  
│   │   ├── delivery           // Слой хэндлеров
│   │   │   └── http
│   │   ├── domain             // Сущности
│   │   ├── repository         // Слой работы с БД
│   │   └── services           // Слой с бизнес логикой
│   │       └── mocks
│   ├── configs
│   └── infrastructure         // Настройки подключения к БД, логгера, сервера
│       ├── database
│       │   └── migrations     // Миграции
│       ├── logger
│       └── server
└── test                       // Тесты
```

## Установка приложения:

1. Склонируйте репозиторий себе на компьютер
    - git clone https://github.com/AndrewTarev/wallet-app.git

2. Установите свои переменные в .env файл


3. Запустите сборку контейнеров
    - docker-compose up --build

API документация (Swagger/OpenAPI) доступна по пути http://localhost:8080/swagger/index.html