#!/bin/sh

export FLAGOPS_ENVIRONMENT=${ARGOCD_ENV_FLAGOPS_ENVIRONMENT}
export FLAGOPS_SOURCE_PATH="."
export FLAGOPS_IDENTITY=${ARGOCD_APP_NAME}

export FLAGOPS_DESTINATION_TYPE="file"
export FLAGOPS_DESTINATION_PATH=$(mktemp -d)

# Run Generation
>&2 echo "Executing FlagOps"
flagops --use-env

>&2 echo "Building Dependencies"
>&2 helm dependency build ${FLAGOPS_DESTINATION_PATH}

>&2 echo "Infalting Helm Chart"
helm template $ARGOCD_APP_NAME -n $ARGOCD_APP_NAMESPACE -a $KUBE_API_VERSIONS ${FLAGOPS_DESTINATION_PATH}