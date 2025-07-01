#!/bin/bash
docker-compose rm -f &&
docker-compose build &&
docker-compose up &&
docker image tag fold-vitest-fold ognerezov/fold-test-node:0.3