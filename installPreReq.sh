printf "Installing Building Artefacts\n\n"

sudo apt -y update

sudo apt -y upgrade

sudo apt install build-essential

printf "\n\nInstalling GoLang\n"

sudo curl -fsSL https://golang.org/dl/go1.15.12.linux-amd64.tar.gz --output go1.15.12.linux-amd64.tar.gz

sudo rm -rf /opt/go

sudo tar -C /opt -xvzf go1.15.12.linux-amd64.tar.gz

printf "\n\nInstalling NodeJs\n"

curl -sL https://deb.nodesource.com/setup_10.x -o nodesource_setup.sh

sudo bash nodesource_setup.sh

sudo apt install nodejs

printf "\n\nInstalling Docker\n"

curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh

printf "\n\nReseting Docker\n"

sudo usermod -aG docker $(whoami)

sudo grpck

sudo grpconv

newgrp docker << END

sudo systemctl restart docker.service

printf "\n\nInstalling Docker-Compose\n"

sudo curl -L "https://github.com/docker/compose/releases/download/1.29.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose

printf "\n\nInstalling Hyperledger Fabric binaries and images\n"

curl -sSL http://bit.ly/2ysbOFE | bash -s -- 1.4.12 1.4.9 0.4.22

printf "\n\nCustomizing enviroment variables\n"

echo "export PATH=$PATH:/opt/go/bin:$HOME/fabric-samples/bin" >> ~/.bashrc

echo "export GOPATH=$HOME" >> ~/.bashrc

echo "export GOROOT=/opt/go" >> ~/.bashrc

source ~/.bashrc

printf "\n\nYou should logout now\n"

exit

END