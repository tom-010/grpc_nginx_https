#!/bin/bash

mkcert -install
mkcert -cert-file=certs/tls.crt -key-file=certs/tls.key example.com "*.example.com"
