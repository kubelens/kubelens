 #!/bin/bash

usage="\n$(basename "$0") build & push a docker image
  Usage $0
  -h, --help      show this help text
  -a, --app       name of the application
  -b, --branch    git working branch
  -t, --tag       tag to use
  -i, --id        docker id
  -u, --user      docker username\n\n"

while [[ $# -gt 0 ]]; do
  opt="$1"
  shift;
  current_arg="$1"
  case "$opt" in
    "-a"|"--app"        ) app=$1; shift;;
    "-b"|"--branch"     ) branch=$1; shift;;
    "-t"|"--tag"        ) tag=$1; shift;;
    "-i"|"--id"         ) id=$1; shift;;
    "-u"|"--user"       ) user=$1; shift;;
    "-h"|"--help"       ) printf "$usage"; exit 0;;
    ":"                 ) printf "missing argument for -%s\n" "$2" >&2 printf "$usage" >&2 exit 1;;
    *                   ) echo "ERROR: Invalid option: \""$opt"\"" >&2 printf "$usage" >&2 exit 1;;
  esac
done

docker version

echo ${DOCKER_PASS} | docker login --username ${user} --password-stdin

docker build -t ${id}/${app}:${tag} .

if [ -z "${DOCKER_PASS}" ]; then
  echo "please provide docker password"
  read -s DOCKER_PASS
fi

docker push ${id}/${app}:${tag}

# tag latest and push on master branch
if [[ ${branch} -eq "master" ]]; then
  docker tag ${id}/${app}:${tag} ${id}/${app}:latest

  docker push ${id}/${app}:latest
fi
