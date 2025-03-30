# Starter Project for GO (golang.org)

### Getting Started

Assume the new project name is `watson`. A github repository exists of name `github.com/unbxd/watson`.

> Note: Having hyphen in project name can cause issues with skaffold, as variables in Helm don't support `-`.

- Clone the Repository to some location
- Use `rsync` to copy code
```
cd watson
rsync -av --progress <source_to_go-starter>/go-starter/* . --exclude '.git'
```
- This will not copy `.github` directory, copy that separately
```
cp -R <source_to_go-starter>/go-starter/.github .
```
- Change Configuration in `Makefile`.
```
bin_name := watson.bin
main_pkg := cmd/watson
proj_name := watson
```
- Rename the cmd/app directory to `cmd/watson`
- Rename all prefix of all Environment Variables in cmd/watson/flags.go
```
		&cli.StringFlag{
			Name:        "log.level",
			Value:       "debug",
			Usage:       "set logging level of application. [info, error, warn, debug]",
			DefaultText: "debug",
			EnvVars:     []string{"WATSON_LOG_LEVEL"},
		},
```
- Rename `module` name in `go.mod` file
```
module github.com/unbxd/watson
```
- Find and replace `go-starter` dependency from the project
```
find . -name "*.go"  -type f -exec sed -i 's/go-starter/watson/g' {} \;
```
- Run `make goensure` in the directory, it should run without error
```
make goensure
 > Ensure: go mod tidy
```
- Run `make gobuild` in the direcotry, it should run without error
```
make gobuild
--------------------------------------------
               WATSON MAKE FILE
--------------------------------------------
Bin:                    watson.bin
Module:                 github.com/unbxd/watson
Proj Dir:               /home/uknth/Workspace/fork/watson
Proj Name:              watson
Main Pkg:               cmd/watson
Git Hash:               96ba80dc7cead712
Git Branch:             dev01
Git Stat:               dirty
Git Tag:                latest
Build Date:             1671484804
Build Version:          latest
--------------------------------------------
 > LD FLags
--------------------------------------------
-X github.com/unbxd/watson/cmd/ldflags.Module=github.com/unbxd/watson -X github.com/unbxd/watson/cmd/ldflags.GitHash=96ba80dc7cead712 -X github.com/unbxd/watson/cmd/ldflags.GitBranch=dev01 -X github.com/unbxd/watson/cmd/ldflags.GitStat=dirty -X github.com/unbxd/watson/cmd/ldflags.GitTag=latest -X github.com/unbxd/watson/cmd/ldflags.BuildDate=1671484804 -X github.com/unbxd/watson/cmd/ldflags.Version=latest
--------------------------------------------
 > Delete: rm -f /home/uknth/Workspace/fork/watson/bin/watson.bin 2>/dev/null
 > Ensure: go mod tidy
 > Run Test: go test ./...
 > ---------
?       github.com/unbxd/watson [no test files]
?       github.com/unbxd/watson/cmd/ldflags     [no test files]
?       github.com/unbxd/watson/cmd/watson      [no test files]
?       github.com/unbxd/watson/crud    [no test files]
?       github.com/unbxd/watson/proxy   [no test files]
?       github.com/unbxd/watson/utils/err       [no test files]
 > ---------
 > Build Binary: env GOOS=linux GOARCH=amd64 go build -o watson.bin -ldflags=-X github.com/unbxd/watson/cmd/ldflags.Module=github.com/unbxd/watson -X github.com/unbxd/watson/cmd/ldflags.GitHash=96ba80dc7cead712 -X github.com/unbxd/watson/cmd/ldflags.GitBranch=dev01 -X github.com/unbxd/watson/cmd/ldflags.GitStat=dirty -X github.com/unbxd/watson/cmd/ldflags.GitTag=latest -X github.com/unbxd/watson/cmd/ldflags.BuildDate=1671484804 -X github.com/unbxd/watson/cmd/ldflags.Version=latest github.com/unbxd/watson/cmd/watson
 > Make Dir:  mkdir -p /home/uknth/Workspace/fork/watson/bin
 > Move: mv watson.bin /home/uknth/Workspace/fork/watson/bin
 > --------------------------------------------
 > Go Binary Generated
 > --------------------------------------------
```
- Run `make gobuild docker-build` and ensure that command works successfully
- Create a repository in ECR and ensure `docker login` on the ECR
```
aws ecr describe-repositories --repository-names watson

{
    "repositories": [
        {
            "repositoryArn": "arn:aws:ecr:us-east-1:012629307706:repository/watson",
            "registryId": "012629307706",
            "repositoryName": "watson",
            "repositoryUri": "012629307706.dkr.ecr.us-east-1.amazonaws.com/watson",
            "createdAt": "2022-12-20T02:55:20+05:30",
            "imageTagMutability": "MUTABLE",
            "imageScanningConfiguration": {
                "scanOnPush": false
            },
            "encryptionConfiguration": {
                "encryptionType": "AES256"
            }
        }
    ]
}
```
- Run `make build` and ensure that `watson:latest` docker image is pushed to ECR repository
```
> --------------------------------------------
 > Go Binary Generated
 > --------------------------------------------
 > Build: docker build -t watson:latest /home/uknth/Workspace/fork/watson
[+] Building 2.5s (11/11) FINISHED
 => [internal] load build definition from Dockerfile                                                                                                                                                                                                     0.0s
 => => transferring dockerfile: 38B                                                                                                                                                                                                                      0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                                        0.0s
 => => transferring context: 2B                                                                                                                                                                                                                          0.0s
 => [internal] load metadata for docker.io/library/ubuntu:jammy                                                                                                                                                                                          2.2s
 => [auth] library/ubuntu:pull token for registry-1.docker.io                                                                                                                                                                                            0.0s
 => [1/5] FROM docker.io/library/ubuntu:jammy@sha256:27cb6e6ccef575a4698b66f5de06c7ecd61589132d5a91d098f7f3f9285415a9                                                                                                                                    0.0s
 => [internal] load build context                                                                                                                                                                                                                        0.1s
 => => transferring context: 17.52MB                                                                                                                                                                                                                     0.1s
 => CACHED [2/5] WORKDIR /service                                                                                                                                                                                                                        0.0s
 => CACHED [3/5] RUN mkdir -p /service/bin     && mkdir -p /service/scripts                                                                                                                                                                              0.0s
 => [4/5] ADD bin/* /service/bin/                                                                                                                                                                                                                        0.0s
 => [5/5] ADD scripts/* /service/scripts/                                                                                                                                                                                                                0.0s
 => exporting to image                                                                                                                                                                                                                                   0.1s
 => => exporting layers                                                                                                                                                                                                                                  0.1s
 => => writing image sha256:b2ee10f58b61a727eae7b2bd4083efc15d758fd2a06248325596b7ab1ab51928                                                                                                                                                             0.0s
 => => naming to docker.io/library/watson:latest                                                                                                                                                                                                         0.0s

Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
 > Docker Push: pushing image : [ 012629307706.dkr.ecr.us-east-1.amazonaws.com/watson:latest ]
The push refers to repository [012629307706.dkr.ecr.us-east-1.amazonaws.com/watson]
d2243f3acf15: Pushed
69b1c2075a8e: Pushed
94be1cd1c2e1: Pushed
32277ba5b7a4: Pushed
6515074984c6: Pushed
latest: digest: sha256:d587351db982a33e4cb790797ff72a3e62c8c53dd44afedf0473eae2928b24c5 size: 1360
 > --------------------------------------------
 > Docker Image Generated
 > --------------------------------------------
```
- Rename directory `build/helm/app` to `build/helm/watson`
```
ls -l ./build/helm/
total 8
drwxr-xr-x 3 uknth uknth 4096 Dec 20 02:24 faker
drwxr-xr-x 4 uknth uknth 4096 Dec 20 02:24 watson
```
- Change variables in `helm charts`
  - chart.yaml [ name => watson ]
  - values.yaml [ config.service => watson ]
  - values.yaml [ image.repository => watson ]
  - values.yaml [ nameOverride => watson ]
  - values.yaml [ fullnameOverride => watson ]
  - values.yaml [ ingress.hosts.host => watson.dev.infra ]

- Change variables in skaffold.yaml
  - metadata.name [ app => watson ]
  - artifact.image [ gostarter => watson ]
  - deploy.helm.name [ app => watson ]
  - deploy.helm.chartPath [ build/helm/app => build/helm/watson ]
  - deploy.helm.setValueTemplates.image.repository [ IMAGE_REPO_gostarter => IMAGE_REPO_watson ]. Leave the double curly braces and `.` in prefix as is.
  - deploy.helm.setValueTemplates.image.tag [ IMAGE_TAG_gostarter => IMAGE_TAG_watson ]. Leave the double curly braces and `.` in prefix as is.
- Add `/etc/hosts` entry 
```
127.0.0.1  watson.dev.infra
```
- Run `make skaffold` to run the Skaffold Pipeline, ensure everything is functional and APIs are responding
```
curl watson.dev.infra/monitor
alive
```
- Create ECR repository for `faker-watson` 
```
aws ecr describe-repositories --repository-names faker-watson
```
- Change variables `skaffold.yaml` 
  - build.artifacts.image [ faker => faker-watson ]
  - deploy.helm.releases.setValueTemplates.image [ faker => faker-watson ]

- Add `/etc/hosts` entry for faker
```
127.0.0.1 faker.dev.infra
```
- Run `make skaffold` and ensure all APIs are functioning
```
curl watson.dev.infra/monitor
alive

curl faker.dev.infra/ping
pong

curl watson.dev.infra/json
{
  "hello": "world"
}
```
