IP CIDR downloader 1.0v / IP地址范围生成器与国家代码查询工具 1.0v

一个用于生成各国IP地址范围（CIDR格式）并提供国家代码查询的工具，支持从五大区域互联网注册管理机构（RIR）获取最新IP分配数据。

## 功能特点

- 从APNIC、ARIN、RIPE NCC、AFRINIC和LACNIC获取最新IP分配数据
- 按国家/地区生成IPv4和IPv6地址范围（CIDR格式）
- 提供国家代码查询功能，支持单个查询和批量列出
- 数据自动按国家分类存储，便于快速访问
- 每个IP文件包含更新时间和作者信息

## 使用方法

### 基本命令

```bash
# 显示帮助信息
go run ipcidr.go -h

# 查询单个国家代码信息
go run ipcidr.go -name CN

# 查看所有国家代码（默认每行5个）
go run ipcidr.go -name all

# 自定义每行显示的国家代码数量
go run ipcidr.go -name all 3

# 下载并生成最新IP数据
go run ipcidr.go -update
```

### 输出文件结构

执行更新命令后，会在当前目录创建`data`文件夹，其结构如下：

```
data/
├── cn/
│   ├── ipv4.txt
│   └── ipv6.txt
├── us/
│   ├── ipv4.txt
│   └── ipv6.txt
...
```

每个国家代码对应的文件夹中包含两个文件，分别存储该国家/地区的IPv4和IPv6地址范围。

文件内容示例（`data/cn/ipv4.txt`）：
```
# last updated: 2023-10-15T10:30:45+08:00
1.0.1.0/24
1.0.2.0/23
...
```

## 在Ubuntu 22.04上部署教学

### 1. 安装必要依赖

```bash
# 更新系统包
sudo apt update && sudo apt upgrade -y

# 安装Go语言环境
sudo apt install -y golang-go

# 验证Go安装
go version  # 应输出类似 "go version go1.18.1 linux/amd64" 的信息
```

### 2. 获取代码

```bash
# 克隆代码库（如果使用Git管理）
git clone https://github.com/saycc1982/ipcidr.git
cd ipcidr

# 或者直接创建文件
nano ipcidr.go
# 将代码粘贴到文件中，按Ctrl+O保存，Ctrl+X退出
```

### 3. 运行程序

```bash
# 查看帮助
go run ipcidr.go -h

# 首次运行更新数据（这可能需要几分钟时间）
go run ipcidr.go -update
```

### 4. 可选：编译为可执行文件

```bash
# 编译
go build -o ipcidr ipcidr.go

# 赋予执行权限
chmod +x ipcidr

# 使用编译后的文件
./ipcidr -update
./ipcidr -name JP
```

## 可能遇到的问题及解决方法

### 1. 网络连接问题

**症状**：下载数据时出现超时或连接失败

**解决方法**：
```bash
# 检查网络连接
ping -c 3 ftp.apnic.net

# 如使用代理，设置环境变量
export http_proxy=http://your-proxy:port
export https_proxy=https://your-proxy:port
```

### 2. 权限问题

**症状**：无法创建目录或文件，出现"permission denied"错误

**解决方法**：
```bash
# 检查当前目录权限
ls -ld .

# 如必要，更改目录权限
chmod 755 .

# 或使用sudo运行（不推荐，仅在必要时使用）
sudo go run ipcidr.go -update
```

### 3. AFRINIC数据下载失败

**症状**：AFRINIC数据下载多次尝试后失败

**解决方法**：
程序会自动尝试30天内的历史数据，如果仍失败，会尝试备选URL。如果所有尝试都失败，程序会忽略AFRINIC数据继续处理其他来源。

### 4. 内存不足

**症状**：处理大量数据时程序崩溃或被系统终止

**解决方法**：
```bash
# 检查系统内存使用情况
free -m

# 如内存不足，可增加系统交换空间
sudo dd if=/dev/zero of=/swapfile bs=1M count=2048
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 5. Go语言版本过低

**症状**：编译时出现语法错误

**解决方法**：
```bash
# 安装较新版本的Go
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go
```

## 反馈与支持

如有建议或BUG反馈，欢迎联络作者：
- 作者: ED
- 联络: https://t.me/hongkongisp
- 服务器推荐: 顺安云 https://www.say.cc
