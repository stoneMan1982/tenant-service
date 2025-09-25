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