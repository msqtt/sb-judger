<div align="center">
<h1 align="center">
<br>SB-JUDGER 🤪</h1>
<h3>◦ ► 使用 Go 语言开发的轻量 OnlineJudge Server.</h3>
<img src="https://img.shields.io/github/license/msqtt/sb-judger?style=flat-square&color=5D6D7E" alt="GitHub license" />
<img src="https://img.shields.io/github/last-commit/msqtt/sb-judger?style=flat-square&color=5D6D7E" alt="git-last-commit" />
<img src="https://img.shields.io/github/commit-activity/m/msqtt/sb-judger?style=flat-square&color=5D6D7E" alt="GitHub commit activity" />
<img src="https://img.shields.io/github/languages/top/msqtt/sb-judger?style=flat-square&color=5D6D7E" alt="GitHub top language" />
</div>

---

## 🖼️ 截图

![image](https://github.com/msqtt/sb-judger/assets/94043894/685b8195-985e-4a01-9b44-436e66b3cdbe)

## 📝 简介

sandbox-judger 使用 LXC 技术，为用户提交每一个的代码进程创建一个沙箱，并在指定的系统资源下运行程序，达到安全判题的效果。

它使用 namespace、cgroup 来隔离资源，使用 overlayfs 作为 Unionfs。

> **目前 sb-judger 还不支持 cgroupv1，在使用前请确保运行环境使用的是 cgroupv2。**

## 🚀 开始

### 🔧 安装

1. 克隆仓库:
```sh
git clone https://github.com/msqtt/sb-judger
```

2. 切换目录:
```sh
cd sb-judger
```

3. 制作 rootfs:

```sh
make rootfs
```
> 如果你有自己的 rootfs 请忽略这步
> `mkdir rootfs` 后，直接把根目录解压到 `rootfs` 即可。

4. 开始构建:
```sh
make build
```

### 🤖 启动

```sh
./sb-judger
```

打开运行代码测试页面 :

```sh
open http://localhost:8080
```

### 🐬 Docker 

#### 用用我的

```sh
docker pull msqt/sb-judger:latest
docker run --privileged -d -p8080:8080 -p9090:9090 msqt/sb-judger
```

#### 自己构建

```sh
make docker
```

---

## 🌐 API

sb-judger 使用 `grpc` 和 `http` 作为通讯协议，且使用 `grpc-http-gateway` 提供 `http` 服务。

- [Http OpenAPI](https://github.com/msqtt/sb-judger/blob/master/api/openapi/v1/judger/judger_service.swagger.json)
  - [Apifox doc](https://4725lf5hpc.apifox.cn)
- [Grpc Protos](https://github.com/msqtt/sb-judger/tree/master/api/protos/v1)

## ⚙️ 配置

- [语言配置](https://github.com/msqtt/sb-judger/blob/master/configs/lang.json)
- [软件配置](https://github.com/msqtt/sb-judger/blob/master/configs/app.env)
  - 可以直接通过环境变量传入，比如: `HTTP_ADDR=0.0.0.0:8080 ./sb-judger`，docker 也是同理。



## 🛣 路线

> - [X] `ℹ️ Support cgroupv2`
> - [ ] `ℹ️ Test`
> - [ ] `ℹ️ Support cgroupv1`
> - [ ] `ℹ️ ...`

## 🧮 支持语言

> - [X] `ℹ️ c/cpp`
> - [X] `ℹ️ golang`
> - [X] `ℹ️ python`
> - [X] `ℹ️ java`
> - [X] `ℹ️ rust`
> - [ ] `ℹ️ ...`
---

## 🤝 Contributing

Contributions are welcome! Here are several ways you can contribute:

- **[Submit Pull Requests](https://github.com/msqtt/sb-judger/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.
- **[Join the Discussions](https://github.com/msqtt/sb-judger/discussions)**: Share your insights, provide feedback, or ask questions.
- **[Report Issues](https://github.com/msqtt/sb-judger/issues)**: Submit bugs found or log feature requests for MSQTT.

#### *Contributing Guidelines*

<details closed>
<summary>Click to expand</summary>

1. **Fork the Repository**: Start by forking the project repository to your GitHub account.
2. **Clone Locally**: Clone the forked repository to your local machine using a Git client.
   ```sh
   git clone <your-forked-repo-url>
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear and concise message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to GitHub**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.

Once your PR is reviewed and approved, it will be merged into the main branch.

</details>

---

## 📄 License


This project is protected under the [MPL2](https://choosealicense.com/licenses/mpl-2.0/) License. For more details, refer to the [LICENSE](./LICENSE) file.

---


