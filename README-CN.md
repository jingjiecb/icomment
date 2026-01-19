[English](README.md) | [中文](README-CN.md)

# iComment - 极简博客评论系统

一个专为博客设计的轻量级单文件评论系统。

## ✨ 特性

- **极简主义**: 单一可执行文件，SQLite 数据库，无复杂依赖。
- **Docker 就绪**: 提供一流的 Docker 支持和自动构建。
- **安全**: 注重隐私，包含评论审核工作流。
- **易于集成**: 任何网站只需嵌入简单的 HTML/JS 代码。
- **通知**: 可选 Bark (iOS) 推送通知。

---

## 📖 用户指南

### 快速开始 (Docker) - 推荐

使用提供的 `docker-compose.yml` 即可立即开始。

```yaml
services:
  icomment:
    image: claws059/icomment:latest
    restart: unless-stopped
    ports:
      - "7001:7001" # 公开 API
      - "7002:7002" # 管理界面 (务必保护此端口!)
    volumes:
      - ./data:/data
    environment:
      - ICOMMENT_BARK=your_bark_key # 可选
```

启动服务：

```bash
docker-compose up -d
```

### 配置

你可以通过环境变量 (Docker) 或命令行参数配置 iComment。

| 环境变量 | 参数标志 | 默认值 | 说明 |
|----------|----------|--------|------|
| `ICOMMENT_DB` | `-db` | `./comments.db` | SQLite 数据库路径 |
| `ICOMMENT_PORT` | `-port` | `7001` | 公开 API 和静态文件端口 |
| `ICOMMENT_ADMIN_PORT` | `-admin-port` | `7002` | 管理控制台端口 |
| `ICOMMENT_BARK` | `-bark` | (空) | 用于通知的 Bark 设备 Key |

### 集成

在你的博客模板中添加以下代码段：

```html
<div id="icomment"></div>
<script src="http://your-domain.com:7001/static/comment.js" 
        data-api="http://your-domain.com:7001"></script>
```

### 管理与安全

访问管理控制台： `http://your-domain.com:7002/`。

> [!WARNING]
> **安全警告**: 管理控制台 (端口 7002) 会监听 **所有接口 (0.0.0.0)**。
> 在生产环境中，你 **必须** 使用防火墙或反向代理 (如 Nginx 带 Basic Auth) 来保护此端口。**不要在无保护的情况下将其暴露在公网上。**

**Nginx 保护示例:**

```nginx
location /admin/ {
    proxy_pass http://127.0.0.1:7002/;
    # 启用 Basic Auth
    auth_basic "Restricted Admin Area";
    auth_basic_user_file /etc/nginx/.htpasswd;
}
```

---

## 🛠 开发者指南

### 从源码构建

```bash
# 创建数据目录
mkdir data

# 构建
go build -o icomment .

# 运行
./icomment
```

### 项目结构

```text
.
├── main.go              # 入口点和服务器设置
├── config.go            # 配置加载 (Env/Flags)
├── handler.go           # HTTP 请求处理器
├── Dockerfile           # 多阶段 Docker 构建
├── .github/workflows/   # CI/CD (Docker 发布)
├── model/               # 数据结构
├── dao/                 # 数据库访问层
├── static/              # 前端资源 (JS)
└── templates/           # 管理界面模板
```

### API 参考

**公开端点 (端口 7001)**
- `GET /api/comments?article_url=...` - 获取评论列表
- `POST /api/comments` - 提交评论

**管理端点 (端口 7002)**
- `GET /comments` - 获取/过滤评论
- `PATCH /comments/:id/approve` - 批准评论
- `DELETE /comments/:id` - 删除评论

### 数据库模式

iComment 使用 SQLite。

```sql
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_url TEXT NOT NULL,
    nickname TEXT NOT NULL,
    email TEXT,
    content TEXT NOT NULL,
    status TEXT DEFAULT 'pending', -- pending/approved
    ...
);
```

---

## License

MIT
