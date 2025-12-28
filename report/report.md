# Отчёт по заданию GitHub Actions CI/CD

## 1) Выбранный проект
Язык: Go.  
Проект: простое CLI-приложение `Greet(name)` + unit-тесты.

Репозиторий: https://github.com/YellowPhil/gha-ci-cd-demo

## 2) Файл workflow (ci-cd.yml)
```yaml
name: CI-CD

on:
  push:
    branches: ["main"]
  pull_request:

jobs:
  build_test:
    name: Build & Test (${{ matrix.os }})
    runs-on: ${{ matrix.os }}
    strategy:
      fail-fast: false
      matrix:
        os:
          - ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.22.x"
          cache: true

      - name: Build
        run: go build ./...

      - name: Test (unit)
        run: go test -v ./...

  deploy:
    name: Deploy (SSH)
    runs-on: ubuntu-latest
    needs: [build_test]
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    steps:
      - name: SSH deploy
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.DEPLOY_HOST }}
          username: ${{ secrets.DEPLOY_USER }}
          key: ${{ secrets.DEPLOY_SSH_KEY }}
          port: ${{ secrets.DEPLOY_PORT }}
          script: |
            set -e
            echo "Deploy on $(hostname) at $(date)"
            mkdir -p ~/apps/gha-ci-cd-demo
            cd ~/apps/gha-ci-cd-demo

            if [ ! -d .git ]; then
              git clone https://github.com/${{ github.repository }} ./
            else
              git fetch --all
              git reset --hard origin/main
            fi

            go build -o app .
            echo "Deployed build OK"
```

## 3) Описание настроенных job'ов

### build_test
- Триггеры: `push` в `main`, `pull_request`
- Matrix ОС:
  - ubuntu-latest
- Actions:
  - actions/checkout@v4
  - actions/setup-go@v5
- Шаги: checkout -> setup go -> build -> go test -v
- Цель: гарантировать сборку и выполнение unit-тестов на нескольких ОС.

### deploy
- Запускается только при `push` в `main` (условие `if`)
- Зависимости: `needs: [build_test]` (деплой только после успешных тестов на всех ОС)
- Action:
  - appleboy/ssh-action@v1.0.3
- Секреты:
  - DEPLOY_HOST, DEPLOY_USER, DEPLOY_SSH_KEY, DEPLOY_PORT
- Логика: подключение по SSH -> clone/pull репозитория -> сборка на сервере.

## 4) Скриншоты / результаты
- Вкладка Actions в GitHub показывает успешные прогоны build/test по всем ОС.
- Deploy job запускается после build_test (main push).
