default: test

depends:
	dep init || dep ensure

.PHONY: test
test:
	bash ../_scripts/go-test-cover.sh

.PHONY: docker-build-push
docker-build-push:
	bash ../_scripts/docker-build-push.sh --app kubelens-api --branch ${GIT_BRANCH} --tag ${TAG} --id ${DOCKER_ID} --user ${DOCKER_USER}

helm-upgrade:
	helm upgrade kubelens-api --install --set ingress.host=${INGRESS_HOST} ./_helm/kubelens-api
	