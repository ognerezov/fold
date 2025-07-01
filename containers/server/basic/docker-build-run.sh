#!/bin/bash
docker-compose rm -f &&
docker-compose build &&
docker-compose up
docker image tag fold-be-fold ognerezov/fold