FROM php:8.3-cli

# Устанавливаем необходимые зависимости для Symfony
RUN apt-get update && apt-get install -y \
    libpq-dev \
    git \
    unzip

# Устанавливаем PHP-расширения для PostgreSQL
RUN docker-php-ext-install pdo pdo_pgsql

# Установка Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

# Задаем рабочую директорию, которая будет использоваться для Symfony
WORKDIR /var/www/symfony

# Копируем файлы проекта Symfony в контейнер
COPY ./symfony /var/www/symfony

# Устанавливаем зависимости проекта Symfony
RUN composer install --no-dev --optimize-autoloader

# Команда для запуска встроенного PHP-сервера Symfony
CMD php -S 0.0.0.0:8000 -t public
