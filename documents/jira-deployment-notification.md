# Jira Deployment Notification via Jenkins

**Caveat - currently this feature is only available for [Entigo](https://www.entigo.com/) customers. Feel free to [write us](mailto:info@entigo.com)  and ask for support.**

This document describes how to use [entigo-k8s-gitops](../readme.md) ```update``` command with notification feature that notifies [Jira](https://www.atlassian.com/software/jira) about specified deployments.

## Setup

Instructions how to set up the integration.

### Plugins setup

1. In Jira go to **Apps** → **Explore more app** and install **[Jenkins for Jira (Official)](https://marketplace.atlassian.com/apps/1227791/jenkins-for-jira-official?hosting=cloud&tab=overview)** plugin.  
    * Then go to **Apps** → **Manage your apps** and on the left sidebar click on **Jenkins for Jira**. 
    * Then click on **Connect a Jenkins server**. Follow the instructions to create a webhook.

2. Install **Atlassian Jira Software Cloud** 
    * Open your Jenkins server
    * Navigate to Manage **Jenkins** -> **Manage plugins**
    * In the **Available** tab, search for “[**Atlassian Jira Software Cloud**](https://plugins.jenkins.io/atlassian-jira-software-cloud/)”
    * Check the "**Install**" checkbox
    * Click "**Download now and install after restart**"

### Jenksinsfile setup

To enable notification feature just add notification spcific flags (```--notify*```) and you are good to go. 

Example [Jenkinsfile](https://www.jenkins.io/doc/book/pipeline/jenkinsfile/) how to use ```entigo-k8s-gitops update``` command with notification:

```groovy
JiraDeploymentInfo jiraDeploymentInfo = new JiraDeploymentInfo('<environmentId>', '<environmentName>', 'environmentType', [])

pipeline {
    agent {
        docker { 
            image 'entigolabs/entigo-k8s-gitops:<tag>' 
            }
    }

    stages {
        stage('Step 1') {
            steps {
                withCredentials(bindings: [sshUserPrivateKey(credentialsId: '<repositoryCredentials>', keyFileVariable: 'SSH_KEY_FOR_GIT')]) {
                    script {
                        def updateWithNotifyCmd = "gitops update " +
                                                        "--git-repo=<repoAddress> " +
                                                        "--git-branch=master --git-key-file=\"$SSH_KEY_FOR_GIT\" " +
                                                        "--app-path=<appPath> " +
                                                        "--images=<imagesToModify> " +
                                                        "--keep-registry=<boolen> " +
                                                        "--notify-env=<notifyEnvName> " +
                                                        "--notify-registry-uri=<registryUri> " +
                                                        "--notify-auth-token=<tokenKey=tokenValue> " +
                                                        "--notify-api-url=<baseUrl>/api/cicd/v1/atlassian/jira/deployments/info"
                        def stdout = executeAndGetStdout(updateWithNotifyCmd) 
                        jiraDeploymentInfo.issueKeys = findJiraIssueKeys(stdout)
                        jiraSendDeploymentInfoCustom(jiraDeploymentInfo) 
                    }
                }
            }
        }
    }
    post {
        success{
            echo 'Step 1 successful'
            jiraSendDeploymentInfoCustom(jiraDeploymentInfo, 'successful')
        } 
        unsuccessful{
            echo 'Step 1 unsuccessful'
            jiraSendDeploymentInfoCustom(jiraDeploymentInfo, 'failed')
        } 
    }
}

class JiraDeploymentInfo {
    String environmentId;
    String environmentName;
    String environmentType;
    List<String> issueKeys;

    public JiraDeploymentInfo(String environmentId, String environmentName, String environmentType, List<String> issueKeys) {
        this.environmentId = environmentId;
        this.environmentName = environmentName;
        this.environmentType = environmentType;
        this.issueKeys = issueKeys;
    }
}

def executeAndGetStdout(String shellScript) {
    try {
        def outfile = 'stdout.out'
        def status = sh(script:"${shellScript} >${outfile} 2>&1", returnStatus:true)
        def output = readFile(outfile).trim()
        if (status != 0) {
            error output
        }
        println(output)
        return output
    } catch (Exception ex) {
        println("Unable to execute script: ${ex}")
    }
}

def findJiraIssueKeys(String updateWithNotifyStdout) {
    def stdoutLastLine = updateWithNotifyStdout.tokenize().last()
    def isSuccessfulResponse = stdoutLastLine.contains('data') && stdoutLastLine.contains('deployedJiraIssueKeys')
    if(!isSuccessfulResponse) {
        error 'gitops update with notify was unsuccesful - notification insertion failed; inspect gitops update command stdout '
    }
    def jiraIssueKeysPattern = ~"[A-Z]{2,}-\\d+"
    def matcher = stdoutLastLine =~ jiraIssueKeysPattern
    return matcher.findAll()
}


def jiraSendDeploymentInfoCustom(JiraDeploymentInfo deploymentInfo, String state = null) {
    if(state == null) {
        jiraSendDeploymentInfo environmentId: deploymentInfo.environmentId, environmentName: deploymentInfo.environmentName, environmentType: deploymentInfo.environmentType, issueKeys: deploymentInfo.issueKeys
    } else {
        jiraSendDeploymentInfo environmentId: deploymentInfo.environmentId, environmentName: deploymentInfo.environmentName, environmentType: deploymentInfo.environmentType, issueKeys: deploymentInfo.issueKeys, state: state
    }
}
```
#### ```Jenkinsfile``` notes
* In this example credentialsId ```<repositoryCredentials>``` means that in **Jeknkins Dashboard** -> **Manage Jenksins** -> **Manage Credentials** you have to add **+ Add Credentials** with kind of ```SSH Username with private key``` where **username** is ```git```, **username as secret** checkbox is ```checked``` and **private key is entered directly**.
* The ```environmentType``` must be one of the following: ```unmapped```, ```development```, ```testing```, ```staging```, ```production```.