# Frontend — Development setup (Flutter)

Коротко: как запустить клиентскую часть проекта (Flutter) для разработки.

## Предварительные требования

- Flutter SDK (последняя стабильная версия). Установка: https://flutter.dev/docs/get-started/install
- Android SDK / Xcode (по целям платформ)
- VS Code или Android Studio (рекомендуется)

## Установка зависимостей и запуск

1. Перейдите в папку `frontend`:

```powershell
cd d:\Documents\Projects\9_Golang\Orders\frontend
```

2. Установите зависимости:

```powershell
flutter pub get
```

3. Настройте API_BASE_URL (если приложение использует константу или переменную окружения). В коде ищите место, где формируются URL для API (например, в `lib/services`).

4. Запустите приложение на эмуляторе или девайсе:

```powershell
flutter run
```

## Аутентификация и тестирование

- Используйте страницу логина (`0_0_login.dart`) для входа через backend `/login`.
- Сохраняйте возвращаемый JWT (в коде это может делать `shared_preferences`) и используйте в заголовках следующих запросов.

## Советы

- Если вы меняете API-адрес на лету, удобно хранить его в `lib/config.dart` или использовать flavour/environment.
- Чтобы дебажить network-запросы, используйте `Flutter DevTools` или логирование в сервисах.
