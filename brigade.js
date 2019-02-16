const {events, Job, Group} = require("brigadier");
const checkRunImage = "deis/brigade-github-check-run:latest"

const dest = "$GOPATH/src/github.com/stephenhillier/geoprojects";

events.on("check_suite:requested", runTests)
events.on("check_suite:created", runTests)
events.on("check_suite:rerequested", runTests)
events.on("check_run:rerequested", runTests)
events.on("exec", runTests)
events.on("pull_request:opened", pullRequestOpened)

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
  checkRequested(e, p, "Build", "run tests", build)
}

function pullRequestOpened(e, p) {
  var build = new Job("Open PR", "alpine:3.9")
  build.tasks = [
    "PR opened"
  ];
  checkRequested(e, p, "Deploy", "create environment", build)
}

function checkRequested(e, p, name, title, job) {
  console.log("check requested")
  // Common configuration
  const env = {
    CHECK_PAYLOAD: e.payload,
    CHECK_NAME: name,
    CHECK_TITLE: title,
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
