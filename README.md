# iComment - 极简博客评论系统

一个简洁的博客评论系统，单一可执行文件，开箱即用。

## 特性

- ✨ 极简设计，专注核心功能
- 🚀 单一可执行文件，无需额外依赖
- 💾 SQLite 数据库，轻量高效
- 🔒 管理端独立端口，仅限 localhost 访问
- ✅ 评论审核机制，保护网站安全
- 📝 支持一层回复
- 🎨 纯净的前端组件，自适应父容器
- 🔌 RESTful API 设计

## 快速开始

### 1. 编译

```bash
go build -o icomment .
```

### 2. 运行

```bash
# 使用默认配置
# 数据库：./comments.db
# 公开端口：8080
# 管理端口：8081 (localhost only)
./icomment

# 自定义配置
./icomment -db /path/to/comments.db -port 3000 -admin-port 3001
```

### 3. 嵌入到博客

在你的博客文章页面中添加以下代码：

```html
<div id="icomment"></div>
<script src="http://your-domain.com:8080/static/comment.js" data-api="http://your-domain.com:8080"></script>
```

Replace `your-domain.com:8080` with your iComment server's domain and port.

## 访问地址

### 公开服务 (端口 8080)
- **API 端点**: `http://0.0.0.0:8080/api/comments`
- **前端脚本**: `http://0.0.0.0:8080/static/comment.js`

### 管理服务 (端口 8081, localhost only)
- **管理界面**: `http://localhost:8081/`
- **API 端点**: `http://localhost:8081/comments`

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
  "nickname": "张三",
  "email": "user@example.com",
  "content": "评论内容"
}
```

评论提交后状态为 `pending`，需要管理员审核后才会显示。

### 管理 API

#### 获取评论列表（带过滤）
```http
GET /comments?status=pending&article_url=xxx&email=xxx&page=1
```

参数：
- `status`: `all` | `pending` | `approved` (默认: `pending`)
- `article_url`: 文章 URL 前缀过滤
- `email`: 邮箱精确匹配
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

访问 `http://localhost:8081/` 可以：

- 📋 查看和过滤评论（按状态、文章、邮箱）
- ✅ 批准待审核评论
- 🗑️ 删除评论
- 📄 分页浏览（每页 10 条）
- 🔍 展开/收起评论内容
- 🔗 点击文章链接直接访问
- 📊 实时统计评论数量

## 数据库结构

```sql
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_url TEXT NOT NULL,
    parent_id INTEGER,
    nickname TEXT NOT NULL,
    email TEXT NOT NULL,
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

```nginx
# 公开 API
location /api/comments {
    proxy_pass http://localhost:8080;
}

location /static/comment.js {
    proxy_pass http://localhost:8080;
}

# 管理端 (可选，如需远程访问)
location /admin/ {
    proxy_pass http://localhost:8081/;
    # 添加 IP 白名单或 HTTP Basic Auth
    allow 192.168.1.0/24;
    deny all;
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

## 命令行参数

```bash
./icomment -h

Usage of ./icomment:
  -db string
        Path to SQLite database file (default "./comments.db")
  -port string
        Public API port (default "8080")
  -admin-port string
        Admin panel port (default "8081")
```

## License

MIT
