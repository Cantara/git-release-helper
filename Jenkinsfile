def vers
def outFile
def release = false
def choices = []

node('master') {
    stage('prepare choices') {
        choices = sh (script: "git tag | git-release-helper | git-release-helper-next", returnStdout: true).trim().split("\n")
    }
}
pipeline {
    agent any
    parameters {
        choice(
            choices: choices,
            description: 'New version',
            name: 'Version'
        )
    }
    tools {
        go 'Go 1.20'
        maven 'Mvn'
    }
    environment {
        NEXUS_CREDS = credentials('Cantara-NEXUS')
    }
    stages {
        stage("pre") {
            steps {
                script {
                    if (env.TAG_NAME) {
                        vers = "${env.TAG_NAME}"
                        release = true
                    } else {
                        vers = sh (script: "git-release-helper-new", returnStdout: true).trim()
                    }
                    artifactId = "git-release-helper"
                    outFile = "${artifactId}-${vers}"
                    echo "New file: ${outFile}"
                }
            }
        }
        stage("test") {
            steps {
                script {
                    testApp()
                }
            }
        }
        stage("build") {
            steps {
                script {
                    echo "V: ${vers}"
                    echo "File: ${outFile}"
                    buildApp(outFile, vers)
                }
            }
        }
        stage("deploy") {
            steps {
                script {
                    echo 'deplying the application...'
                    echo "deploying version ${vers}"
                    if (release) {
                        sh "find . -name '${outFile}-*' -type f -exec curl -v -u "+'$NEXUS_CREDS'+" --upload-file {} https://mvnrepo.cantara.no/content/repositories/releases/no/cantara/gotools/${artifactId}/${vers}/{}  \\;"
                        sh "cd next && find . -name '${outFile}-*' -type f -exec curl -v -u "+'$NEXUS_CREDS'+" --upload-file {} https://mvnrepo.cantara.no/content/repositories/releases/no/cantara/gotools/${artifactId}/${vers}/next/{}  \\;"
                        sh "cd new && find . -name '${outFile}-*' -type f -exec curl -v -u "+'$NEXUS_CREDS'+" --upload-file {} https://mvnrepo.cantara.no/content/repositories/releases/no/cantara/gotools/${artifactId}/${vers}/new/{}  \\;"
                    } else {
                        sh "find . -name '${outFile}-*' -type f -exec curl -v -u "+'$NEXUS_CREDS'+" --upload-file {} https://mvnrepo.cantara.no/content/repositories/snapshots/no/cantara/gotools/${artifactId}/${vers}/{}  \\;"
                        sh "cd next && find . -name '${outFile}-*' -type f -exec curl -v -u "+'$NEXUS_CREDS'+" --upload-file {} https://mvnrepo.cantara.no/content/repositories/snapshots/no/cantara/gotools/${artifactId}/${vers}/next/{}  \\;"
                        sh "cd new && find . -name '${outFile}-*' -type f -exec curl -v -u "+'$NEXUS_CREDS'+" --upload-file {} https://mvnrepo.cantara.no/content/repositories/snapshots/no/cantara/gotools/${artifactId}/${vers}/new/{}  \\;"
                    }
                    sh "rm ${outFile}-*"
                }
            }
        }
    }
}

def testApp() {
    echo 'testing the application...'
    sh './testRecursive.sh'
}

def buildApp(outFile, vers) {
    echo 'building the application...'
    //buildFlags = "-X 'github.com/cantara/gober/webserver/health.Version=${vers}' -X 'github.com/cantara/gober/webserver/health.BuildTime=\$(date)' -X 'github.com/cantara/gober/webserver.Name=${artifactId}' "
    sh "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${outFile}-linux-amd64"
    sh "CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${outFile}-linux-arm64"
    sh "CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${outFile}-darwin-amd64"
    sh "CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${outFile}-darwin-arm64"

    sh "cd next && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${outFile}-linux-amd64"
    sh "cd next && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${outFile}-linux-arm64"
    sh "cd next && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${outFile}-darwin-amd64"
    sh "cd next && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${outFile}-darwin-arm64"

    sh "cd new && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${outFile}-linux-amd64"
    sh "cd new && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${outFile}-linux-arm64"
    sh "cd new && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${outFile}-darwin-amd64"
    sh "cd new && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ${outFile}-darwin-arm64"
}
