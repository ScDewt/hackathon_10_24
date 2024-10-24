#!/bin/bash

# Путь к репозиторию
REPO_PATH="/var/www/hackathon"
DOCKER_COMPOSE_PATH="docker-compose.yml"

# Переходим в директорию с репозиторием
cd $REPO_PATH

# Сохраняем старый хеш последнего коммита
LAST_COMMIT=$(git rev-parse HEAD)

# Выполняем git pull
git pull origin main

# Сохраняем новый хеш последнего коммита
NEW_COMMIT=$(git rev-parse HEAD)

# Сравниваем старый и новый хеши
if [ "$LAST_COMMIT" != "$NEW_COMMIT" ]; then
    echo "Изменения найдены. Перезапуск Docker контейнеров..."
    docker-compose -f $DOCKER_COMPOSE_PATH pull
    docker-compose -f $DOCKER_COMPOSE_PATH up -d --build --remove-orphans
else
    echo "Изменений нет. Контейнеры не будут перезапущены."
fi
