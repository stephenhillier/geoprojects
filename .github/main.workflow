workflow "Build & deploy to GKE" {
  resolves = ["Rollout API server"]
  on = "push"
}

action "Filter for API folder" {
  uses = "docker://gcr.io/cloud-builders/git"
  runs = "sh -l -c"
  args = ["git diff-tree --name-only HEAD | grep api"]
}

action "Build image" {
  needs = ["Filter for API folder"]
  uses = "actions/docker/cli@04185cf"
  args = ["build -t gcr.io/islandcivil-223001/earthworks-api:$(echo ${GITHUB_SHA} | head -c7) ./api/"]
}

action "Setup gcloud" {
  uses = "actions/gcloud/auth@8ec8bfa"
  needs = ["Build image"]
  secrets = ["GCLOUD_AUTH"]
}

action "GKE Credentials" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["Setup gcloud"]
  args = "container clusters get-credentials earthworks --zone us-west1-a --project islandcivil-223001"
}

action "GKE Docker" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["GKE Credentials"]
  args = "auth configure-docker -q"
}

action "Push to GCR.io" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["GKE Docker"]
  args = "docker -- push gcr.io/islandcivil-223001/earthworks-api:$(echo ${GITHUB_SHA} | head -c7)"
}

action "Apply deployment config" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["Push to GCR.io"]
  runs = "sh -l -c"
  args = ["SHORT_REF=$(echo ${GITHUB_SHA} | head -c7) && cat kubernetes/pipeline/api.istio.yaml | sed 's/IMAGE_VERSION/'\"$SHORT_REF\"'/' | kubectl apply -f - "]
}

action "Rollout API server" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["Apply deployment config"]
  args = "rollout -n earthworks status deploy/earthworks-api"
}
