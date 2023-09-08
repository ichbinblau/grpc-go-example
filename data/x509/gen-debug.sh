#!/bin/bash

rm -f *.key
rm -f *.crt
rm -f *.csr

# 1. Generate Root CA's private key and self-signed certificate
openssl genrsa -out root-ca.key 3072
openssl req -new -x509 -sha384 -days 365 -key root-ca.key -out root-ca.crt -subj "/C=CN/L=Shanghai/O=Intel/CN=Root"

# # 2. Generate Aggregator CA private key and certificate
# openssl genrsa -out aggregator-ca.key 3072
# openssl req -new -sha384 -key aggregator-ca.key -out aggregator-ca.csr -subj "/C=CN/L=Shanghai/O=Intel/CN=Aggregator"
# openssl x509 -req -in aggregator-ca.csr -CA root-ca.crt -CAkey root-ca.key -CAcreateserial -out aggregator-ca.crt -days 365 -sha384

# 2. Generate Aggregator server key and certificate
openssl genrsa -out aggr-server.key 3072
openssl req -new -sha384 -key aggr-server.key -out aggr-server.csr -subj "/C=CN/L=Shanghai/O=Intel/CN=AggregatorServer"
openssl x509 -req -in aggr-server.csr -CA root-ca.crt -CAkey root-ca.key -CAcreateserial -out aggr-server.crt -days 365 -sha384 -extfile debug.cnf

# # 4. Generate Node Agent CA private key and certificate
# openssl genrsa -out nodeagent-ca.key 3072
# openssl req -new -sha384 -key nodeagent-ca.key -out nodeagent-ca.csr -subj "/C=CN/L=Shanghai/O=Intel/CN=NA"
# openssl x509 -req -in nodeagent-ca.csr -CA root-ca.crt -CAkey root-ca.key -CAcreateserial -out nodeagent-ca.crt -days 365 -sha384

# 3. Generate Node Agent server key and certificate
openssl genrsa -out na-server.key 3072
openssl req -new -sha384 -key na-server.key -out na-server.csr -subj "/C=CN/L=Shanghai/O=Intel/CN=NodeAgentServer"
openssl x509 -req -in na-server.csr -CA root-ca.crt -CAkey root-ca.key -CAcreateserial -out na-server.crt -days 365 -sha384 -extfile debug.cnf

# 3. Generate scheduler plugin key and certificate
openssl genrsa -out plugin.key 3072
openssl req -new -sha384 -key plugin.key -out plugin.csr -subj "/C=CN/L=Shanghai/O=Intel/CN=SchedulerPlugin"
openssl x509 -req -in plugin.csr -CA root-ca.crt -CAkey root-ca.key -CAcreateserial -out plugin.crt -days 365 -sha384
