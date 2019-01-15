workflow "Deploy API" {
  resolves = [
    "Rollout API server"
  ]
  on = "push"
}

# backend API pipeline

action "Filter for API folder" {
  uses = "netlify/actions/diff-includes@exit-code-78"
  args = "api"
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

# Frontend pipeline


workflow "Deploy web" {
  on = "push"
  resolves = ["web - deployment status"]
}


action "Filter for web folder" {
  uses = "netlify/actions/diff-includes@exit-code-78"
  args = "api"
}

action "web - build image" {
  needs = ["Filter for web folder"]
  uses = "actions/docker/cli@04185cf"
  args = ["build -t gcr.io/islandcivil-223001/earthworks-web:$(echo ${GITHUB_SHA} | head -c7) ./api/"]
}

action "web - setup gcloud" {
  uses = "actions/gcloud/auth@8ec8bfa"
  needs = ["web - build image"]
  secrets = ["GCLOUD_AUTH"]
}

action "web - GKE Credentials" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["web - setup gcloud"]
  args = "container clusters get-credentials earthworks --zone us-west1-a --project islandcivil-223001"
}

action "web - GKE Docker" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["web - GKE Credentials"]
  args = "auth configure-docker -q"
}

action "web - Push to GCR.io" {
  uses = "actions/gcloud/cli@8ec8bfa"
  needs = ["web - GKE Docker"]
  args = "docker -- push gcr.io/islandcivil-223001/earthworks-web:$(echo ${GITHUB_SHA} | head -c7)"
}

action "web - apply k8s/Istio config" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["web - Push to GCR.io"]
  runs = "sh -l -c"
  args = ["SHORT_REF=$(echo ${GITHUB_SHA} | head -c7) && cat kubernetes/pipeline/web.istio.yaml | sed 's/IMAGE_VERSION/'\"$SHORT_REF\"'/' | kubectl apply -f - "]
}

action "web - deployment status" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["web - apply k8s/Istio config"]
  args = "rollout -n earthworks status deploy/earthworks-web"
}
