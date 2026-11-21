#!/bin/sh
# Генерация сертификатов для задания 4: свой CA, сертификат сервера и клиента.
# Нужен OpenSSL 3 (на Mac: brew install openssl, если системный не справится).
# Запуск из папки certs: sh gen-certs.sh
set -e
DAYS=365

# 1. Свой удостоверяющий центр (CA): ключ + самоподписанный корневой сертификат
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out ca.key
openssl req -x509 -new -key ca.key -days $DAYS -subj "/CN=Lab10 CA" -out ca.crt \
  -addext "basicConstraints=critical,CA:TRUE" \
  -addext "keyUsage=critical,keyCertSign,cRLSign"

# 2. Сервер: ключ -> запрос на сертификат (CSR) -> CA подписывает
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out server.key
openssl req -new -key server.key -subj "/CN=localhost" -out server.csr
printf "subjectAltName=DNS:localhost,IP:127.0.0.1\nextendedKeyUsage=serverAuth\n" > server.ext
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -days $DAYS -extfile server.ext -out server.crt

# 3. Клиент: то же самое, но с правом clientAuth (для mTLS)
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out client.key
openssl req -new -key client.key -subj "/CN=lab10-client" -out client.csr
printf "extendedKeyUsage=clientAuth\n" > client.ext
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -days $DAYS -extfile client.ext -out client.crt

# Промежуточные файлы больше не нужны
rm -f server.csr client.csr server.ext client.ext ca.srl
echo "Готово: ca.crt ca.key server.crt server.key client.crt client.key"
