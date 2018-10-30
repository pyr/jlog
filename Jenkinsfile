@Library('jenkins-pipeline') _

repo = "pyr/jlog"

node {
  // Wipe the workspace so we are building completely clean
  cleanWs()

  try {
    dir('src') {
      stage('checkout source code') {
        checkout scm
      }

      build(repo)

      stage('build deb package') {
        parallel(
	  "xenial": { gitPbuilder('xenial', false, '../build-area-xenial') },
	  "bionic": { gitPbuilder('bionic', false, '../build-area-bionic') }
	)
      }

    }
    stage('upload packages') {
      parallel(
        "xenial": { aptlyBranchUpload('xenial','main','build-area-xenial/*.deb') },
	"bionic": { aptlyBranchUpload('bionic','main','build-area-bionic/*.deb') }
	)
      }
  }
  catch (err) {
    currentBuild.result = 'FAILURE'
    throw err
  }

  finally {
    if (currentBuild.result != 'FAILURE') {
    }
  }
}

def build(repo) {
  docker.withRegistry('https://registry.internal.exoscale.ch') {
    def image = docker.image('registry.internal.exoscale.ch/exoscale/golang:1.11')
    image.pull()
    image.inside('-u root --net=host') { sh 'make' }
  }
}
