#!/bin/bash

usage="\n$(basename "$0") creates the config.json file for kubelens-api or kubelens-web. assumes variables are available as environment variables
  Usage $0
  -h, --help          show this help text
  -t, --type   the config to overwrite, valid values are 'web, api'
  -f, --file          the file to overwrite\n\n"

while [[ $# -gt 0 ]]; do
  opt="$1"
  shift;
  current_arg="$1"
  case "$opt" in
    "-t"|"--type" ) type=$1; shift;;
    "-f"|"--file" ) file=$1; shift;;
    "-h"|"--help" ) printf "$usage"; exit 0;;
    ":"           ) printf "missing argument for -%s\n" "$2" >&2 printf "$usage" >&2 exit 1;;
    *             ) echo "ERROR: Invalid option: \""$opt"\"" >&2 printf "$usage" >&2 exit 1;;
  esac
done

if [[ "${type}" -eq "web" ]]; then
cat > $file <<- EOF
{
  "oAuthJwtIssuer": "${OAUTH_JWT_ISSUER}",
  "oAuthAudience": "${OAUTH_AUDIENCE}",
  "oAuthClientId": "${OAUTH_CLIENT_ID}",
  "oAuthRedirectUri": "${OAUTH_REDIRECT_URI}",
  "oAuthResponseType": "${OAUTH_RESPONSE_TYPE}",
  "oAuthRequestType": "${OAUTH_REQUEST_TYPE}",
  "oAuthScope": "${OAUTH_SCOPE}",
  "oAuthProvider": "${OAUTH_PROVIDER}",
  "oAuthConnection": "${OAUTH_CONNECTION}",
  "availableClusters": [
    ${AVAILABLE_CLUSTERS}
  ],
  "deployerLinkName": "${DEPLOYER_LINK_NAME}"
}
EOF
fi

if [[ "${type}" -eq "api" ]]; then
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
  "oAuthJwk": "${OAUTH_JWK}",
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
fi

echo "Overwrite config at $file"

cat $file
