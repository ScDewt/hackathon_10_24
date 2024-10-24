sudo apt update
sudo apt upgrade -y

echo "Install dependencies Docker"
sudo apt install apt-transport-https ca-certificates curl software-properties-common -y
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "Install Docker"
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io -y

echo "Check Docker"
sudo systemctl status docker

echo "Install Docker Compose"
sudo curl -L "https://github.com/docker/compose/releases/download/v2.21.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

echo "Check Docker Compose"
docker-compose --version

echo "Create user hackathon"
sudo adduser hackathon
sudo usermod -aG sudo hackathon
sudo usermod -aG docker hackathon



mkdir -p /var/www/html/symfony

echo "Start AutoDeploy"
chmod +x check-and-reload.sh
touch /var/log/pull-and-reload.log
(crontab -l 2>/dev/null; echo "* * * * * /bin/bash -c 'for i in {1..12}; do /var/www/hackathon/system/check-and-reload.sh; sleep 5; done' >> /var/log/pull-and-reload.log 2>&1") | crontab -




sudo reboot
