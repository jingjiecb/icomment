# iComment - 极简博客评论系统

一个简洁的博客评论系统，单一可执行文件，开箱即用。

## 特性

- ✨ 极简设计，专注核心功能
- 🚀 单一可执行文件，无需额外依赖
- 💾 SQLite 数据库，轻量高效
- 🔒 管理端独立端口，仅限 localhost 访问
- ✅ 评论审核机制，保护网站安全
- 📝 支持一层回复
- 🎨 纯净的前端组件，自适应父容器，支持暗黑模式
- 🔌 RESTful API 设计
- 🔐 隐私保护：Email 地址不对外暴露
- ⚡ 静态资源缓存优化
- 🔔 Bark 推送通知（可选）

## 快速开始

### 1. 编译

```bash
go build -o icomment .
```

### 2. 运行

```bash
# 使用默认配置
# 数据库：./comments.db
# 公开端口：7001
# 管理端口：7002 (localhost only)
./icomment

# 自定义配置
./icomment -db /path/to/comments.db -port 3000 -admin-port 3001
```

### 3. 嵌入到博客

在你的博客文章页面中添加以下代码：

```html
<div id="icomment"></div>
<script src="http://your-domain.com:7001/static/comment.js" data-api="http://your-domain.com:7001"></script>
```

Replace `your-domain.com:7001` with your iComment server's domain and port.

## 访问地址

### 公开服务 (端口 7001)
- **API 端点**: `http://0.0.0.0:7001/api/comments`
- **前端脚本**: `http://0.0.0.0:7001/static/comment.js`

### 管理服务 (端口 7002, localhost only)
- **管理界面**: `http://localhost:7002/`
- **API 端点**: `http://localhost:7002/comments`

## API 文档

### 公开 API

#### 获取评论列表
```http
GET /api/comments?article_url=<article_url>
```

返回指定文章的所有已批准评论。

#### 提交评论
```http
POST /api/comments
Content-Type: application/json

{
  "article_url": "https://example.com/post",
  "parent_id": 123,  // 可选，回复评论的 ID
  "nickname": "张三",  // 必填，最多 50 字符
  "email": "user@example.com",  // 可选，最多 100 字符
  "content": "评论内容"  // 必填，最多 2000 字符
}
```

评论提交后状态为 `pending`，需要管理员审核后才会显示。

**响应示例：**
```json
{
  "message": "Comment created, pending approval"
}
```

### 管理 API

#### 获取评论列表（带过滤）
```http
GET /comments?status=pending&article_url=xxx&page=1
```

参数：
- `status`: `all` | `pending` | `approved` (默认: `pending`)
- `article_url`: 文章 URL 前缀过滤
- `page`: 页码 (默认: 1, 每页 10 条)

#### 批准评论
```http
PATCH /comments/:id/approve
```

#### 删除评论
```http
DELETE /comments/:id
```

## 管理功能

访问 `http://localhost:7002/` 可以：

- 📋 查看和过滤评论（按状态、文章 URL）
- ✅ 批准待审核评论（一键操作）
- 🗑️ 删除评论（级联删除回复）
- 📄 分页浏览（每页 10 条）
- 🔍 展开行查看完整详情
- 🔗 点击文章链接直接访问
- 📊 实时统计评论数量

## 数据库结构

```sql
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_url TEXT NOT NULL,
    parent_id INTEGER,
    nickname TEXT NOT NULL,
    email TEXT,  -- 可选
    content TEXT NOT NULL,
    status TEXT DEFAULT 'pending',  -- 'pending' 或 'approved'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_status ON comments(status);
CREATE INDEX idx_article_url ON comments(article_url);
CREATE INDEX idx_email ON comments(email);
```

## 项目结构

```
.
├── main.go              # 主程序入口
├── config.go            # 配置管理
├── db.go                # 数据库初始化
├── handler.go           # HTTP 处理器
├── model/
│   └── comment.go       # 数据模型
├── dao/
│   └── comment.go       # 数据访问层
├── static/
│   └── comment.js       # 前端组件
├── templates/
│   └── admin.html       # 管理界面
└── sql/
    └── create_table_comment.sql
```

## 设计理念

- **Less is more** - 只保留核心功能
- **Simple is better than complex** - 简单直接的实现
- **开箱即用** - 默认配置即可运行
- **安全第一** - 评论审核机制，管理端隔离

## 安全特性

- ✅ **评论审核**: 所有评论默认待审核，只有批准后才会在前端显示
- 🔒 **管理端隔离**: 管理 API 运行在独立端口，仅监听 localhost
- 🛡️ **CORS 支持**: 公开 API 支持跨域，方便嵌入任何网站
- 🔗 **级联删除**: 删除评论时自动删除其回复

## 部署建议

### 反向代理配置 (Nginx)

#### 基础配置

```nginx
# 公开 API
location /api/comments {
    proxy_pass http://localhost:7001;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}

location /static/comment.js {
    proxy_pass http://localhost:7001;
    proxy_set_header Host $host;
}
```

#### 管理端远程访问配置

默认情况下，管理端只能通过 `localhost:7002` 访问。如果需要远程访问管理界面，可以通过 Nginx 反向代理并添加安全认证。

**方案一：IP 白名单**

```nginx
location /admin/ {
    proxy_pass http://localhost:7002/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    
    # 只允许特定 IP 访问
    allow 192.168.1.0/24;  # 允许内网
    allow 203.0.113.10;    # 允许特定公网 IP
    deny all;              # 拒绝其他所有 IP
}
```

**方案二：HTTP Basic Auth（推荐）**

1. 创建密码文件：

```bash
# 安装 htpasswd 工具（如果没有）
# Ubuntu/Debian
sudo apt-get install apache2-utils

# CentOS/RHEL
sudo yum install httpd-tools

# macOS (通常已预装)
# 如果没有：brew install httpd

# 创建密码文件和第一个用户
sudo htpasswd -c /etc/nginx/.htpasswd admin

# 添加更多用户（不要使用 -c 参数，会覆盖文件）
sudo htpasswd /etc/nginx/.htpasswd another_user
```

2. 配置 Nginx：

```nginx
location /admin/ {
    proxy_pass http://localhost:7002/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    
    # 启用 HTTP Basic Auth
    auth_basic "iComment Admin Area";
    auth_basic_user_file /etc/nginx/.htpasswd;
}
```

3. 重载 Nginx 配置：

```bash
# 测试配置文件语法
sudo nginx -t

# 重载配置
sudo nginx -s reload
```

#### HTTPS 配置（强烈推荐）

如果管理端需要远程访问，强烈建议启用 HTTPS：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    # 公开 API（无需认证）
    location /api/comments {
        proxy_pass http://localhost:7001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /static/comment.js {
        proxy_pass http://localhost:7001;
        proxy_set_header Host $host;
    }
    
    # 管理端（需要认证）
    location /admin/ {
        proxy_pass http://localhost:7002/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        
        auth_basic "iComment Admin Area";
        auth_basic_user_file /etc/nginx/.htpasswd;
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

### 使用 systemd 管理服务

```ini
[Unit]
Description=iComment Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/icomment
ExecStart=/opt/icomment/icomment -db /var/lib/icomment/comments.db
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

**重要：确保数据库目录权限正确**
```bash
# 创建数据目录
mkdir -p /var/lib/icomment

# 设置正确的所有者（与 systemd User 一致）
chown www-data:www-data /var/lib/icomment

# 确保有写权限
chmod 755 /var/lib/icomment
```

## 命令行参数

```bash
./icomment -h

Usage of ./icomment:
  -db string
        Path to SQLite database file (default "./comments.db")
  -port string
        Public API port (default "7001")
  -admin-port string
        Admin panel port (default "7002")
  -bark string
        Bark device key for notifications (optional)
```

### Bark 推送通知

iComment 支持通过 [Bark](https://github.com/Finb/Bark) 接收新评论通知。

**启用方式：**

```bash
# 启动时添加 -bark 参数，填入你的 Bark Device Key
./icomment -bark "your_bark_device_key"
```

**获取 Bark Device Key：**

1. 在 iOS App Store 下载 Bark 应用
2. 打开应用，复制显示的 Device Key（格式如：`aBcDeFgHiJkLmN`）
3. 将 Device Key 作为 `-bark` 参数传入

**通知内容：**

- 新评论时推送标题："新评论"
- 新回复时推送标题："新回复"
- 通知内容包含：昵称 + 评论内容预览（前 100 字符）
- 点击通知可直接跳转到文章页面
- 所有通知归类到 "iComment" 分组

**示例：**

```bash
# 完整配置示例
./icomment \
  -db /var/lib/icomment/comments.db \
  -port 7001 \
  -admin-port 7002 \
  -bark "aBcDeFgHiJkLmN"
```

启动后，每当有新评论提交时，你的 iPhone 会立即收到推送通知。

## License

MIT
