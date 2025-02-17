#!/bin/bash

# install curl unzip nginx
apt-get update && apt-get install -y curl unzip nginx

if [ ! -d "trojan-install" ]; then
    mkdir trojan-install
fi

# Change to the "trojan-install" directory
cd trojan-install

# Define the destination directory
DEST_DIR="/usr/local/bin"

# Define the URL for the lego tarball
LEGO_URL="https://github.com/go-acme/lego/releases/download/v4.21.0/lego_v4.21.0_linux_amd64.tar.gz"

# Download the tarball
curl -L $LEGO_URL -o lego.tar.gz

# Extract the tarball
tar -xzf lego.tar.gz

# Move the lego binary to the destination directory
mv lego $DEST_DIR

# Clean up
rm lego.tar.gz

# Verify installation
$DEST_DIR/lego --version

# 安装 trojan-go service
# Define the URL for the trojan-go tarball
TROJAN_GO_URL="https://github.com/p4gefau1t/trojan-go/releases/download/v0.10.6/trojan-go-linux-amd64.zip"

# Download the tarball
curl -L $TROJAN_GO_URL -o trojan-go.zip

# Extract the tarball
unzip trojan-go.zip

mv trojan-go $DEST_DIR
# Verify installation
$DEST_DIR/trojan-go --version

DOMAIN="t.dl.tizizz.top"
WILDCARD_DOMAIN="*.dl.tizizz.top"
CONFIG="/etc/trojan-go/config.json"

if [ ! -d $(dirname ${CONFIG}) ]; then
    mkdir -p $(dirname ${CONFIG})
fi

# lego generate tls cert
EMAIL="janeysesions@gmail.com"
# ALICLOUD_ACCESS_KEY=LTAI5uFjRX4DJpm5 \
# ALICLOUD_SECRET_KEY=your-secret-key
# lego --email="${EMAIL}" --domains="${DOMAIN}" --domains="${WILDCARD_DOMAIN}" --http run

# if use dns provider, you can generate wildcard DOMAIN certs
# lego --email="${EMAIL}" --domains="${DOMAIN}" --dns alidns run

DOMAIN_CERT="$(pwd)/.lego/certificates/${DOMAIN}.crt"
DOMAIN_KEY="$(pwd)/.lego/certificates/${DOMAIN}.crt"

# write trojan-go config
cat > ${CONFIG} <<EOF
{
    "run_type": "server",
    "local_addr": "0.0.0.0",
    "local_port": 443,
    "remote_addr": "127.0.0.1",
    "remote_port": 80,
    "password": [
        "15660811",
        "a1b2c3",
        "e5f6g7",
        "i9j0k1",
        "m3n4o5",
        "q7r8s9",
        "u1v2w3",
        "y5z6a7",
        "c9d0e1",
        "g3h4i5",
        "k7l8m9"
    ],
    "ssl": {
        "cert": "${DOMAIN_CERT}",
        "key": "${DOMAIN_KEY}",
        "sni": "${DOMAIN}"
    }
}
EOF

SERVICE="/etc/systemd/system/trojan-go.service"

# write trojan-go.service
cat > ${SERVICE} <<EOF
[Unit]
Description=Trojan-Go - An unidentifiable mechanism that helps you bypass GFW
Documentation=https://p4gefau1t.github.io/trojan-go/
After=network.target nss-lookup.target

[Service]
User=root
CapabilityBoundingSet=CAP_NET_ADMIN CAP_NET_BIND_SERVICE
AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE
NoNewPrivileges=true
ExecStart=/usr/local/bin/trojan-go -config /etc/trojan-go/config.json
Restart=on-failure
RestartSec=10s
LimitNOFILE=infinity

[Install]
WantedBy=multi-user.target
EOF

# reload systemd and start trojan-go
systemctl daemon-reload && systemctl enable trojan-go && systemctl start trojan-go



