# 🍔 麦当劳订单管理系统

基于 **Go** 后端和 **Vue 3** 前端的订单处理系统，支持 **VIP 优先级排队**、**动态增减机器人**、**订单完成时间记录**。  
**双运行模式**：默认启动 HTTP 服务器（配合前端界面），支持 `--cli` 参数启动交互式命令行。

## ✨ 功能特性

- **订单优先级**：VIP 订单自动排在普通订单之前，VIP 内部 FIFO，普通内部 FIFO
- **机器人厨师**：
  - 每个机器人同时只能处理一个订单，每个订单耗时 **10 秒**
  - 可动态增加或减少机器人，机器人空闲时自动从待处理区取单
  - 减少机器人时，若该机器人正在处理订单，则订单**原样归还**到队列头部（保持优先级）
- **实时状态**：默认 HTTP 模式提供 REST API，配合 Vue 3 前端界面实时刷新
- **交互式 CLI**：通过 `--cli` 参数启动纯命令行界面，支持完整的管理命令
- **订单记录**：所有订单完成时间以 `HH:MM:SS` 格式记录在 `result.txt` 文件中
- **订单号规则**：20 位唯一编码（17 位毫秒时间戳 + 3 位序列号），趋势递增

## 🏗️ 项目结构
```text
.
├── backend/ # Go 后端服务
│ ├── main.go # 主程序（双模式：HTTP / CLI）
│ ├── main_test.go # 单元测试（优先级、归还逻辑）
│ ├── go.mod
│ ├── build.sh # 编译脚本
│ ├── run.sh # 运行脚本（默认启动 HTTP 服务器）
│ ├── run-cli.sh # 运行脚本（启动 CLI 交互模式）
│ ├── test.sh # 测试脚本
│ └── result.txt # 运行时生成的订单完成日志
├── frontend/
│ └── index.html # Vue 3 前端界面（支持中英文切换）
└── README.md
```

## 🚀 快速开始

### 编译后端

```bash
cd backend
chmod +x build.sh run.sh run-cli.sh test.sh
./build.sh          # 编译 order-system 二进制
```

### 运行模式
#### 1. HTTP 服务器模式（默认，配合前端界面）
```bash
   ./run.sh            # 启动 HTTP 服务器，监听 :8080
```
然后打开前端页面：直接在浏览器
中打开 frontend/index.html，或使用静态服务器（如 npx serve frontend）。
前端会轮询后端 API 并实时展示订单状态。
#### 2. 交互式 CLI 模式
```bash
./run-cli.sh        # 启动纯命令行界面
```
CLI 模式下支持以下命令（支持简写）：
```
命令	简写	说明
new normal	n normal	添加普通订单
new vip	n vip	添加 VIP 订单
add bot	ab	增加一个机器人
remove bot	rb	移除最新添加的机器人
list pending	lp	列出待处理订单
list done	ld	列出已完成订单
status	s	查看系统状态
exit	q	退出程序
```

运行示例：
```text
🍔 麦当劳订单处理系统 CLI
命令列表 (支持简写):
  new normal / n normal     - 添加普通订单
  new vip    / n vip        - 添加VIP订单
  add bot    / ab           - 增加一个机器人
  remove bot / rb           - 移除一个机器人
  list pending / lp         - 列出待处理订单
  list done    / ld         - 列出已完成订单
  status       / s          - 查看系统状态
  exit         / q          - 退出程序

> n vip
🌟 VIP订单 20250428143025001 已添加到待处理区（优先级高于普通订单）
> add bot
🤖 + 机器人 #2 已添加，当前机器人总数: 2
> status

========== 系统状态 ==========
机器人数量: 2
待处理 VIP 订单: 1, 普通订单: 0
已完成订单总数: 0
处理中订单: 20250428143025001 
==============================
```

### 运行测试
```bash
cd backend
./test.sh       # 执行单元测试，验证订单优先级和机器人归还逻辑
```

📋 订单完成日志
所有订单完成时，会将记录追加到 backend/result.txt，格式示例：
```text
订单 20250428143025001 (VIP) 完成时间 14:30:35
订单 20250428143025002 (普通) 完成时间 14:30:45
```

🔌 HTTP API 端点（默认模式）
```text
方法	端点	说明
POST	/api/order/normal	添加普通订单
POST	/api/order/vip	添加 VIP 订单
POST	/api/bot/add	增加一个机器人
POST	/api/bot/remove	移除最近添加的机器人
GET	/api/state	获取当前状态（待处理、已完成、机器人数等）
```
### 🧪 单元测试说明
 main_test.go 包含两个核心测试：
- TestOrderPriority：验证订单优先级（VIP 先于普通，且各自保持 FIFO）
- TestReturnOrder：验证机器人被移除时，正在处理的订单能正确返回到队列头部

### 🤖 GitHub Actions 工作流
创建 .github/workflows/backend-verify.yml 以自动验证后端代码：
```yaml
name: Backend Verify
on:
  pull_request:
    paths:
      - 'backend/**'
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: cd backend && ./build.sh
      - name: Test
        run: cd backend && ./test.sh
```

### 📝 注意事项
后端使用内存存储，无数据持久化（满足需求）
机器人数量没有上限，减少机器人时遵循 LIFO（后进先出）
订单号 20 位，基于毫秒时间戳 + 序列号，保证唯一性且趋势递增
前端通过 轮询（每秒）更新状态，避免 WebSocket 复杂性
CLI 模式下所有输出（包括订单完成）均带时间戳并写入 result.txt

### 🛠️ 技术栈
后端：Go 1.21 + 标准库（net/http, sync, time, flag）
前端：Vue 3 (CDN) + 原生 CSS（响应式设计，支持中英文切换）
并发模型：机器人 goroutine + sync.Cond 高效等待订单

### 👨‍🍳 设计决策
双模式设计：默认 HTTP 保留前端交互，--cli 参数满足面试对 CLI 应用程序的强制要求
使用 sync.Cond 实现机器人阻塞等待订单，避免 CPU 空转
机器人 goroutine 监听取消信号，保证移除时优雅中止并归还订单
订单队列分离为 vipQueue 和 normalQueue，获取时先取 VIP 后取普通