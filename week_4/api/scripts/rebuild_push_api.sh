#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
API_DIR="$ROOT_DIR/week_4/api"
KUSTOMIZATION_FILE="$ROOT_DIR/week_4/manifests/kustomization.yaml"

IMAGE_REPO="registry.onstackit.cloud/level-3-paas/paas-api"
TAG=""
UPDATE_KUSTOMIZE="false"
CONTAINER_CLI="podman"

usage() {
  cat <<USAGE
Usage: $(basename "$0") [options]

Build and push the week_4 API image.

Options:
  -t, --tag <tag>           Image tag to use. Default: git-<shortsha>-<timestamp>
  -r, --repo <image-repo>   Image repo. Default: $IMAGE_REPO
  -c, --container-cli <cli> Container CLI to use. Default: $CONTAINER_CLI
  -u, --update-kustomize    Update week_4/manifests/kustomization.yaml newTag to pushed tag
  -h, --help                Show this help

Examples:
  $(basename "$0")
  $(basename "$0") --tag v4
  $(basename "$0") --container-cli podman --tag v4
  $(basename "$0") --tag v4 --update-kustomize
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -t|--tag)
      TAG="${2:-}"
      shift 2
      ;;
    -r|--repo)
      IMAGE_REPO="${2:-}"
      shift 2
      ;;
    -c|--container-cli)
      CONTAINER_CLI="${2:-}"
      shift 2
      ;;
    -u|--update-kustomize)
      UPDATE_KUSTOMIZE="true"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$TAG" ]]; then
  GIT_SHA="$(git -C "$ROOT_DIR" rev-parse --short HEAD 2>/dev/null || echo nogit)"
  TS="$(date +%Y%m%d%H%M%S)"
  TAG="git-${GIT_SHA}-${TS}"
fi

IMAGE="$IMAGE_REPO:$TAG"

if ! command -v "$CONTAINER_CLI" >/dev/null 2>&1; then
  echo "Container CLI not found: $CONTAINER_CLI" >&2
  exit 1
fi

echo "Building $IMAGE"
"$CONTAINER_CLI" build -t "$IMAGE" "$API_DIR"

echo "Pushing $IMAGE"
"$CONTAINER_CLI" push "$IMAGE"

if [[ "$UPDATE_KUSTOMIZE" == "true" ]]; then
  if [[ ! -f "$KUSTOMIZATION_FILE" ]]; then
    echo "Kustomization file not found: $KUSTOMIZATION_FILE" >&2
    exit 1
  fi

  awk -v repo="$IMAGE_REPO" -v tag="$TAG" '
    $0 ~ "name:[[:space:]]*" repo { print; getline; sub(/newTag:[[:space:]]*.*/, "newTag: " tag); print; next }
    { print }
  ' "$KUSTOMIZATION_FILE" > "$KUSTOMIZATION_FILE.tmp"

  mv "$KUSTOMIZATION_FILE.tmp" "$KUSTOMIZATION_FILE"
  echo "Updated kustomization tag to $TAG"
fi

echo "Done. Image: $IMAGE"
