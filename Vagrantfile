# -*- mode: ruby -*-
# vi: set ft=ruby :

$script = <<-'SCRIPT'
# Set DNS
sudo sed -i '/^DNS=/c\DNS=8.8.8.8' "/etc/systemd/resolved.conf"
sudo sed -i '/^FallbackDNS=/c\FallbackDNS=1.1.1.1' "/etc/systemd/resolved.conf"
sudo systemctl restart systemd-resolved

# Upgrade system
sudo apt update
sudo apt upgrade -y
sudo apt install -y ca-certificates git curl

# Install docker
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
# Add Docker's official GPG key:
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Install mattermost on port 8065
(
  cd /opt/
  git clone https://github.com/mattermost/docker mattermost
)
(
  cd /opt/mattermost
  cp env.example .env
  sed -i '/^DOMAIN=/c\DOMAIN=192.168.56.4' /opt/mattermost/.env
  sed -i '/^MM_SERVICESETTINGS_SITEURL=/c\MM_SERVICESETTINGS_SITEURL=http://${DOMAIN}' /opt/mattermost/.env

  mkdir -p ./volumes/app/mattermost/{config,data,logs,plugins,client/plugins,bleve-indexes}
  sudo chown -R 2000:2000 ./volumes/app/mattermost

  sudo docker compose -f docker-compose.yml -f docker-compose.without-nginx.yml up -d
)

# Install sentry
(
  cd /opt/
  VERSION="24.2.0"
  git clone https://github.com/getsentry/self-hosted.git sentry
  cd sentry
  git checkout ${VERSION}
  sudo REPORT_SELF_HOSTED_ISSUES=0 SKIP_USER_CREATION=1 ./install.sh
  sed -i '/^system.internal-url-prefix:/c\system.internal-url-prefix: '"'"'http://192.168.56.4:9000'"'"'' sentry/config.yml
  sed -i '/^\(# \)\?CSRF_TRUSTED_ORIGINS =/c\CSRF_TRUSTED_ORIGINS = ["http://192.168.56.4:9000"]' sentry/sentry.conf.py
  docker compose up -d
)
SCRIPT

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/jammy64"
  config.vm.box_version = "20240131.0.0"

  config.vm.network "private_network", ip: "192.168.56.4"

  config.vm.provision "shell", inline: $script

  config.vm.provider "virtualbox" do |vb|
    vb.cpus = "12"
    vb.memory = "16384"
  end
end
