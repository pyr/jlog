@Library('jenkins-pipeline') _

repo = "pyr/jlog"

node {
  // Wipe the workspace so we are building completely clean
  cleanWs()

  try {
    dir('src') {
      stage('checkout') { checkout scm }
      stage('build')    { build(repo) }
      stage('package')  {
        parallel (
          "bionic": {
              gitPbuilder('bionic', false, '../build-area-bionic')
          },
          "focal": {
              gitPbuilder('focal', false, '../build-area-focal')
          }
        )
      }
      stage('upload')   {
        aptlyUpload('staging', 'bionic', 'main', '../build-area-bionic/*deb')
        aptlyUpload('staging', 'focal', 'main', '../build-area-focal/*deb')
      }
    }
  }
  catch (err) {
    currentBuild.result = 'FAILURE'
    throw err
  }

  finally {
    if (currentBuild.result != 'FAILURE') {
      currentBuild.result = 'SUCCESS'
    }
    updateGithubCommitStatus(currentBuild.result, "${env.WORKSPACE}/src")
  }
}

def build(repo) {
  docker.withRegistry('https://registry.internal.exoscale.ch') {
    def image = docker.image('registry.internal.exoscale.ch/exoscale/golang:1.11')
    image.pull()
    image.inside('-u root --net=host') { sh 'make' }
  }
}
