# Вступление
Привет участник хакатона!

Репозиторий - https://github.com/ScDewt/hackathon_10_24
Изначально этот репозиторий - это заготовка для разработки сервисов на разных языках - php / go / python / nodejs / etc.

Внутри репозитория есть соответствующие каталоги для этих языков.  
Внутри каталога каждого языка есть уже некоторые предустановленные библиотеки/фреймворки и сделан апи метод /api/health, где проверяется подключение к БД.

Внутри [docker-compose.yml](docker-compose.yml) ты найдешь порты по которым опубликованы эти языки и их api.

# Python
```bash
curl -XGET localhost:8080/api/health
curl -XGET http://hackathon.scdewt.ru:8080/api/health
```

# Php
```bash
curl -XGET localhost:8081/api/health
curl -XGET http://hackathon.scdewt.ru:8081/api/health
```

# Go
```bash
curl -XGET localhost:8082/api/health
curl -XGET http://hackathon.scdewt.ru:8082/api/health
```

# NodeJS
```bash
curl -XGET localhost:8083/api/health
curl -XGET http://hackathon.scdewt.ru:8083/api/health
```

<br><br><br>

**Удачного участия и хорошо провести время!**

<br><br><br>