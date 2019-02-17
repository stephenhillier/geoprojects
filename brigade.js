const {events, Job, Group} = require("brigadier");
const checkRunImage = "deis/brigade-github-check-run:latest"

const dest = "$GOPATH/src/github.com/stephenhillier/geoprojects";

events.on("check_suite:requested", runTests)
events.on("check_suite:created", runTests)
events.on("check_suite:rerequested", runTests)
events.on("check_run:rerequested", runTests)
events.on("exec", runTests)
events.on("pull_request", deployPullRequest)

function runTests(e, p) {
  var build = new Job("test", "golang:1.11")
  build.tasks = [
    "mkdir -p " + dest,
    "cp -a /src/* " + dest,
    "cd " + dest + "/api",
    "go get -u github.com/golang/dep/cmd/dep",
    "dep ensure",
    "go test"
  ];
  runJobWithCheck(e, p, "Build", "run tests", build)
}

function deployPullRequest(e, p) {
  var build = new Job("provision", "gcr.io/cloud-builders/kubectl")
  console.log(e.payload)
  build.tasks = [
    "NAME=earthworks-pr-" + e.payload['number'] + " && cat /src/operator/deploy/crds/earthworks_v1alpha1_earthworks_cr.yaml | sed 's/DEPLOYMENT_NAME/'\"$NAME\"'/' | cat",
    "NAME=earthworks-pr-" + e.payload['number'] + " && cat /src/operator/deploy/crds/earthworks_v1alpha1_earthworks_cr.yaml | sed 's/DEPLOYMENT_NAME/'\"$NAME\"'/' | kubectl apply -f -",
    "echo " + e.payload['action']
  ];
  build.run()
}

function runJobWithCheck(e, p, checkName, checkTitle, job) {
  console.log("check requested")
  // Common configuration
  const env = {
    CHECK_PAYLOAD: e.payload,
    CHECK_NAME: checkName,
    CHECK_TITLE: checkTitle,
  }


  // For convenience, we'll create three jobs: one for each GitHub Check
  // stage.
  const start = new Job("start-run", checkRunImage)
  start.imageForcePull = true
  start.env = env
  start.env.CHECK_SUMMARY = "Beginning test run"

  const end = new Job("end-run", checkRunImage)
  end.imageForcePull = true
  end.env = env

  // Now we run the jobs in order:
  // - Notify GitHub of start
  // - Run the test
  // - Notify GitHub of completion
  //
  // On error, we catch the error and notify GitHub of a failure.
  start.run().then(() => {
    return job.run()
  }).then( (result) => {
    end.env.CHECK_CONCLUSION = "success"
    end.env.CHECK_SUMMARY = "Build completed"
    end.env.CHECK_TEXT = result.toString()
    return end.run()
  }).catch( (err) => {
    // In this case, we mark the ending failed.
    end.env.CHECK_CONCLUSION = "failure"
    end.env.CHECK_SUMMARY = "Build failed"
    end.env.CHECK_TEXT = `Error: ${ err }`
    return end.run()
  })
}
