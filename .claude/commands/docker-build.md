---
description: 构建 Docker 镜像并支持自定义 tag
argument-hint: [image_tag]
model: sonnet
allowed-tools: Bash(docker:*), Bash(make:docker-build)
---

你现在是`issue2md`项目的 DevOps 工程师，你的任务是构建 Docker 镜像。

## 任务目标

使用 Makefile 的 `docker-build` 目标构建 issue2md 项目的 Docker 镜像。支持自定义镜像 tag，如果未提供则默认使用 `latest`。

## 参数说明

- **`$1` (image_tag)**: 可选参数，指定 Docker 镜像的 tag
  - 如果用户提供了参数（如 `v1.0.0`、`test`、`dev`），使用该值作为 tag
  - 如果用户未提供参数或为空，默认使用 `latest`

## 执行步骤

### 步骤 1: 确定镜像 tag

首先判断用户是否提供了镜像 tag 参数：

```
如果 $1 为空或未提供：
  TAG=latest
否则：
  TAG=$1
```

### 步骤 2: 执行构建命令

使用以下命令构建 Docker 镜像：

```bash
# 方式 1: 直接使用 make（推荐）
make docker-build

# 方式 2: 使用 docker build
docker build -t issue2md:$TAG .
```

**注意**: Makefile 的 `docker-build` 目标会自动：
- 注入 VERSION、BUILD_DATE、GIT_COMMIT 等构建参数
- 同时打两个 tag：`issue2md:$VERSION` 和 `issue2md:latest`

### 步骤 3: 验证构建结果

**如果构建成功（退出码为 0）：**
- 显示构建成功的镜像信息
- 显示镜像大小、ID、创建时间
- 可选：运行 `docker run --rm issue2md:$TAG --version` 验证镜像
- 可选：显示镜像分层信息

**如果构建失败（退出码非 0）：**
进入步骤 4 进行错误分析

### 步骤 4: 分析构建失败原因

当 Docker 构建失败时，系统地分析错误：

#### 4.1 Dockerfile 错误分析

- **语法错误**: Dockerfile 指令语法是否正确
- **路径错误**: COPY、ADD 等指令的路径是否存在
- **基础镜像**: FROM 指定的镜像是否可访问
- **权限问题**: 文件复制和执行权限

#### 4.2 构建上下文问题

- **上下文过大**: 是否因不必要的文件导致构建缓慢
- **.dockerignore**: 是否正确配置以排除不需要的文件
- **文件位置**: 源文件是否在正确的位置

#### 4.3 依赖和下载问题

- **网络问题**: 基础镜像或依赖下载失败
- **Go modules**: go mod download 是否成功
- **DNS/代理**: 是否存在网络连接问题

#### 4.4 多阶段构建问题

- **阶段依赖**: builder 阶段是否成功完成
- **文件复制**: `COPY --from=builder` 是否正确
- **架构匹配**: builder 和 final 阶段的架构是否一致

#### 4.5 资源问题

- **磁盘空间**: Docker 守护进程存储空间是否充足
- **内存**: 构建过程中内存是否不足
- **Docker 守护进程**: Docker 服务是否正常运行

### 步骤 5: 提供解决方案

根据错误分析，提供具体的解决方案：

1. **Dockerfile 问题**: 提供修复 Dockerfile 的具体建议
2. **依赖问题**: 提供网络或依赖相关的修复命令
3. **环境问题**: 提供 Docker 环境配置或安装的步骤
4. **资源问题**: 提供清理或优化的建议

## 输出格式

### 构建成功时

```
✅ Docker 镜像构建成功！

镜像信息：
- Repository: issue2md
- Tag: $TAG
- Image ID: [镜像ID]
- Size: [镜像大小]
- Created: [创建时间]

镜像详情：
[运行 docker images issue2md:$TAG 的输出]

验证：
[运行 docker run --rm issue2md:$TAG --version 的输出]
```

### 构建失败时

```
❌ Docker 镜像构建失败！

## 错误分析

### 错误阶段
[识别失败的构建阶段：builder/final]

### 错误类型
[识别的错误类型：语法/路径/网络/依赖等]

### 错误详情
[完整的错误输出]

### 根本原因
[分析错误的根本原因]

### 解决方案

#### 方案 1: [推荐方案]
[具体的解决步骤]

#### 方案 2: [备选方案]
[替代的解决方法]

### 优化建议
[如何改进 Dockerfile 或构建流程]
```

## 常见问题和解决方案

### 问题 1: 基础镜像下载失败

**错误示例**:
```
ERROR: golang:1.24.9-alpine: error pulling image
```

**解决方案**:
```bash
# 检查网络连接
ping -c 3 registry-1.docker.io

# 配置 Docker 镜像加速器（中国大陆）
sudo systemctl edit docker
# 添加镜像加速配置

# 手动拉取基础镜像
docker pull golang:1.24.9-alpine
```

### 问题 2: go mod download 失败

**错误示例**:
```
#11 ERROR: go mod download: context deadline exceeded
```

**解决方案**:
```bash
# 使用 Go proxy
docker build --build-arg GOPROXY=https://goproxy.cn,direct -t issue2md:$TAG .

# 或者预下载依赖
go mod download
docker build -t issue2md:$TAG .
```

### 问题 3: 文件复制失败

**错误示例**:
```
COPY failed: file not found
```

**解决方案**:
```bash
# 检查 .dockerignore 配置
cat .dockerignore

# 确认文件存在
ls -la cmd/issue2md/main.go

# 检查构建上下文
docker build --no-cache -t issue2md:$TAG .
```

### 问题 4: 磁盘空间不足

**错误示例**:
```
no space left on device
```

**解决方案**:
```bash
# 清理 Docker 未使用的镜像
docker system prune -a

# 查看 Docker 磁盘使用情况
docker system df

# 只清理悬空镜像
docker image prune
```

## 使用示例

### 示例 1: 默认 tag (latest)

```
用户: /docker-build

执行: make docker-build
结果: 构建镜像 issue2md:latest
```

### 示例 2: 自定义 tag

```
用户: /docker-build v1.0.0

执行: make docker-build
结果: 构建镜像 issue2md:v1.0.0
```

### 示例 3: 开发版本

```
用户: /docker-build dev

执行: make docker-build
结果: 构建镜像 issue2md:dev
```

## 构建优化建议

### 1. 利用 Docker 层缓存

- **优先复制依赖**: 先复制 `go.mod` 和 `go.sum`
- **后复制源码**: 最后复制源代码
- **避免无效变更**: 将不常变化的文件放在前面

### 2. 多阶段构建优势

- **减小镜像体积**: 只复制编译后的二进制文件
- **提高安全性**: 最终镜像不包含源代码和构建工具
- **加快分发**: 更小的镜像传输更快

### 3. 构建参数调优

```bash
# 不使用缓存构建
docker build --no-cache -t issue2md:$TAG .

# 显示构建详情
docker build --progress=plain -t issue2md:$TAG .

# 并行构建
docker build --build-arg BUILDKIT_INLINE_CACHE=1 -t issue2md:$TAG .
```

## 相关文件

在构建和分析过程中，你可能需要查看以下文件：
- `Dockerfile` - Docker 镜像构建配置
- `.dockerignore` - Docker 构建排除文件
- `Makefile` - docker-build 目标定义
- `go.mod` - Go 依赖定义
- `cmd/issue2md/main.go` - 主入口文件

## 验证和测试

构建成功后，建议执行以下验证：

```bash
# 1. 验证镜像存在
docker images | grep issue2md

# 2. 验证版本信息
docker run --rm issue2md:$TAG --version

# 3. 验证帮助信息
docker run --rm issue2md:$TAG --help

# 4. 测试实际功能
docker run --rm issue2md:$TAG https://github.com/golang/go/issues/1

# 5. 检查镜像安全（可选）
docker scan issue2md:$TAG
```

## 注意事项

1. **tag 规范**: 使用语义化版本号（如 v1.0.0）或描述性标签（如 dev、test）
2. **大小监控**: 构建完成后关注镜像大小，正常应在 15-20MB 之间
3. **安全扫描**: 定期使用 `docker scan` 检查镜像安全漏洞
4. **多架构支持**: 如需支持 ARM64，可使用 `docker buildx`
5. **清理**: 定期清理旧的测试镜像以节省磁盘空间
