pipeline {
    agent any
    environment {
        BBSCOUT_SSH = credentials('bbscout_key')
        BBSCOUT_PASS = credentials('bbscout_password')
    }
    stages {
        stage('Install Dependencies') {
            steps {
                echo "Installing required dependencies..."
                sh "apt-get update && apt-get install -y sshpass"
            }
        }
        stage('Deploy') {
            
            steps {
                echo "Deploying..."

                script {
                    // Parse the BBSCOUT_SSH into components
                    def parts = BBSCOUT_SSH.split('@|:')
                    def user = parts[0]
                    def host = parts[1]
                    def port = parts[2]

                    sh """
                        sshpass -p "$BBSCOUT_PASS" ssh -o StrictHostKeyChecking=no -p $port $user@$host '
                            cd bbscout/backend &&
                            git pull origin main &&
                            echo "$BBSCOUT_PASS" | sudo -S docker compose up --build -d
                        '
                    """
                }
            }
        }
    }

    post {
        success {
            echo "Deployment successful!"
        }
        failure {
            echo "Deployment failed!"
        }
    }
}
