FROM php:8.1-fpm

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y \
    libpq-dev \
    git \
    unzip \
    nginx

# Устанавливаем расширения PHP
RUN docker-php-ext-install pdo pdo_pgsql

# Установка Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

# Копируем файлы приложения
WORKDIR /var/www/html
COPY . .

# Устанавливаем зависимости Symfony
RUN composer install

# Копируем конфигурацию Nginx
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf

# Команды для запуска PHP-FPM и Nginx
CMD service nginx start && php-fpm
