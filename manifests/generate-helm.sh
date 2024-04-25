#!/bin/sh

export FLAGOPS_ENVIRONMENT=${ARGOCD_ENV_FLAGOPS_ENVIRONMENT}
export FLAGOPS_SOURCE_PATH="."
export FLAGOPS_IDENTITY=${ARGOCD_APP_NAME}

export FLAGOPS_DESTINATION_TYPE="file"
export FLAGOPS_DESTINATION_PATH=$(mktemp -d)

>&2 echo "DEBUG PWD: $(pwd)"
>&2 echo "DEBUG LS: $(ls)"

# Run Generation
>&2 echo "INFO Running FlagOps"
flagops --use-env

>&2 echo "DEBUG LS: $(ls ${FLAGOPS_DESTINATION_PATH})"

>&2 echo "INFO Inflating rendered templates"
helm template $ARGOCD_APP_NAME -n $ARGOCD_APP_NAMESPACE ${FLAGOPS_DESTINATION_PATH}