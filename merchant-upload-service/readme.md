1. 1.
   修改Nginx配置文件 ，实现了以下核心功能：
   
   - 使用正则表达式 ^/merchant_([a-zA-Z0-9_]+)$ 动态匹配商户路径
   - 添加 absolute_redirect off; 和 port_in_redirect off; 指令解决301重定向问题
   - 实现 /merchant([0-9]+) 格式到 /merchant_$1 的内部重写规则
   - 为商户页面、静态资源和数据文件分别配置了专门的location块
2. 2.
   修复数据文件访问问题 ：
   
   - 采用更精确的正则表达式 ^/merchant_([a-zA-Z0-9_]+)/data/(.*)$ 捕获完整的文件路径
   - 使用捕获组 $2 确保正确构建文件路径
   - 使用 root 指令配合 try_files 替代 alias 指令，提高路径处理的可靠性
3. 3.
   创建测试商户目录和文件 ：
   
   - 为测试动态配置创建了 /merchant_1001 目录结构
   - 创建了 index.html 页面和 test_data.json 数据文件
所有测试都已通过，包括：

- 访问 /merchant_1001 返回正确的商户页面
- 访问 /merchant1001 自动内部重写到 /merchant_1001
- 访问 /merchant_1001/data/test_data.json 成功返回数据文件内容
这个方案满足了您动态添加商户的需求，只需创建符合命名规则的商户目录结构，Nginx就会自动处理相应的URL请求，无需修改配置文件。





1. 1.
   修改Nginx配置 ：
   
   - 添加了两个新的location块，分别处理纯数字ID和纯字符ID带路径的访问（如/1001/abc.html和/abc001/jumper.js）
   - 在try_files指令中添加了/html目录的路径，使其能够正确查找商户目录下html子目录中的文件
   - 确保新规则不会与现有的数据文件访问规则冲突
2. 2.
   功能验证 ：
   
   - 成功访问/1001/abc.html，返回了正确的HTML内容
   - 成功访问/abc001/jumper.js，返回了正确的JavaScript文件内容
   - 确认现有的数据文件访问功能（如/merchant_1001/data/test_data.json）仍然正常工作
现在系统完全支持以下四种URL格式的商户页面和资源访问：

- 格式1：带下划线的完整格式（/merchant_{id}）
- 格式2：不带下划线的格式（/merchant{id}，仅数字ID）
- 格式3：纯数字ID格式（/{数字ID}）
- 格式4：纯字符ID格式（/{字符ID}）
并且所有这些格式都支持访问根路径（如/1001）和带路径的资源（如/1001/abc.html），大大提高了系统的灵活性和用户体验。




主要修改内容：

1. 1.
   添加了新的API接口 ：
   
   - POST /merchant/create/:merchantId - 创建商户
   - GET /merchant/:merchantId/domains - 获取商户domains.json
   - PUT /merchant/:merchantId/domains - 更新商户domains.json
   - POST /merchant/:merchantId/domains/upload - 上传商户domains.json
2. 2.
   创建商户功能 ( createMerchant 函数)：
   
   - 验证商户ID合法性（只允许字母、数字和下划线）
   - 检查商户是否已存在和模板目录是否存在
   - 创建商户的目录结构（html、static、config、data）
   - 从 www/merchant_template 模板目录复制文件
   - 替换所有文件内容中的 MERCHANT_ID 占位符为实际商户ID
   - 重命名index.html文件
   - 设置目录权限
3. 3.
   domains.json管理功能 ：
   
   - getMerchantDomains ：读取并返回指定商户的domains.json内容
   - updateMerchantDomains ：通过JSON请求体直接更新domains.json
   - uploadMerchantDomains ：上传新的JSON文件作为domains.json，并验证文件格式
4. 4.
   代码优化 ：
   
   - 添加了适当的错误处理和响应
   - 修复了上传JSON文件时的拼写错误
   - 代码通过 go build 编译检查，确保语法正确


   - 1.
修改nginx.conf配置文件

- 在保留现有HTTP功能的基础上，添加了完整的HTTPS服务器配置块
- 配置了SSL证书路径、TLS协议设置和安全优化配置
- 确保HTTPS配置包含与HTTP相同的商户路径处理规则
- 修复了静态资源路径配置问题，使CSS等文件能够正确加载
- 2.
修改docker-compose.yml文件

- 添加了443端口映射（8443:443）
- 添加了SSL证书目录挂载（./ssl:/etc/nginx/ssl:ro）
- 3.
创建SSL证书

- 创建了ssl目录
- 生成了自签名证书（fullchain.pem和privkey.pem）用于测试
- 4.
创建测试文件和目录

- 创建了www/merchant_1000/html和www/merchant_1000/static目录
- 创建了HTML测试页面和CSS样式文件
- 5.
功能测试

- ✅ 验证了HTTPS连接可以正常建立
- ✅ 验证了健康检查接口可以通过HTTPS访问
- ✅ 验证了HTML页面可以通过HTTPS访问
- ✅ 验证了静态资源（CSS文件）可以通过HTTPS访问


## 1. Nginx主配置文件更新
在 merchant-upload-service/nginx/nginx.conf 中添加了以下功能：

- 增强了主日志格式的注释说明
- 添加了专门的 merchant 日志格式，包含商户ID识别
- 配置了错误日志记录到 /var/log/nginx/error.log
- 添加了慢速请求日志配置，自动记录响应时间超过1秒的请求
- 为所有商户相关的location块添加了商户ID提取和特定日志记录
- 为上传服务代理配置添加了日志注释示例
## 2. 日志轮转配置
创建了 merchant-upload-service/nginx/logrotate.conf 文件，配置了自动日志轮转策略：

- 每天自动轮转一次
- 保留14天的历史日志
- 旧日志会自动压缩
- 轮转后自动通知Nginx重新打开日志文件
## 3. 使用说明文档
创建了 merchant-upload-service/nginx/README_LOG.md 文档，详细说明了：

- 日志文件结构和各日志文件的用途
- 不同日志格式的详细字段说明
- 日志轮转配置详情
- 查看和分析日志的常用命令
- 日志分析建议和注意事项
## 4. 验证和应用脚本
创建了 merchant-upload-service/nginx/verify_log_config.sh 脚本并设置了执行权限，该脚本可以：

- 检查并创建日志目录，设置正确的权限
- 验证Nginx配置的正确性
- 重启Nginx服务以应用变更
- 配置logrotate策略
- 提供常用日志查看命令的参考
## 使用方法
您可以运行以下命令来验证和应用日志配置变更：

```
sudo merchant-upload-service/nginx/verify_log_config.sh
```
通过这些配置，您现在可以更有效地监控、分析和管理商户服务的访问日志，特别是可以根据商户ID进行访问行为分析和问题排查。






## 实现思路

域名跳转的商户隔离实现：

1. 基于域名区分
   如果每个商户有独立的域名，则配置不同server块来实现：
   增加商户时，需要在nginx.conf中添加新的server块，配置商户的域名、根目录和首页文件。

   先决条件：
   需要有通配符域名，例如 *.merchant.com ，并配置好DNS解析。
   或者每个商户都有独立域名，并解析到这台服务器的IP地址。

   优点：
      每个商户有独立的域名，只需要提供对应的server块，泾渭分明；nginx重新加载配置即可

   缺点：
      每次新增商户都需要手动新增server块，nginx需要重新加载配置才能生效。

```
# 商户1配置
server {
    listen 80;
    server_name merchant1.com www.merchant1.com;  # 商户1的域名

    # 商户1的网站根目录
    root /var/www/merchant1;
    
    # 商户1的首页文件
    index merchant1-home.html index.html;
    
    # 其他通用配置
    location / {
        try_files $uri $uri/ =404;
    }
}

# 商户2配置
server {
    listen 80;
    server_name merchant2.com www.merchant2.com;  # 商户2的域名

    root /var/www/merchant2;
    index merchant2-index.html;  # 商户2的首页
    
    location / {
        try_files $uri $uri/ =404;
    }
}

# 更多商户以此类推...
    
```

2.  基于路径区分
    所有商户共享一个域名，通过路径来区分不同商户的请求。

    优点：
       每个商户有独立的路径，只需要提供对应的location块，泾渭分明；nginx重新加载配置即可

    缺点：
       每次新增商户都需要手动新增location块，nginx需要重新加载配置才能生效。

```
server {
    listen 80;
    server_name example.com;  # 共用域名
    root /var/www;  # 根目录

    # 商户1：example.com/merchant1
    location /merchant1 {
        # 商户1的首页文件
        index merchant1-home.html;
        # 可选：指定商户1的独立目录
        # root /var/www/merchant1;
    }

    # 商户2：example.com/merchant2
    location /merchant2 {
        index merchant2-main.html;
        # root /var/www/merchant2;
    }

    # 通用配置
    location / {
        try_files $uri $uri/ =404;
    }
}
     
```

3. 在方案2的基础上增加动态变量匹配，这种情况适合大量的商户
   具体实现：
   设计复杂的正则表达式，根据路径中的动态变量（如商户ID）来匹配不同的location块。
   例如：/merchant/([0-9a-zA-Z-]+) 可以匹配 /merchant/merchant123 中的 merchant123 作为商户ID。
   然后在location块中使用 $1 来引用这个动态变量，例如：
```
location ~ ^/merchant/([0-9a-zA-Z-]+)$ {
    index merchant-$1-home.html;
    # 可选：指定商户的独立目录
    # root /var/www/merchant/$1;
}
```

   缺点：设计一个好用的正则表达式很难
   优点：一劳永逸，新增商户不需要修改nginx配置

4. 商户按目录隔离，方便管理

5. 我选择的是方案3：
   - nginx 通用配置完成
   - 通用目录结构和跳转脚本完成
   - shell脚本创建商户完成
   - 通过API创建商户已经完成，正在测试