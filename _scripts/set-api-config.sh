#!/bin/bash

usage="\n$(basename "$0") creates the config.json file for kubelens-api. assumes variables are available as environment variables
  Usage $0
  -h, --help      show this help text
  -f, --file      the file to overwrite\n\n"

while [[ $# -gt 0 ]]; do
  opt="$1"
  shift;
  current_arg="$1"
  case "$opt" in
    "-f"|"--file"     ) file=$1; shift;;
    "-h"|"--help"     ) printf "$usage"; exit 0;;
    ":"               ) printf "missing argument for -%s\n" "$2" >&2 printf "$usage" >&2 exit 1;;
    *                 ) echo "ERROR: Invalid option: \""$opt"\"" >&2 printf "$usage" >&2 exit 1;;
  esac
done

cat > $file <<- EOF
{
  "serverPort": ${SERVER_PORT},
  "allowedOrigins": [${ALLOWED_ORIGINS}],
  "allowedMethods": [
    "GET",
    "OPTIONS"
  ],
  "allowedHeaders": [
    "Accept",
    "Content-Type",
    "Authorization",
    "Pragma",
    "Cache-Control",
    "X-Requested-With"
  ],
  "allowedHosts": [${ALLOWED_HOSTS}],
  "oAuthJwk": "https://${OAUTH_JWT_ISSUER}/.well-known/jwks.json",
  "oAuthJwtIssuer": "https://${OAUTH_JWT_ISSUER}/",
  "oAuthAudience": "${OAUTH_AUDIENCE}",
  "enableAuth": ${ENABLE_OAUTH},
  "enableRBAC": ${ENABLE_RBAC},
  "rbacClaim": "${RBAC_CLAIM}",
  "serviceNameRegex": "${SERVICE_NAME_REGEX}",
  "projectSlugRegex": "${PROJECT_SLUG_REGEX}",
  "deployerLink": "${DEPLOYER_LINK}",
  "healthRoute": "${HEALTH_ROUTE}",
  "defaultSearchLabels": [${DEFAULT_SEARCH_LABELS}],
  "contentSecurityPolicy": "${CONTENT_SECURITY_POLICY}",
  "enableTLS": ${ENABLE_TLS},
  "tlsCert": "${TLS_CERT}",
  "tlsKey": "${TLS_KEY}",
  "enableCertPinning": ${ENABLE_CERT_PINNING},
  "publicKeyHPKP": "${PUBLIC_KEY_HPKP}"
}
EOF

echo "Overwrite config at $file"

cat $file
