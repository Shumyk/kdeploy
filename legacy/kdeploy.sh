#!/bin/zsh
#
# kdeploy - deploy from the terminal to Kubernetes.
#           searches for images of requested microservice in Google Container Registry,
#           prompts interactively you to select an image for deployment (arrows navigation, search features),
#           sets the selected image in workload, and conveniently shows you watch over deploying microservices.
#
#           also has a deploy previous mode - when you need quickly redeploy what was before your last deployment.
#           however, it has goldfish memory - can redeploy only the previous deployment.

### CONSTANTS
##
DOMAIN="your-domain"
STATEFULSET_MS="your-microservice"
GCR_URL="us.gcr.io/your-url"

SELECTED_IMAGE_FILE="tmp.selected.image"
PREVIOUS_DEPLOYMENT_FILE="${HOME}/.kdeploy"
###

### FUNCTIONS
##
# creates watch over deploying pods
watchPods() {
  watch -n1 "kubectl get pods | grep ${MICROSERVICE}"
}
# wraps arguments or stdin into color
greenish() {
  GREEN='\033[1;92m' DEFAULT='\033[0m'
  if [[ "${#}" -eq 0 ]]; then
    echo "${GREEN}$(cat)${DEFAULT}"
  else
    echo "${GREEN}${*}${DEFAULT}"
  fi
}
purplish() {
  PURPLE='\033[1;35m' DEFAULT='\033[0m'
  echo "${PURPLE}${*}${DEFAULT}"
}
error() {
  BOLD_RED='\033[1;31m' DEFAULT='\033[0m'
  echo "${BOLD_RED}${*}${DEFAULT}" >&2
}
# prints horizontal line
horizontalLine() {
  printf %"${COLUMNS}"s |tr " " "-"
}
# wraps input into horizontal lines
wrapHorizontalLines() {
  ON_GREEN='\033[42m' BOLD_BLACK='\033[1;30m' DEFAULT='\033[0m'
  horizontalLine
  echo "${ON_GREEN}${BOLD_BLACK}${1}${DEFAULT}"
  horizontalLine
}
# general info printing
printEnvironmentInfo() {
  wrapHorizontalLines "|     ENVIRONMENT        |"
  printf              "service                  :   %s\n" "$(greenish "${MICROSERVICE}")"
  printf              "namespace                :   %s\n" "$(greenish "${NAMESPACE}")"
}
# parses image path and prints tag and digest
printImageInfo() {
  printf "tag                      :   %s\n" "$(echo "${1}" | awk -F: '{print $2}' | awk -F@ '{print $1}' | greenish)"
  printf "digest                   :   %s\n" "$(echo "${1}" | awk -F: '{print $NF}' | greenish)"
  horizontalLine
}
# sets image and watch over deployment
kubectlSetImage() {
  kubectl set image "${WORKLOAD}/${NAMESPACE}-${MICROSERVICE}" "${MICROSERVICE}=${NEW_IMAGE}" 1> /dev/null

  wrapHorizontalLines  "|     DEPLOYED IMAGE     |"
  printImageInfo       "${NEW_IMAGE}"

  watchPods
}
# validates previous deployment file and redeploys it
deployPrevious() {
  if [[ ! -f "${PREVIOUS_DEPLOYMENT_FILE}" ]]; then
    error "No previous deployment file found"
    return 1
  fi
  if [[ $(grep -vc ^$ < "${PREVIOUS_DEPLOYMENT_FILE}") -ne 4 ]]; then
    error "Previous deployment file corrupted: ${PREVIOUS_DEPLOYMENT_FILE}"
    error "should contain exactly four lines"
    return 1
  fi

  WORKLOAD=$(     sed '1q;d' "${PREVIOUS_DEPLOYMENT_FILE}")
  NAMESPACE=$(    sed '2q;d' "${PREVIOUS_DEPLOYMENT_FILE}")
  MICROSERVICE=$( sed '3q;d' "${PREVIOUS_DEPLOYMENT_FILE}")
  NEW_IMAGE=$(    sed '4q;d' "${PREVIOUS_DEPLOYMENT_FILE}")
  printEnvironmentInfo

  rm "${PREVIOUS_DEPLOYMENT_FILE}"
  kubectlSetImage
  return 0
}
# saves current deployment into previous deployment file
savePreviousDeployment() {
  echo "${WORKLOAD}"         >"${PREVIOUS_DEPLOYMENT_FILE}"
  echo "${NAMESPACE}"       >>"${PREVIOUS_DEPLOYMENT_FILE}"
  echo "${MICROSERVICE}"    >>"${PREVIOUS_DEPLOYMENT_FILE}"
  echo "${CURRENT_IMAGE}"   >>"${PREVIOUS_DEPLOYMENT_FILE}"
}
# resolves workload resource type, as some microservices could be stateful sets
# if you have more than one stateful set - alter condition
resolveWorkload() {
  if [[ "${MICROSERVICE}" == "${STATEFULSET_MS}" ]]; then
    echo "statefulset"
  else
    echo "deployment"
  fi
}
# if tag absent semicolon should be omitted
appendSemicolon() {
  if [[ ${1} ]]; then
    echo ":${1}"
  fi
}
###


### SCRIPT BEGINNING
##
if [[ "${#}" -eq 0 ]]; then
  wrapHorizontalLines "DEFAULT DEPLOY MODE"                                                            >&2
  echo "usage:        " "$(greenish "${0} <microservice> [size]")"                                     >&2
  echo "                          microservice     - name"                                             >&2
  echo "                          size             - optional, amount of images to fetch from GCP"     >&2
  wrapHorizontalLines "PREVIOUS DEPLOY MODE"                                                           >&2
  echo "usage:        " "$(greenish "${0} previous")"                                                  >&2
  return 1
fi

# deploy previous mode - different routine
if [[ "${1}" == "previous" ]]; then
  deployPrevious
  return "${?}"
fi

# parameters
MICROSERVICE="${1}"
SHOW_IMAGES="${2:-20}"

NAMESPACE=$(kubectl config view --minify --output 'jsonpath={..namespace}'; echo)
WORKLOAD=$(resolveWorkload)
printEnvironmentInfo

# gets currently deployed image
CURRENT_IMAGE=$(kubectl get "${WORKLOAD}" "${NAMESPACE}-${MICROSERVICE}" -o jsonpath="{..image}" | tr -s "[:space:]" "\n" | grep "${DOMAIN}-${MICROSERVICE}")
wrapHorizontalLines  "|CURRENTLY DEPLOYED IMAGE|"
printImageInfo       "${CURRENT_IMAGE}"
# additional line so go program will eat it
horizontalLine

# gets images from GCR in specific format
IMAGES_LIST=$(gcloud alpha container images list-tags "${GCR_URL}${MICROSERVICE}" --show-occurrences-from="${SHOW_IMAGES}" --format="value[separator=|](timestamp, digest, tags)")

# runs go program to interactively select image
select_image_prompt "${SELECTED_IMAGE_FILE}" $(echo "${IMAGES_LIST}")
SELECTED_IMAGE_INFO=($(cat "${SELECTED_IMAGE_FILE}"))
rm "${SELECTED_IMAGE_FILE}"

# terminate execution if nothing selected
if [[ "${#SELECTED_IMAGE_INFO[@]}" -eq 0 ]]; then
  echo
  purplish "heh, ctrl+C combination was gently pressed. see you"
  return 0
fi

# SELECTED_IMAGE_INFO contents:
#     1: timestamp, not needed anymore
#     2: short digest of selected image, full digest should be resolved
#     3: image tag, optional
# Note: arrays start from 1 in .zsh, modify for other shells
IMAGE_TAG=$(appendSemicolon "${SELECTED_IMAGE_INFO[3]}")
SHORT_DIGEST="${SELECTED_IMAGE_INFO[2]}"
FULL_DIGEST=$(gcloud container images describe "${GCR_URL}${MICROSERVICE}@sha256:${SHORT_DIGEST}" --format="value(image_summary.digest)")

NEW_IMAGE="${GCR_URL}${MICROSERVICE}${IMAGE_TAG}@${FULL_DIGEST}"

savePreviousDeployment
kubectlSetImage
return 0
