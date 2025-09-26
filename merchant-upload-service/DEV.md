

## 实现思路

域名跳转的商户隔离实现：

1. 基于域名区分
   如果每个商户有独立的域名，则配置不同server块来实现：
   增加商户时，需要在nginx.conf中添加新的server块，配置商户的域名、根目录

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
优点：一劳永逸，新增商户不需要修改nginx配置，不用重启，不用重新加载

4. 商户按目录隔离，方便管理

5. 我选择的是方案3：
    - nginx 通用配置完成
    - 通用目录结构和跳转脚本完成
    - 支持任意数量的跳转池子完成
    - shell脚本创建商户完成
    - 商户名支持字母数字加下划线的组合
    - 通过API创建商户已经完成，正在测试
 