<<<<<<< HEAD
# Komari

![Badge](https://hitscounter.dev/api/hit?url=https%3A%2F%2Fgithub.com%2FFearless743%2Fkomari&label=&icon=github&color=%23a370f7&message=&style=flat&tz=UTC)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Fearless743/komari)

![komari](https://socialify.git.ci/Fearless743/komari/image?description=1&font=Inter&forks=1&issues=1&language=1&logo=https%3A%2F%2Fraw.githubusercontent.com%2FFearless743%2Fkomari-web%2Fd54ce1288df41ead08aa19f8700186e68028a889%2Fpublic%2Ffavicon.png&name=1&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Auto)

[简体中文](./docs/README_zh.md) | [繁體中文](./docs/README_zh-TW.md) | [日本語](./docs/README_ja.md)

Komari is a lightweight, self-hosted server monitoring tool designed to provide a simple and efficient solution for monitoring server performance. It supports viewing server status through a web interface and collects data through a lightweight agent.

[Documentation](https://komari-document.pages.dev/) | [文档(镜像站 By Geekertao)](https://www.komari.wiki) | [Telegram Group](https://t.me/komari_monitor)

## Features

- **Lightweight and Efficient**: Low resource consumption, suitable for servers of all sizes.
- **Self-hosted**: Complete control over data privacy, easy to deploy.
- **Web Interface**: Intuitive monitoring dashboard, easy to use.

## Quick Start

### 0. One-click Deployment with Cloud Hosting

- Rainyun - CNY 4.5/month

[![](https://rainyun-apps.cn-nb1.rains3.com/materials/deploy-on-rainyun-cn.svg)](https://app.rainyun.com/apps/rca/store/6780/NzYxNzAz_)

- 1Panel App Store

Available on 1Panel App Store. Install via **App Store > Utilities > Komari**.

### 1. Use the One-click Install Script

Suitable for distributions using systemd (Ubuntu, Debian...).

```bash
curl -fsSL https://raw.githubusercontent.com/Fearless743/komari/main/install-komari.sh -o install-komari.sh
chmod +x install-komari.sh
sudo ./install-komari.sh
```

### 2. Docker Deployment

1. Create a data directory:
   ```bash
   mkdir -p ./data
   ```
2. Run the Docker container:
   ```bash
   docker run -d \
     -p 25774:25774 \
     -v $(pwd)/data:/app/data \
     --name komari \
     ghcr.io/Fearless743/komari:latest
   ```
3. View the default username and password:
   ```bash
   docker logs komari
   ```
4. Access `http://<your_server_ip>:25774` in your browser.

> [!NOTE]
> You can also customize the initial username and password through the environment variables `ADMIN_USERNAME` and `ADMIN_PASSWORD`.

### 3. Binary File Deployment

1. Visit Komari's [GitHub Release page](https://github.com/Fearless743/komari/releases) to download the latest binary for your operating system.
2. Run Komari:
   ```bash
   ./komari server -l 0.0.0.0:25774
   ```
3. Access `http://<your_server_ip>:25774` in your browser. The default port is `25774`.
4. The default username and password can be found in the startup logs or set via the environment variables `ADMIN_USERNAME` and `ADMIN_PASSWORD`.

> [!NOTE]
> Ensure the binary has execute permissions (`chmod +x komari`). Data will be saved in the `data` folder in the running directory.

### Manual Build

#### Dependencies

- Go 1.18+ and Node.js 20+ (for manual build)

1. Build the frontend static files:
   ```bash
   git clone https://github.com/Fearless743/komari-web
   cd komari-web
   npm install
   npm run build
   ```
2. Build the backend:
   ```bash
   git clone https://github.com/Fearless743/komari
   cd komari
   ```
   Copy the static files generated in step 1 to the `/public/defaultTheme/dist` folder in the root of the `komari` project, and copy `komari-theme.json` + `preview.png`/`perview.png` to `/public/defaultTheme`.
   ```bash
   go build -o komari
   ```
3. Run:
   ```bash
   ./komari server -l 0.0.0.0:25774
   ```
   The default listening port is `25774`. Access `http://localhost:25774`.

## Frontend Development Guide

[Komari Theme Development Guide | Komari](https://komari-document.pages.dev/dev/theme.html)

## Client Agent Development Guide

[Komari Agent Information Reporting and Event Handling Documentation](https://komari-document.pages.dev/dev/agent.html)

## Contributing

Issues and Pull Requests are welcome!

## Acknowledgements

### 破碎工坊云

[破碎工坊云 - 专业云计算服务平台，提供高效、稳定、安全的高防服务器与CDN解决方案](https://www.crash.work/)

### DreamCloud

[DreamCloud - 极高性价比解锁直连亚太高防](https://as211392.com/)

### 🚀 Sponsored by SharonNetworks

[![Sharon Networks](https://raw.githubusercontent.com/Fearless743/public/refs/heads/main/images/sharon-networks.webp)](https://sharon.io)

SharonNetworks 为您的业务起飞保驾护航！

亚太数据中心提供顶级的中国优化网络接入 · 低延时&高带宽&提供Tbps级本地清洗高防服务, 为您的业务保驾护航, 为您的客户提供极致体验. 加入社区 [Telegram群组](https://t.me/SharonNetwork) 可参与公益募捐或群内抽奖免费使用

### The open source software community

All the developers who submitted PRs and created themes

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Fearless743/komari&type=Date)](https://www.star-history.com/#Fearless743/komari&Date)
=======
# Komari Web UI

参与翻译Komari？
- 直接提PR

We use AI to assist with translations. If you find any issues, please let us know!

How to contribute to Komari translations?
- Directly PR

## 开发环境配置

> 我不是计科专业的，代码质量可能达不到平均水平，React是边学边写的，在此之前我从未接触过前端开发，请多包涵。

### 前置 Nodejs

如果未安装，请访问 [Node.js 官网](https://nodejs.org/) 下载并安装。版本建议为 22 及以上。

### 安装依赖

```bash
npm install
```

> 所有指令均在项目根目录下执行

### 修改API地址

1. 复制 `.env.example` 文件并重命名为 `.env.development`。

2. 修改 `.env.development` 文件中的 `VITE_API_TARGET` 为你的开发环境地址。

### 启动开发服务器

```bash
npm run dev
```

### 构建

```bash
npm run build
```

## 主题相关

如果你需要基于本项目进行二次开发，可以参考以下步骤：

1. 完成开发环境配置

> 如果你是在 Linux 系统下开发，可以直接运行脚本 `build-theme.sh` 快速生成主题包。

2. 修改 `komari-theme.json` 中的相关配置，具体可参考 [主题配置文件 | Komari](https://komari-document.pages.dev/dev/theme.html#%E4%B8%BB%E9%A2%98%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)

3. 发挥你的想象和创造力，设计并实现你独特的主题风格！

4. 构建主题

   ```bash
   npm run build
   ```

5. 生成的主题文件位于 `dist` 目录下，创建一个新的文件夹 `my-theme`（名称自定），将 `dist` 目录下复制到 `my-theme` 文件夹中。

6. 将 `komari-theme.json` 文件复制到 `my-theme` 文件夹中。

7. 将 `my-theme` 文件夹打包为 ZIP 文件。

8. 在 Komari 的主题管理页面上传并应用你的自定义主题。
>>>>>>> origin/radix
