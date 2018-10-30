@Library('jenkins-pipeline') _

node {
  // Wipe the workspace so we are building completely clean
  cleanWs()

  try {
    dir('src') {
      stage('checkout source code') {
        checkout scm
      }

      Build()

      stage('build deb package') {
        parallel(
	  "xenial": { gitPbuilder('xenial', false '../build-area-xenial') },
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

def Build() {
  def golang = docker.image('golang:latest')
  golang.pull()
  stage('build') {
    golang.inside() {
      sh 'mkdir -p /go/src/jlog'
      sh 'cp Jenkinsfile  Makefile  README.md jlog.go /go/src/jlog'
      sh 'cd /go/src/jlog && make'
      sh 'cp /go/src/jlog/jlog .'
      sh 'make version > .version'
    }
  }
}
