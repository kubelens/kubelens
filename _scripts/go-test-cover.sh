 #!/bin/bash

usage="\n$(basename "$0") run go tests
  Usage $0
  -h, --help      show this help text\n\n"

while [[ $# -gt 0 ]]; do
  opt="$1"
  shift;
  current_arg="$1"
  case "$opt" in
    "-h"|"--help"       ) printf "$usage"; exit 0;;
    ":"                 ) printf "missing argument for -%s\n" "$2" >&2 printf "$usage" >&2 exit 1;;
    *                   ) echo "ERROR: Invalid option: \""$opt"\"" >&2 printf "$usage" >&2 exit 1;;
  esac
done

go version
        
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

if [ -z "${GOBIN}" ]; then
  $GOBIN="$(go env GOPATH)/bin"
fi

$GOBIN/dep init || $GOBIN/dep ensure

# make -s test
echo "mode: count" > coverage-all.out

for pkg in $(go list ./... | grep -v -e /vendor/ -e "fakes")
do
  go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg} || exit 1
  tail -n +2 coverage.out >> coverage-all.out
done

COVERAGE=$(go tool cover -func=coverage-all.out | tail -1 | tr -d '[:space:]' | tr -d '()' | tr -d '%' | tr -d ':' | sed -e 's/total//g' | sed -e 's/statements//g')

if [[ ${COVERAGE%.*} -lt 70 ]]; then 
  echo "Insufficient Test Coverage: ${COVERAGE}"
  exit 1
else
  echo "Total Coverage: ${COVERAGE}"
fi