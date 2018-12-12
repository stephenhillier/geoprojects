workflow "Build & deploy to GKE" {
  on = "push"
  resolves = ["Rollout API server"]
}

action "Build image" {
  uses = "actions/docker/cli@04185cf"
  args = "build -t gcr.io/islandcivil-223001/earthworks-api ./api/"
}

action "Setup gcloud" {
  uses = "actions/gcloud/auth@8ec8bfa"
  needs = ["Build image"]
  secrets = ["GCLOUD_AUTH"]
}

action "GKE Credentials" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["Setup gcloud"]
  args = "container clusters get-credentials islandcivil --zone us-west1-a --project islandcivil-223001"
}

action "GKE Docker" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["GKE Credentials"]
  args = "auth configure-docker -q"
}

action "Push to GCR.io" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["GKE Docker"]
  args = "docker -- push gcr.io/islandcivil-223001/earthworks-api"
}

action "Apply deployment config" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["Push to GCR.io"]
  args = "apply -f kubernetes/02-api-deploy.yaml"
}

action "Rollout API server" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["Apply deployment config"]
  args = "rollout -n earthworks status deploy/earthworks-api"
}
