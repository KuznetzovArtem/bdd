FROM golang:1.20

# Создаем и переходим в рабочую директорию
WORKDIR /app

# Устанавливаем Git
RUN apt-get update && apt-get install -y git

# Клонируем репозиторий с GitHub
RUN git clone https://github.com/KuznetzovArtem/web-server.git

# Переходим в директорию с приложением
WORKDIR /app/web-server

RUN go mod vendor

WORKDIR /app/web-server/cmd
# Собираем приложение
RUN go build -o web-server


CMD ["./web-server"]

