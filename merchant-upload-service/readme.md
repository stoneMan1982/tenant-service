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