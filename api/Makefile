default: test

depends:
	dep init || dep ensure

.PHONY: test
test:
	sh ../_scripts/go-test-cover.sh

.PHONY: docker-build-push
docker-build-push:
	sh ../_scripts/docker-build-push.sh --app kubelens-api --branch ${GIT_BRANCH} --tag ${TAG} --id ${DOCKER_ID} --user ${DOCKER_USER}

set-config: 
	sh ../_scripts/set-api-config.sh --file ./_helm/kubelens-api/files/config.json

helm-upgrade:
	helm upgrade kubelens-api --install --recreate-pods --reset-values --force --set ingress.host=${INGRESS_HOST} ./_helm/kubelens-api
	