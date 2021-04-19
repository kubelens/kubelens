 #!/bin/bash

go version
        
# curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# make -s test
echo "mode: count" > coverage-all.out

for pkg in $(go list ./... | grep -v -e /vendor/ -e "fakes")
do
  go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg} || exit 1
  tail -n +2 coverage.out >> coverage-all.out
done

COVERAGE=$(go tool cover -func=coverage-all.out | tail -1 | tr -d '[:space:]' | tr -d '()' | tr -d '%' | tr -d ':' | sed -e 's/total//g' | sed -e 's/statements//g')

if [[ ${COVERAGE%.*} -lt 65 ]]; then 
  echo "Insufficient Test Coverage: ${COVERAGE}"
  exit 1
else
  echo "Total Coverage: ${COVERAGE}"
fi