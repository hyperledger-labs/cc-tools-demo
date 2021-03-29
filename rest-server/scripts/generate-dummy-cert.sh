usage() {
    echo "usage: ./generate-dummy-cert.sh -d DOMAIN"
}

if [ $# -eq 0 ]
  then
    usage
fi

while getopts ":d:" opt; do
    case $opt in
        d) DOMAIN=$OPTARG;;
        *)
            usage
            exit ;;
    esac
done

sed "s/DOMAIN/${DOMAIN}/g" ../template/app-template.conf > ../data/nginx/app.conf


domains=($DOMAIN)
rsa_key_size=4096
data_path="../data/certbot"
email="test@goledger.com.br" # Adding a valid address is strongly recommended
staging=0 # Set to 1 if you're testing your setup to avoid hitting request limits

if [ ! -e "$data_path/conf/options-ssl-nginx.conf" ] || [ ! -e "$data_path/conf/ssl-dhparams.pem" ]; then
  echo "### Downloading recommended TLS parameters ..."
  mkdir -p "$data_path/conf"
  curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/options-ssl-nginx.conf > "$data_path/conf/options-ssl-nginx.conf"
  curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot/ssl-dhparams.pem > "$data_path/conf/ssl-dhparams.pem"
  echo
fi

echo "### Creating dummy certificate for $domains ..."
path="/etc/letsencrypt/live/$domains"
mkdir -p "$data_path/conf/live/$domains"
docker-compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:1024 -days 1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot
echo

echo 'Need root access to change ssl cert ownership'
set -x
sudo chown -R $USER:$USER ..
