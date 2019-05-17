#!/bin/bash

#set -x

RED='\033[1;91m'
WT='\033[1;97m'
NC='\033[0m'

kubectl delete ksvc prime-ui prime-api
clear

echo -e "${RED}The simple HTTP API backend${WT}\n"

cat http_service.yaml

read  -p "====="

ko apply -f ./http_service.yaml

read  -p "====="

echo -e "${RED}The simple HTTP UI Client${WT}\n"

cat ../prime-client/ui.yaml

read -p "====="

ko apply -f ../prime-client/ui.yaml

echo -e "${RED}Examine the Knative Services...${WT}\n"

kubectl get ksvc

read -p "====="
echo -n "Probing..."

r=$(curl http://prime-ui.dangerd.net -o /dev/null -s --write-out '%{http_code}')

while [ $r -ne "200" ]; do
  r=$(curl http://prime-ui.dangerd.net -o /dev/null -s --write-out '%{http_code}')
done

echo -e " done.\n"

read -p "====="

echo -e "${RED}We developed a new negative backend${WT}\n"

cat ./http_service3.yaml

read -p "====="

echo -e "${RED}Deploying new version, 60/40 split${WT}\n"

ko apply -f ./http_service3.yaml

read -p "====="

echo -e "${RED}Now, let us see how this app scales${WT}\n"

read -p "====="

hey -t 90 -n 100000 -c 350 "http://prime-ui.dangerd.net/prime?query=611101"&

read -p "===== hey launched..."

echo -e "${RED}Now we decided to develop a gRPC backend${WT}\n"

read -p "====="

cat ./grpc_service.yaml

read -p "====="

ko apply -f ./grpc_service.yaml

read -p "====="

echo -e "${RED}Checking the endpoints...${WT}\n"

kubectl get ksvc prime-api -oyaml | grep "url:"

read -p "====="

echo -e "${RED}The load balanced HTTP and gRPC UI Client${WT}\n"

cat ../prime-client/ui3.yaml

read -p "====="

ko apply -f ../prime-client/ui3.yaml

read -p "====="

echo -e "\n\t\t\t\t${NC}That's all folks!\n"

