#!/usr/bin/env bash

cd ../vendor/k8s.io/code-generator/

chmod o+x generate-groups.sh


./generate-groups.sh all \
hidevops.io/cube/pkg/client \
hidevops.io/cube/pkg/apis \
cube:v1alpha1