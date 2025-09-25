# 域名跳转系统 (Domain Jumper)

一个支持动态多级域名跳转的Web系统，可以按照预设的路径在多个域名池之间进行跳转，最终到达目标域名。

## 🚀 功能特性

- **动态域名池支持**: 支持任意数量的中间域名池
- **智能域名检测**: 自动检测每个域名池中可用的域名
- **安全随机跳转**: 使用加密级随机字符串替换子域名
- **实时进度显示**: 显示当前跳转进度和状态
- **容错机制**: 自动处理不可用域名，选择备用域名
- **分布式部署**: 支持多服务器部署

## 📋 跳转流程

系统按照以下流程进行跳转：

```
startDomain → pool1 → pool2 → ... → poolN → target_pool → 最终目标
```

- **中间池** (`type: "intermediate"`): 可以有任意多个，依次跳转
- **目标池** (`type: "target"`): 只有一个，是跳转的终点

## 📁 项目结构

```
distributed_test/
├── server1/                 # 服务器1
│   ├── config/
│   │   └── domains.json     # 域名池配置
│   └── static/
│       ├── index.html       # 主页面
│       └── jumper.js        # 跳转逻辑
├── server2/                 # 服务器2 (相同结构)
├── server3/                 # 服务器3 (相同结构)
└── README.md               # 本文档
```

## ⚙️ 配置文件

### domains.json 配置示例

```json
{
  "startDomain": "http://test1.localhost:8001",
  "domainPools": [
    {
      "name": "pool1",
      "type": "intermediate",
      "domains": [
        "http://pool1-1.test1.localhost:8001",
        "http://pool1-2.test1.localhost:8001",
        "http://pool1-3.test1.localhost:8001"
      ]
    },
    {
      "name": "pool2", 
      "type": "intermediate",
      "domains": [
        "http://pool2-1.test2.localhost:8002",
        "http://pool2-2.test2.localhost:8002"
      ]
    },
    {
      "name": "target_pool",
      "type": "target",
      "domains": [
        "http://target1.final.localhost:8003",
        "http://target2.final.localhost:8003"
      ]
    }
  ],
  "timeoutMs": 3000,
  "jumpDelayMs": 1000
}
```

### 配置字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `startDomain` | string | 起始域名 |
| `domainPools` | array | 域名池数组 |
| `name` | string | 池子名称 |
| `type` | string | 池子类型：`intermediate`(中间池) 或 `target`(目标池) |
| `domains` | array | 该池子中的域名列表 |
| `timeoutMs` | number | 域名检测超时时间（毫秒） |
| `jumpDelayMs` | number | 跳转延迟时间（毫秒） |

## 🔧 部署说明

### 1. 基本部署

每个server目录都是一个独立的服务，可以部署到不同的服务器上：

```bash
# 进入任一server目录
cd server1

# 启动简单HTTP服务器 (Python)
python3 -m http.server 8001

# 或者使用Node.js
npx http-server -p 8001
```

### 2. 配置域名池

根据您的需求修改 `config/domains.json`：

1. **添加中间池**: 在 `domainPools` 数组中添加更多 `type: "intermediate"` 的池子
2. **设置目标池**: 确保有且仅有一个 `type: "target"` 的池子
3. **配置域名**: 在每个池子的 `domains` 数组中添加实际可用的域名

### 3. 多服务器同步

所有server目录的配置和代码应该保持同步：

```bash
# 同步配置文件
cp server1/config/domains.json server2/config/domains.json
cp server1/config/domains.json server3/config/domains.json

# 同步代码文件
cp server1/static/jumper.js server2/static/jumper.js  
cp server1/static/jumper.js server3/static/jumper.js
```

## 🚦 使用方法

### 1. 启动系统

1. 部署并启动各个服务器
2. 确保所有域名都能正常访问
3. 访问起始域名开始跳转流程

### 2. 跳转过程

1. **初始化**: 系统加载配置，计算总步骤数
2. **域名检测**: 并发检测当前池子中所有域名的可用性
3. **选择域名**: 选择第一个可用的域名
4. **随机跳转**: 使用随机字符串替换子域名后跳转
5. **重复流程**: 在每个服务器上重复上述过程
6. **到达目标**: 到达目标池后进行最终跳转

## 🔒 安全特性

### 随机字符串生成

系统使用 `window.crypto.getRandomValues()` 生成加密级随机字符串：

```javascript
function generateRandomString() {
    const length = Math.floor(Math.random() * 4) + 3; // 3-6位随机长度
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const randomValues = new Uint32Array(length);
    window.crypto.getRandomValues(randomValues);
    
    let result = '';
    for (let i = 0; i < length; i++) {
        result += chars[randomValues[i] % chars.length];
    }
    return result;
}
```

### 子域名替换

跳转时会自动替换子域名：
```
原域名: http://pool1-1.test1.localhost:8001
跳转后: http://Kx9aB.test1.localhost:8001
```

## 🐛 故障排除

### 常见问题

1. **域名不可用**
   - 系统会自动尝试池子中的其他域名
   - 如果整个池子都不可用，会显示错误信息

2. **跳转卡住**
   - 检查网络连接
   - 确认目标域名是否正确配置
   - 查看浏览器控制台的错误信息

3. **配置错误**
   - 确保JSON格式正确
   - 检查域名格式是否有效
   - 确保有且仅有一个 `type: "target"` 的池子

### 调试模式

打开浏览器开发者工具，在Console中可以看到详细的跳转日志。

## 📝 开发说明

### 代码结构

- **DomainJumper类**: 处理初始跳转逻辑
- **IntermediateJumper类**: 处理中间跳转逻辑
- **域名检测**: 使用fetch检测域名可用性
- **进度管理**: 动态计算和显示跳转进度

### 扩展功能

如需添加新功能，主要修改 `jumper.js` 文件：

1. **修改检测逻辑**: 修改 `checkDomainAvailability` 方法
2. **自定义跳转规则**: 修改 `startJumpSequence` 方法  
3. **UI增强**: 修改相关的状态更新方法

## 📜 许可证

本项目仅供学习和研究使用。

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个项目。

---

**注意**: 请确保在合法和授权的环境中使用此系统。