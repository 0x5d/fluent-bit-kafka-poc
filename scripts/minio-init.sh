#!/bin/bash

mc alias set myminio http://minio:9000 $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD
mc mb --ignore-existing --region=sa-east-1 myminio/access-logs
