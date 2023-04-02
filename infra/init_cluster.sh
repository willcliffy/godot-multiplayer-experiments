#!/bin/bash

kubectl config use-context do-nyc1-kilnwood-cluster

kubectl create secret generic registry-credentials \
  --from-file=.dockerconfigjson=secrets/dockerRegistryCredentials.json \
  --type=kubernetes.io/dockerconfigjson || exit

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.5.1/deploy/static/provider/do/deploy.yaml
