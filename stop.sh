#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# Shut down the Docker containers that might be currently running.
docker-compose -f docker-compose.yml stop
yes|docker-compose -f docker-compose.yml rm
docker rm -f $(docker ps -aqf name=dev)
docker rmi $(docker images --format '{{.Repository}}'|grep 'dev-')
docker rm -f $(docker ps -aqf name=peer)
docker rm -f $(docker ps -aqf name=couchdb)
