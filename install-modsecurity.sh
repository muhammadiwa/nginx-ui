#!/bin/bash

# Exit on error
set -e

echo "Installing ModSecurity for NGINX..."

# Clean up any previous installation files
echo "Cleaning up previous installation files..."
rm -rf /tmp/ModSecurity /tmp/ModSecurity-nginx /tmp/nginx-1.24.0 /tmp/nginx-1.24.0.tar.gz

# Install dependencies
echo "Installing dependencies..."
sudo apt-get update
sudo apt-get install -y \
    autoconf automake build-essential \
    libcurl4-openssl-dev liblmdb-dev libpcre3-dev \
    libxml2-dev libyajl-dev libtool libtool-bin \
    libgeoip-dev pkgconf zlib1g-dev git wget

# Build ModSecurity
echo "Building ModSecurity..."
cd /tmp
git clone --depth 1 https://github.com/SpiderLabs/ModSecurity
cd ModSecurity
git submodule init
git submodule update
./build.sh
./configure
make
sudo make install

# Build ModSecurity NGINX connector
echo "Building ModSecurity NGINX connector..."
cd /tmp
git clone --depth 1 https://github.com/SpiderLabs/ModSecurity-nginx

# Get NGINX source matching our installed version
NGINX_VERSION=$(nginx -v 2>&1 | grep -oP '(?<=nginx/)[0-9]+\.[0-9]+\.[0-9]+')
echo "Using NGINX version: $NGINX_VERSION"

wget "http://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz"
tar xzf "nginx-${NGINX_VERSION}.tar.gz"
cd "nginx-${NGINX_VERSION}"

# Configure and build the module
./configure \
    --prefix=/usr/share/nginx \
    --conf-path=/etc/nginx/nginx.conf \
    --http-log-path=/var/log/nginx/access.log \
    --error-log-path=/var/log/nginx/error.log \
    --lock-path=/var/lock/nginx.lock \
    --pid-path=/run/nginx.pid \
    --modules-path=/usr/lib/nginx/modules \
    --with-compat \
    --with-http_ssl_module \
    --add-dynamic-module=/tmp/ModSecurity-nginx

make modules

# Install the module
sudo mkdir -p /usr/lib/nginx/modules
sudo cp objs/ngx_http_modsecurity_module.so /usr/lib/nginx/modules/

# Create ModSecurity configuration directory
sudo mkdir -p /etc/nginx/modsecurity

# Copy ModSecurity configuration files
sudo cp /tmp/ModSecurity/modsecurity.conf-recommended /etc/nginx/modsecurity/modsecurity.conf
sudo cp /tmp/ModSecurity/unicode.mapping /etc/nginx/modsecurity/

# Create ModSecurity module configuration
sudo mkdir -p /usr/share/nginx/modules-available
echo "load_module modules/ngx_http_modsecurity_module.so;" | sudo tee /usr/share/nginx/modules-available/mod-modsecurity.conf

# Enable the module
sudo ln -sf /usr/share/nginx/modules-available/mod-modsecurity.conf /etc/nginx/modules-enabled/

# Create main ModSecurity configuration
echo 'Include "/etc/nginx/modsecurity/modsecurity.conf"

SecRuleEngine On
SecRequestBodyAccess On
SecResponseBodyAccess On
SecResponseBodyMimeType text/plain text/html text/xml application/json
SecDataDir /tmp' | sudo tee /etc/nginx/modsecurity/main.conf

# Update NGINX configuration
if ! grep -q "modsecurity on;" /etc/nginx/nginx.conf; then
    sudo sed -i '/http {/a \    modsecurity on;\n    modsecurity_rules_file /etc/nginx/modsecurity/main.conf;' /etc/nginx/nginx.conf
fi

# Set proper permissions
echo "Setting permissions..."
sudo chown -R www-data:www-data /etc/nginx/modsecurity
sudo chmod -R 644 /etc/nginx/modsecurity/*
sudo chmod 755 /etc/nginx/modsecurity

# Test NGINX configuration
echo "Testing NGINX configuration..."
sudo nginx -t

echo "ModSecurity installation completed!"
echo "Now you can restart NGINX with: sudo systemctl restart nginx"
