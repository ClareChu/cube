#!/usr/bin/env bash

cd ../vendor/k8s.io/code-generator/

chmod o+x generate-groups.sh


./generate-groups.sh all \
hidevops.io/mio/pkg/client \
hidevops.io/mio/pkg/apis \
mio:v1alpha1