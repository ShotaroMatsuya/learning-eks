# GitOps Demo

## Jenkins installation

### First, Creating a key pair

To create your key pair:

1. Open the Amazon EC2 console at https://console.aws.amazon.com/ec2/ and sign in.

2. In the navigation pane, under NETWORK & SECURITY, select Key Pairs.

3. Select Create key pair.

4. For Name, enter a descriptive name for the key pair. Amazon EC2 associates the public key with the name that you specify as the key name. A key name can include up to 255 ASCII characters. It cannot include leading or trailing spaces.

5. For File format, select the format in which to save the private key.

- For OpenSSH compatibility, select pem.

- For PuTTY compatibility, select ppk.

6. Select Create key pair.

7. The private key file downloads automatically. The base file name is the name you specified as the name of your key pair, and the file name extension is determined by the file format you chose. Save the private key file in a safe place.

8. If you use an SSH client on a macOS or Linux computer to connect to your Linux instance, run the following command to set the permissions of your private key file so that only you can read it.

```bash
$ chmod 400 <key_pair_name>.pem
```

---

### Second, Creating a security group

To create and configure your security group:

1. Decide who may access your instance. For example, a single computer or all trusted computers on a network. For this tutorial, you can use the public IP address of your computer.

2. Sign in to the AWS Management Console.

3. Open the Amazon EC2 console by selecting EC2 under Compute.

4. In the left-hand navigation bar, select Security Groups, and then select Create Security Group.

5. In Security group name, enter `WebServerSG` or any preferred name of your choice, and provide a description.

6. Select your VPC from the list. You can use the default VPC.

7. On the Inbound tab, add the rules as follows:

- Select Add Rule, and then select SSH from the Type list.

- Under Source, select Custom, and in the text box, enter the IP address from step 1.

- Select Add Rule, and then select HTTP from the Type list.

- Select Add Rule, and then select Custom TCP Rule from the Type list.

- Under Port Range, enter 8080.

8. Select Create.

---

## Third, Launching an Amazon EC2 instance

To launch an EC2 instance:

1. Sign in to the the AWS Management Console.

2. Open the Amazon EC2 console by selecting EC2 under Compute.

3. From the Amazon EC2 dashboard, select Launch Instance.

4. The Choose an Amazon Machine Image (AMI) page displays a list of basic configurations called Amazon Machine Images (AMIs) that serve as templates for your instance. Select the HVM edition of the Amazon Linux AMI.

5. On the Choose an Instance Type page, the t3.small instance is selected. Once chosen, you can select Review and Launch.

6. On the Review Instance Launch page, select Edit security groups.

7. On the Configure Security Group page:

- Select Select an existing security group.

- Select the WebServerSG security group that you created.

- Select Review and Launch.

8. On the Review Instance Launch page, select Launch.

9. In the Select an existing key pair or create a new key pair dialog box, select Choose an existing key pair. Then select the key pair you created in the creating a key pair section above or any existing key pair you intend to use.

10. In the left-hand navigation bar, choose Instances to view the status of your instance. Initially, the status of your instance is pending. After the status changes to running, your instance is ready for use.

---

## Next, Installing and configuring Jenkins

In this step you will deploy Jenkins on your EC2 instance by completing the following tasks:

### Connecting to your Linux instance

1. Before you connect to your instance, get the public DNS name of the instance using the Amazon EC2 console.

- Select the instance and locate Public DNS.

2. Using SSH to connect to your instance

- Use the ssh command to connect to the instance. You will specify the private key (.pem) file and ec2-user@public_dns_name.

```bash
$ ssh -i /path/my-key-pair.pem ec2-user@ec2-198-51-
100-1.compute-1.amazonaws.com
```

- You will receive a response like the following:

```bash
The authenticity of host 'ec2-198-51-100-1.compute1.amazonaws.com (10.254.142.33)' cant be
established.

RSA key fingerprint is 1f:51:ae:28:bf:89:e9:d8:1f:25:5d:37:2d:7d:b8:ca:9f:f5:f1:6f.

Are you sure you want to continue connecting
(yes/no)?
```

- Enter yes.

You will receive a response like the following:

```bash
Warning: Permanently added 'ec2-198-51-100-1.compute1.amazonaws.com' (RSA) to the list of known hosts.
```

### Downloading and installing Jenkins

1. Ensure that your software packages are up to date on your instance by uing the following command to perform a quick software update:

```bash
[ec2-user ~]$ sudo yum update –y
```

2. Add the Jenkins repo using the following command:

```bash
[ec2-user ~]$ sudo wget -O /etc/yum.repos.d/jenkins.repo \
    https://pkg.jenkins.io/redhat-stable/jenkins.repo
```

3. Import a key file from Jenkins-CI to enable installation from the package:

```bash
[ec2-user ~]$ sudo rpm --import https://pkg.jenkins.io/redhat-stable/jenkins.io.key

[ec2-user ~]$ sudo yum upgrade
```

4. Install Java & Jenkins:

```bash
[ec2-user ~]$ sudo amazon-linux-extras install java-openjdk11 -y

[ec2-user ~]$ sudo yum install jenkins -y
```

5. Enable the Jenkins service to start at boot:

```bash
[ec2-user ~]$ sudo systemctl enable jenkins
```

6. Start Jenkins as a service:

```bash
[ec2-user ~]$ sudo systemctl start jenkins
```

You can check the status of the Jenkins service using the command:

```bash
[ec2-user ~]$ sudo systemctl status jenkins
```

### Configuring Jenkins

1. Connect to http://<your_server_public_DNS>:8080 from your browser. You will be able to access Jenkins through its management interface:

2. As prompted, enter the password found in /var/lib/jenkins/secrets/initialAdminPassword.

- Use the following command to display this password:

```bash
[ec2-user ~]$ sudo cat /var/lib/jenkins/secrets/initialAdminPassword
```

3. The Jenkins installation script directs you to the Customize Jenkins page. Click Install suggested plugins.

4. Once the installation is complete, the Create First Admin User will open. Enter your information, and then select Save and Continue.

5. On the left-hand side, select Manage Jenkins, and then select Manage Plugins.

6. Select the Available tab, and then enter Amazon EC2 plugin at the top right.

7. Select the checkbox next to Amazon EC2 plugin, and then select Install without restart.

8. Once the installation is done, select Back to Dashboard.

9. Select Configure a cloud if there are no existing nodes or clouds.

10. If you already have other nodes or clouds set up, select Manage Jenkins.

- After navigating to Manage Jenkins, select Configure Nodes and Clouds from the left hand side of the page.
- From here, select Configure Clouds.

11. Select Add a new cloud, and select Amazon EC2. A collection of new fields appears.

12. Fill out all the fields. You will need to input the Amazon EC2 Credentials, ensuring they are AWS Credentials.

---

Jenkins is installed on EC2. Follow the instructions on https://www.jenkins.io/doc/tutorials/tutorial-for-installing-jenkins-on-AWS/ .

You can skip "Configure a Cloud" part for this demo. Please note some commands from this link might give errors, below are the workarounds:

---

### Jenkins plugins

Install the following plugins for the demo.

- Amazon EC2 plugin (No need to set up Configure Cloud after)
- Docker plugin
- Docker Pipeline
- GitHub Integration Plugin
- Parameterized trigger Plugin

## ArgoCD installation

Install ArgoCD in your Kubernetes cluster following this link - https://argo-cd.readthedocs.io/en/stable/getting_started/

## How to run!

Follow along with my Udemy Kubernetes course lectures (GitOps Chapter) to understand how it works, detailed setup instructions, with step by step demo. My highest rated Kubernetes EKS discounted Udemy course link in www.cloudwithraj.com

```

```
