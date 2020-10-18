# k8s-crawler-operator

----

k8s-crawler-operator是一个部署在k8s集群上operator, 实现了分布式爬虫的自由发布与资源垃圾回收

- 充分利用k8s平台的分布式、水平扩展能力、以及容器的隔离性，实现抓取资源配额限制、进程隔离、垃圾回收等功能
- 方便地发布爬虫定时任务`cronjob`以及一次性任务`job`, 自定义爬虫运行运行列表，无需再编写跨主机发布脚本。注意，同一个名称空间`ns`下的相同名称`name`的运行清单发布时，会回收上一个同名任务的所有资源，是覆盖式的。
- 采用生产者-消费者模型解耦爬虫任务发布和爬虫任务运行逻辑，发布者用redis-tools镜像构建，用来发布消息给共享队列，消费者则包括用来拉取爬虫仓库代码的`initContainer`以及基于`scrapy`的爬虫运行环境的容器。
- 利用configmap资源特性，实现发布任务清单的动态更新
----

## 部署operator

- 手动部署

    - 拉取仓库代码，并进入代码路径

    - 使用`kubectl apply -k`按需依次创建资源，所需资源在config目录下
        - config/crd
        - config/manager
        - config/rbac

- 自动部署

    - 安装[`operator-sdk`](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)开发环境
    - 本项目中差少`bin/manager`，请访问官网后进行补全
    - 按需构建相应的资源，参考手动部署部分，也可以按照官网自行构建operator
    - `deploy`
      ```
       make deploy IMG=quay.io/zhu733756/crawler-operator:v0.0.1
      ```

## 部署样例

- 样例目录在`config/samples`，自定义你的 `run-list`, 并修改资源配置
    ```
    apiVersion: kustomize.config.k8s.io/v1beta1
    kind: Kustomization
    namespace: kube-crawler
    resources:
    - cronjob_v1_croncrawlerjob.yaml
    configMapGenerator:
    - name: crawler-run-list
    files:
    # 暴露给爬虫pod的环境变量配置，也就是运行清单，生产者通过读取该目录来发布任务
    # 实现发布任务动态更新和热加载
    - newspapers-test=run_list/newspapers/test
    - newspapers-main=run_list/newspapers/main
    - website-test=run_list/website/test
    - website-main=run_list/website/main
    generatorOptions:
        # disableNameSuffixHash is true disables the default behavior of adding a
        # suffix to the names of generated resources that is a hash of
        # the resource contents.
        disableNameSuffixHash: true
    # +kubebuilder:scaffold:manifestskustomizesamples
    ```
- 部署样例
    - `kubectl apply -k config/samples`

- 详细版配置以及省略版

    - 详细版，可按需定制
    ``` 详细版
    kind: CronCrawlerJob
    metadata:
        name: crawler-newspapers-test
        namespace: datacrawl
    spec:
        # job, 普通任务
        # cronjob，定时任务
        jobType: "cronjob" 
        schedule: "*/10 * * * *"
        duration: "3s"
        publisher:
            publisherTTL: 600
            # default containers, removes it if you don't want to change it
            containers:
            - name: crawler-publisher
            image: zhu733756/crawler-publisher:v1
            command: ["/bin/sh","/app/prepare.sh"]
            volumeMounts:
            - name:  run-list
              mountPath: /run_list
    consumer:
        consumerTTL: 600
        parallelism: 2
        # default initContainers and containers, removes it if you don't want to change it
        initContainers:
        - name: pull-code
        image: zhu733756/pull-code:v1
        command: ["/bin/sh","/app/pull.sh"]
        volumeMounts:
        - name: code
            mountPath: /app/code
        containers:
        - name: crawler-runner
        image: zhu733756/crawler-runner:v1
        command: ["/bin/bash","/app/getJob.sh"]
        volumeMounts:
        - name: code
            mountPath: /app/code
    # default volumes, removes it if you don't want to change it
    volumes:
        - name: code
        emptyDir: {}
        - name: run-list
        configMap:
            name: crawler-run-list
    env:
        - name: GIT_REPO
            value: "your crawler git repo"
        - name: GIT_BRANCH
            value: "your branch"
        - name: REDIS_URL
            value: "redis://xxx.xx.xx.xx:6379/0"
        - name: PROJECT
            value: "your spider project name"
        - name: CATEGORY
            value: "run list in this project"
    ```
    - 简化版

    ``` 简化版
    kind: CronCrawlerJob
    metadata:
        name: crawler-newspapers-test
        namespace: datacrawl
    spec:
        # job, 普通任务
        # cronjob，定时任务
        jobType: "cronjob" 
        schedule: "*/10 * * * *"
        duration: "3s"
        publisher:
            publisherTTL: 600
        consumer:
            consumerTTL: 600
            parallelism: 2
    env:
    - name: GIT_REPO
        value: "your crawler git repo"
    - name: GIT_BRANCH
        value: "your branch"
    - name: REDIS_URL
        value: "redis://xxx.xx.xx.xx:6379/0"
    - name: PROJECT
        value: "your spider project name"
    - name: CATEGORY
        value: "run list in this project"
    ```
- todo
    - 日志收集
    - 添加爬虫自定义参数




