def COLOR_MAP = [
    'SUCCESS': 'good',
    'FAILURE': 'danger',
]
pipeline {
    agent none
    environment {
         DOCKER_CREDENTIALS = credentials('docker-builder')
         BUILD_USER         = 'Jenkins'
         PROJECT            = 'simple-bank'
         VERSION            = 'latest'
    }
    stages {
        stage("Docker build") {
            agent none
            steps {
              sh "docker build --network=host --tag docker-host.banvien.com.vn/${PROJECT}:${VERSION} ."
            }
        }
        stage("Docker Push") {
            agent none
            steps {
                sh "docker push docker-host.banvien.com.vn/${PROJECT}/${ENVIRONMENT}/${SERVICE}:${VERSION}"
            }
        }
        stage("Deploy") {
            agent none
            steps {
                sh "kubectl delete -f kubernetes-${ENVIRONMENT}.yaml | echo IGNORE"
                sh "kubectl apply -f kubernetes-${ENVIRONMENT}.yaml"
            }
        }
    }
    post {
        always {
            node('master'){
                slackSend channel: '#server-dev',
                    color: COLOR_MAP[currentBuild.currentResult],
                    message: "*${currentBuild.currentResult}:* Job ${env.JOB_NAME} build ${env.BUILD_NUMBER} by ${BUILD_USER}\n More info at: ${env.BUILD_URL}"
            }
        }
    }
}
