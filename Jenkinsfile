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
        gitPbuilder('xenial')
      }

      stage('upload packages') {
                sh 'cp ../build-area/*.deb .'
		aptlyUpload('staging')
      }
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
      sh 'mkdir -p /go/src/jlog
      sh 'cp Jenkinsfile  Makefile  README.md jlog.go /go/src/jlog'
      sh 'cd /go/src/jlog && make'
      sh 'cp /go/src/jlog/jlog .'
      sh 'make version > .version'
    }
  }
}
