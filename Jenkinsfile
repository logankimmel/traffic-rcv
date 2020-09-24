def containerRepo = 'docker.io/logankimmel/traffic-rcv'
def roxCentral = 'central.rox.binbytes.io'
def gitAddress = 'https://github.com/logankimmel/traffic-rcv.git'

podTemplate(
    namespace: 'cicd', 
    envVars: [
        secretEnvVar(key: 'ROX_API_TOKEN', secretName: 'rox-api', secretKey: 'data')
    ], 
    volumes: [
        secretVolume(secretName: 'regcred', mountPath: '/kaniko/.docker')
    ],
    containers: [
        containerTemplate(name: 'kaniko', image: 'gcr.io/kaniko-project/executor:debug', ttyEnabled: true, command: '/busybox/cat'),
        containerTemplate(name: 'roxctl', image: 'stackrox.io/main:3.0.46.0', ttyEnabled: true, command: '/bin/cat', runAsUser: '0'),
        containerTemplate(name: 'helm', image: 'alpine/helm:3.2.4', ttyEnabled: true, command: 'cat')
    ],
    activeDeadlineSeconds: 300,
    podRetention: onFailure(),
    imagePullSecrets: ['stackrox-io']
) 
{
    node(POD_LABEL) {
        // this will check out a git repo into the working directory for all containers
        def scm = git "${gitAddress}"
        
        stage('Build and Push') {
            container('kaniko') {
                sh "/kaniko/executor -c `pwd` --dockerfile `pwd`/Dockerfile --destination docker.io/${containerRepo}:${scm.GIT_COMMIT}"
            }
        }

        stage('StackRox image check') {
            container('roxctl') {
                sh "/assets/downloads/cli/roxctl-linux image check --image ${containerRepo}:${scm.GIT_COMMIT} -e ${roxCentral}:443"
            }
        }

        stage('StackRox image scan') {
            container('roxctl') {
                sh "/assets/downloads/cli/roxctl-linux image scan --image ${containerRepo}:${scm.GIT_COMMIT} -e ${roxCentral}:443"
            }
        }

        stage('Render Helm') {
            container('helm') {
                sh "helm template traffic-rcv `pwd`/traffic-rcv --set=image.repository=docker.io/${containerRepo}:${scm.GIT_COMMIT} > `pwd`/chart.yaml"
            }
        }

         stage('StackRox deployment check') {
            container('roxctl') {
                sh "/assets/downloads/cli/roxctl-linux deployment check --file `pwd`/chart.yaml -e ${roxCentral}:443"
            }
        }        
    }
}