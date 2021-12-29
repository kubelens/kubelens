default: test

depends:
	go mod tidy

.PHONY: test
test:
	go version

	echo "mode: count" > coverage-all.out
	
	$(foreach pkg,$(shell go list ./... | grep -v -e /vendor/ -e "fakes"), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg} || exit 1; \
		tail -n +2 coverage.out >> coverage-all.out;)

	COVERAGE=$$(go tool cover -func=coverage-all.out | tail -1 | tr -d '[:space:]' | tr -d '()' | tr -d '%' | tr -d ':' | sed -e 's/total//g' | sed -e 's/statements//g'); \
	if [ $$(echo "$${COVERAGE} < 65" | bc) -eq 1 ]; then \
		echo "Insufficient Test Coverage: $${COVERAGE}"; \
		exit 1; \
	else \
		echo "Total Coverage: $${COVERAGE}"; \
	fi

.PHONY: docker-build-push
docker-build-push:
	bash ../_scripts/docker-build-push.sh --app kubelens-api --branch ${GIT_BRANCH} --tag ${TAG} --id ${DOCKER_ID} --user ${DOCKER_USER}

helm-upgrade:
	helm upgrade kubelens-api --install --set ingress.host=${INGRESS_HOST} ./_helm/kubelens-api
	