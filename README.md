部署服务
----
## Knative Service

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: deploy # 服务名称
spec:
  template:
    spec:
      containers:
        - image: dongfg-docker.pkg.coding.net/serverless-func/docker/func-deploy:23.11.4 # 服务镜像
          ports:
            - containerPort: 9000 # 运行端口
          env:
            - name: GIT_REPO
              value: "" # IaC 仓库地址(只支持 http 目前), 仓库结构见下方
            - name: GIT_USER
              value: "" # Iac 仓库用户
            - name: GIT_EMAIL
              value: "" # Iac 仓库提交邮箱
            - name: GIT_PASSWORD
              value: "" # Iac 仓库密码
            - name: KUBE_CONFIG
              value: "" # kubernetes 集群配置
            - name: KUBE_NAMESPACE
              value: "" # knative 服务的 namespace, 默认 func
```

## IaC Repo Structure

```shell
$ tree -L 2 
.
├── README.md
└── func
    ├── func-xxx.yaml
    ├── func-xxx.yaml
    ├── func-xxx.yaml
```