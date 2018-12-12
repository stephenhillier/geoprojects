workflow "New workflow" {
  on = "push"
  resolves = ["Build image"]
}

action "Build image" {
  uses = "actions/docker/cli@04185cf"
  args = "build -t islandcivil-223001/earthworks-api ./api/"
}
