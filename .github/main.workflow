workflow "Deploy API" {
  resolves = [
    "Rollout API server",
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
  args = ["build -t stephenhillier/earthworks-api:$(echo ${GITHUB_SHA} | head -c7) ./api/"]
}

action "Docker login" {
  needs = ["Build image"]
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Push to registry" {
  uses = "actions/docker/cli@04185cf"
  needs = ["Docker login"]
  args = "push stephenhillier/earthworks-api:$(echo ${GITHUB_SHA} | head -c7)"
}

action "Get DO kubeconfig" {
  needs = ["Push to registry"]
  uses = "digitalocean/action/doctl@master"
  secrets = ["DIGITALOCEAN_ACCESS_TOKEN"]
  env = {
    CLUSTER_NAME = "island"
  }
  args = ["kubernetes cluster kubeconfig show $CLUSTER_NAME > $HOME/.kube/config"]
}

action "Apply deployment config" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["Get DO kubeconfig"]
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
  args = "web"
}

action "web - build image" {
  needs = ["Filter for web folder"]
  uses = "actions/docker/cli@04185cf"
  args = ["build -t stephenhillier/earthworks-web:$(echo ${GITHUB_SHA} | head -c7) ./web/"]
}

action "web - docker login" {
  needs = ["web - build image"]
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "web - Push to registry" {
  uses = "actions/docker/cli@04185cf"
  needs = ["web - docker login"]
  args = "push stephenhillier/earthworks-web:$(echo ${GITHUB_SHA} | head -c7)"
}

action "web - Get DO kubeconfig" {
  needs = ["web - Push to registry"]
  uses = "digitalocean/action/doctl@master"
  secrets = ["DIGITALOCEAN_ACCESS_TOKEN"]
  env = {
    CLUSTER_NAME = "island"
  }
  args = ["kubernetes cluster kubeconfig show $CLUSTER_NAME > $HOME/.kube/config"]
}

action "web - apply k8s/Istio config" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["web - Get DO kubeconfig"]
  runs = "sh -l -c"
  args = ["SHORT_REF=$(echo ${GITHUB_SHA} | head -c7) && cat kubernetes/pipeline/web.istio.yaml | sed 's/IMAGE_VERSION/'\"$SHORT_REF\"'/' | kubectl apply -f - "]
}

action "web - deployment status" {
  uses = "docker://gcr.io/cloud-builders/kubectl"
  needs = ["web - apply k8s/Istio config"]
  args = "rollout -n earthworks status deploy/earthworks-web"
}
