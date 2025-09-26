# Nginx日志系统配置说明

本文档详细说明商户服务中Nginx日志系统的配置、使用方法和管理策略。

## 日志文件结构

Nginx配置了以下几类日志文件：

- **全局访问日志**：`/var/log/nginx/access.log` - 记录所有访问请求
- **商户访问日志**：`/var/log/nginx/merchant_access.log` - 记录商户相关的访问请求，包含商户ID
- **错误日志**：`/var/log/nginx/error.log` - 记录所有错误信息
- **慢速请求日志**：`/var/log/nginx/slow.log` - 记录响应时间超过1秒的请求
- **上传服务日志**（若启用）：`/var/log/nginx/upload_access.log` - 记录上传相关请求

## 日志格式说明

### 1. 主日志格式 (main)
包含基本的访问信息：
- 客户端IP地址
- 请求用户
- 请求时间
- 请求URL和方法
- 响应状态码
- 响应大小
- 来源页面
- 用户代理信息
- X-Forwarded-For信息

### 2. 商户访问日志格式 (merchant)
在主日志格式基础上增加了：
- merchant_id: 提取的商户ID
- uri: 请求的完整路径

### 3. 慢速请求日志格式 (slow)
在主日志格式基础上增加了：
- request_time: 请求处理时间（秒）

## 日志轮转配置

为防止日志文件过大，系统配置了自动轮转策略：

- 每天自动轮转一次
- 保留14天的历史日志
- 旧日志会自动压缩
- 轮转后自动通知Nginx重新打开日志文件

轮转配置文件位置：`nginx/logrotate.conf`

## 查看日志的常用命令

```bash
# 实时查看访问日志
tail -f /var/log/nginx/access.log

# 实时查看商户访问日志
tail -f /var/log/nginx/merchant_access.log

# 实时查看错误日志
tail -f /var/log/nginx/error.log

# 查看慢速请求日志
tail -f /var/log/nginx/slow.log

# 统计特定商户的访问次数
grep 'merchant_id: 1000' /var/log/nginx/merchant_access.log | wc -l

# 查找特定状态码的请求
grep ' 404 ' /var/log/nginx/access.log
```

## 日志分析建议

1. **商户访问分析**：使用`merchant_access.log`分析各商户的访问量、热门资源和访问模式
2. **性能监控**：通过`slow.log`识别响应缓慢的请求，进行性能优化
3. **错误排查**：通过`error.log`排查系统错误和异常情况
4. **安全审计**：定期检查日志中的异常访问模式和可疑请求

## 注意事项

1. 确保日志目录存在且Nginx有写入权限
2. 定期检查磁盘空间，特别是日志保留策略可能需要根据实际情况调整
3. 日志中可能包含敏感信息，请确保适当的访问控制
4. 对于生产环境，建议考虑使用专业的日志收集和分析工具（如ELK Stack）