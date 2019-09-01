#!/usr/bin/env bash
set -e

if [[ ${TRAVIS_EVENT_TYPE} == "cron" ]]; then
    exit 0
fi

IMAGE_NAME="messagebird/sachet"

docker build --build-arg TAG="${TRAVIS_TAG}" -t ${IMAGE_NAME}:latest .

docker tag  ${IMAGE_NAME}:latest  ${IMAGE_NAME}:"${TRAVIS_TAG}"

echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin

docker push  ${IMAGE_NAME}:"${TRAVIS_TAG}"
docker push  ${IMAGE_NAME}:latest
