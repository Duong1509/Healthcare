**NOTE**
recommend for nodejs ver 12.12.0
---
## Start network and deploy CC

Start network
```sh
cd blockchain/script
./network.sh up createChannel -ca -s couchdb
```
deploy Chaincode
```sh
./network.sh deployCC
```

start MongoDB 
```sh
cd blockchain/docker
docker-compose -f docker-compose-mongo.yaml up -d
```

Enroll admin
```sh
cd handlers/users/
node EnrollAdmin.js
```
## Start server

Start api server
```sh
pm2 start server.js
```
