<p align="center">
  <img width="180px" height="180px" src=".github/logo.png"/>
</p>

<p align="center">一款基于Golang开发、集成Eino框架实现原生AI编排的高性能即时通讯系统服务端</p>
<br>

## 🚀 项目介绍
Linyu 2.0 是基于 Golang 高性能即时通讯 (IM) 系统。它不仅提供稳定的单聊、群聊及音视频通话能力，更通过集成开源的 Eino 框架，实现了原生级 AI 编排能力。系统采用分层架构设计，具备高并发处理能力与极强的横向扩展性，旨在为开发者提供一个开箱即用的工业级通讯底座。

## 🌟 核心优势
- 原生 AI 驱动: 基于 Eino 框架实现模型编排，内置 AI 机器人、意图识别及长短期记忆管理。

- 卓越通信性能: 基于 WebSocket 与高效的消息 ACK 机制，保障消息的不丢、不重、有序。

- 全栈音视频: 集成 LiveKit，支持低延迟、高可靠的多人视频会议与语音通话。

- 模块化设计: 采用多模块开发来进行解耦。

## 🛠️ 技术栈
- 核心语言: Golang (高并发编程模型)

- 消息引擎: RocketMQ (保障分布式消息的一致性与削峰填谷)

- 存储矩阵:
  - MySQL: 核心业务数据的关系型存储
  - Redis: 会话缓存及热点数据加速
  - Weaviate: 向量数据库，支撑 AI 语义搜索与长期记忆

- 实时通信: WebSocket (linyu-im) + LiveKit (音视频)

- AI 编排: Eino

## 🏗️ 模块导览
### linyu-gateway - API 网关
- 统一入口: 实现路由转发。
- 安全防护: 鉴权拦截、频率限制及日志审计。

### linyu-im - 通讯核心
- 长连接管理: 支撑海量 WebSocket 连接。
- 消息路由: 精准的消息下发策略。

### linyu-ai - 智能大脑
- 模型编排: 基于 Eino 灵活构建 AI 任务流。
- AI机器人: 快速接入大模型，实现智能问答与自动化交互。
- 向量增强: 对接 Weaviate 实现长期记忆等。

### linyu-voip-chat - 音视频服务
- 多人会议: 基于 LiveKit 的低延迟视频通话方案。
- 信令管理: 完整的通话建立、维持与挂断逻辑。

### linyu-auth - 认证中心
- 多模式登录: 支持标准账号、Gitee OAuth、LDAP 等多种身份验证方式。

### linyu-basic-service - 业务基础
- 社交核心: 涵盖单聊、群聊，通讯录、群组管理、过往（朋友圈）等社交功能。

### linyu-cloud-drive - 云端网盘
- 网盘功能: 文件的上传、与在线预览等网盘功能。

# 📊 系统目录结构
~~~
linyu-server/
├── assets/                  # 配置与资源文件
├── linyu-ai/                # AI服务
├── linyu-auth/              # 认证模块
├── linyu-basic/             # IM基础业务
├── linyu-cloud-drive/       # 云网盘
├── linyu-common/            # 基础设施 (DB, MQ, Log, I18n)
├── linyu-gateway/           # API统一网关
├── linyu-im/                # WebSocket通信层
└── linyu-voip-chat/         # 音视频聊天实现
~~~

# 📄 免责声明
- 本项目仅供学习与技术交流，开发者不对使用过程中产生的任何直接或间接损失负责。
- 任何基于本项目的二次开发需遵守当地法律法规。
- Linyu团队拥有对本声明的最终解释权。