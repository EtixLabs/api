node {
  env.PROD_SERVER = 'etix@10.0.5.57'
  env.APP_NAME = 'pepersalt-api'

   stage('Pull SCM') {
       git credentialsId: '9221578f-38d8-4978-a5f6-823dad6dbd23', url: 'ssh://git@git.etixlabs.com/pep/api.git'
       sh 'git rev-parse HEAD > commit'
       env.COMMIT_ID = readFile('commit').trim()
       println 'The commit number is ${COMMIT_ID}'
   }

   stage('Build') {
      sh 'docker-compose run --rm build -t registry.etixlabs.com/${APP_NAME}:${COMMIT_ID} .'
      sh 'docker build -t registry.etixlabs.com/${APP_NAME}:${COMMIT_ID} .'
   }

//   stage('Test') {
//       sh "docker run -v `pwd`/:/jenkins/ --entrypoint 'bash' registry.etixlabs.com/${APP_NAME}:${COMMIT_ID} -c \"bundle exec rspec -f h -o /jenkins/index.html\""
//   }

   stage('Push') {
       sh 'docker tag registry.etixlabs.com/${APP_NAME}:${COMMIT_ID} registry.etixlabs.com/${APP_NAME}'
       sh 'docker push registry.etixlabs.com/${APP_NAME}:${COMMIT_ID}'
       sh 'docker push registry.etixlabs.com/${APP_NAME}'
   }

   stage('Deploy') {
       input 'Do you want to deploy?'
       sh 'ssh ${PROD_SERVER} docker-compose pull ${APP_NAME}'
       sh 'ssh ${PROD_SERVER} docker-compose up -d ${APP_NAME}'
   }
}
