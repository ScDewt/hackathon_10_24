# Подготовка и запуск
mkdir -p /var/www
cd /var/www
git clone https://github.com/ScDewt/hackathon_10_24.git hackathon
bash /var/www/system/init.sh


docker exec -it hackathon-postgresql-1 bash


# Читка кешей Docker
docker-compose down
docker builder prune --all
