#!/bin/bash
# 功能：创建商户目录（含html/static/config/data）、生成默认首页、重载Nginx

# 1. 检查参数是否正确
if [ $# -ne 1 ]; then
    echo "使用方法: $0 <商户ID>"
    echo "支持格式: 字母（大小写）、数字、下划线（例：A123、shop_B89、M_001）"
    exit 1
fi

MERCHANT_ID=$1
MERCHANT_DIR="merchant_${MERCHANT_ID}"  # 商户目录名（加前缀避免冲突）
BASE_DIR="../www/${MERCHANT_DIR}"       # 商户完整路径

# 2. 验证商户ID合法性（禁止特殊字符，防路径注入）
if ! [[ $MERCHANT_ID =~ ^[a-zA-Z0-9_]+$ ]]; then
    echo "错误: 商户ID仅支持字母、数字、下划线，不允许特殊字符！"
    exit 1
fi

# 3. 检查商户是否已存在
if [ -d "$BASE_DIR" ]; then
    echo "警告: 商户${MERCHANT_ID}已存在（目录：${BASE_DIR}），无需重复创建！"
    exit 0
fi

# 4. 自动创建商户子目录（html/static/config/data）
echo "1/4：创建商户${MERCHANT_ID}目录结构..."
mkdir -p "$BASE_DIR"/{html,static,config,data}

# 5. 生成商户默认首页（放在html目录下，与Nginx配置对应）
INDEX_FILE="$BASE_DIR/html/${MERCHANT_DIR}_index.html"
echo "2/4：生成默认首页：${INDEX_FILE}"
cat > "$INDEX_FILE" << EOF
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>商户${MERCHANT_ID}首页</title>
    <style>body{font-family:Arial,sans-serif;text-align:center;padding:5rem;}</style>
</head>
<body>
    <h1>欢迎访问商户 ${MERCHANT_ID}</h1>
    <p>静态资源路径：/static/（图片、CSS等）</p>
    <p>数据文件路径：/data/（配置、日志等）</p>
</body>
</html>
EOF

# 6. 设置权限（适配Nginx容器默认用户UID=101，避免权限拒绝）
echo "3/4：设置目录权限（Nginx可访问）..."
chmod -R 755 "$BASE_DIR"          # 所有者读写，其他只读+执行
chown -R 101:101 "$BASE_DIR"      # 匹配Nginx容器用户ID

# 7. 重载Nginx配置（无需重启容器，立即生效）
echo "4/4：重载Nginx配置..."
if docker compose exec -T nginx nginx -s reload > /dev/null 2>&1; then
    echo "Nginx配置重载成功！"
else
    echo "警告: Nginx重载失败，请检查merchants.conf配置！"
fi

# 输出最终结果
echo -e "\n✅ 商户${MERCHANT_ID}创建完成！"
echo "📁 目录结构："
ls -l "$BASE_DIR"  # 显示子目录（若有tree命令，可替换为tree "$BASE_DIR"）
echo -e "\n🌐 访问地址：https://你的域名/${MERCHANT_DIR}"
echo "🖼️  静态资源示例：https://你的域名/${MERCHANT_DIR}/static/logo.png"