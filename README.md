
# LitTime - Оценка Времени Чтения

**LitTime** — это инструмент для оценки времени, необходимого для прочтения текстового документа. Программа использует различные параметры, такие как скорость чтения, наличие визуальных элементов, и позволяет как из строки, так и в интерактивном режиме оценивать время на чтение документов.

## Навигация
- [Демонстрация](#демонстрация)
- [Основные возможности](#основные-возможности)
- [Установка и запуск](#установка-и-запуск)
- [Параметры командной строки](#параметры-командной-строки)
- [Конфигурация](#конфигурация)
- [Результаты](#результаты)
- [Интерактивный режим](#интерактивный-режим)
- [Зависимости](#зависимости)
- [Лицензия](#лицензия)


## Демонстрация
Вы можете посмотреть видео демонстрацию работы программы:
https://github.com/user-attachments/assets/dfcc016f-7a65-45ca-b1d7-b15f93b3e4fc



## Основные возможности

- Оценка времени чтения на основе скорости чтения пользователя.
- Учет наличия визуальных элементов (например, картинок) в тексте, что может увеличивать время чтения.
- Поддержка параллельной обработки текста для ускорения вычислений.
- Интерактивный режим для удобного выбора параметров без необходимости указывать их через командную строку.
- Поддержка русского и английского языков.

## Установка и запуск

1. **Клонируйте репозиторий:**

   ```bash
   git clone https://github.com/yourusername/littime.git
   cd littime
   ```

2. **Установка зависимостей:**

   Убедитесь, что у вас установлен Go. Для установки зависимостей используйте команду:

   ```bash
   go mod tidy
   ```

3. **Запуск программы:**

   Чтобы запустить программу, используйте следующую команду:

   ```bash
   go run main.go run --file yourfile.txt
   ```

   Вы можете также включить интерактивный режим:

   ```bash
   go run main.go run --interactive
   ```

## Параметры командной строки

- `--file` (`-f`) — Путь к текстовому файлу, который нужно проанализировать.
- `--speed` (`-s`) — Скорость чтения в словах в минуту (по умолчанию — 180).
- `--visuals` (`-v`) — Учитывать наличие визуальных элементов (устанавливайте как `true` для этого).
- `--workers` (`-w`) — Количество горутин для параллельной обработки (по умолчанию — 4).
- `--interactive` (`-i`) — Включение интерактивного режима для ввода параметров через интерфейс.

Пример:

```bash
go run main.go run --file yourfile.txt --speed 200 --visuals true --workers 6
```

## Конфигурация

Конфигурация проекта загружается из файла `config.yaml`, который может быть размещен в текущей директории или другой, указанной в коде. Если файл не найден, программа использует значения по умолчанию.

Пример файла конфигурации:

```yaml
default_reading_speed: 180
default_workers: 4
output_file: littime_results.json
```

## Результаты

После выполнения программы результат будет сохранен в указанный файл, например, `littime_results.json`, в формате JSON. Пример результата:

```json
{
  "ReadingTime": 12.34,
  "WordCount": 2500,
  "SentenceCount": 120,
  "SyllableCount": 4000,
  "FleschKincaidIndex": 72.5
}
```

## Интерактивный режим

В интерактивном режиме программа предоставляет удобный интерфейс для ввода необходимых параметров. Пользователь может последовательно ввести путь к файлу, скорость чтения, информацию о наличии визуальных элементов, а также количество потоков для обработки. Это можно посмотреть в [демонстрации](#демонстрация).

## Зависимости

- [Cobra](https://github.com/spf13/cobra) — для работы с CLI.
- [Viper](https://github.com/spf13/viper) — для загрузки конфигурации.
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — для создания интерактивного терминального интерфейса.




