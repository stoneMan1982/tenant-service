#!/bin/bash
# 功能：创建商户目录（含html/static/config/data）、从模板复制内容、替换商户ID、重载Nginx

# 1. 检查参数是否正确
if [ $# -ne 1 ]; then
    echo "使用方法: $0 <商户ID>"
    echo "支持格式: 字母（大小写）、数字、下划线（例：A123、shop_B89、M_001）"
    exit 1
fi

MERCHANT_ID=$1
MERCHANT_DIR="merchant_${MERCHANT_ID}"  # 商户目录名（加前缀避免冲突）
BASE_DIR="../www/${MERCHANT_DIR}"          # 商户完整路径（当前目录下的www）
TEMPLATE_DIR="merchant_template"     # 独立的模板目录，与具体商户分离

# 2. 要求用户输入域名和端口号
read -p "请输入域名（如：example.com）: " DOMAIN
if [ -z "$DOMAIN" ]; then
    echo "错误: 域名不能为空！"
    exit 1
fi

read -p "请输入端口号（默认8443）: " PORT
if [ -z "$PORT" ]; then
    PORT=8443  # 默认端口号
fi

# 验证端口号是否为数字
if ! [[ $PORT =~ ^[0-9]+$ ]]; then
    echo "错误: 端口号必须为数字！"
    exit 1
fi

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

# 4. 检查模板目录是否存在
if [ ! -d "$TEMPLATE_DIR" ]; then
    echo "错误: 模板目录${TEMPLATE_DIR}不存在，请先确保模板存在！"
    exit 1
fi

# 5. 自动创建商户子目录（html/static/config/data）
echo "1/5：创建商户${MERCHANT_ID}目录结构..."
mkdir -p "$BASE_DIR"/{html,static,config,data}

# 6. 从模板复制文件到新商户目录
echo "2/5：从模板复制文件..."
# 复制html目录下的文件
if [ -d "$TEMPLATE_DIR/html" ]; then
    cp -R "$TEMPLATE_DIR/html/"* "$BASE_DIR/html/"
else
    echo "警告: 模板html目录不存在，跳过复制html文件"
fi

# 复制data目录下的文件
if [ -d "$TEMPLATE_DIR/data" ]; then
    cp -R "$TEMPLATE_DIR/data/"* "$BASE_DIR/data/"
else
    echo "警告: 模板data目录不存在，跳过复制data文件"
fi

# 复制static目录下的文件
if [ -d "$TEMPLATE_DIR/static" ]; then
    cp -R "$TEMPLATE_DIR/static/"* "$BASE_DIR/static/"
else
    echo "警告: 模板static目录不存在，跳过复制static文件"
fi

# 7. 替换文件内容中的商户ID、域名和端口号
echo "3/5：替换文件中的商户ID、域名和端口号..."
# 重命名index.html文件
template_index_file="$BASE_DIR/html/merchant_MERCHANT_ID_index.html"
if [ -f "$template_index_file" ]; then
    mv "$template_index_file" "$BASE_DIR/html/${MERCHANT_DIR}_index.html"
fi

# 在html目录中替换商户ID、域名和端口号
for file in "$BASE_DIR/html"/*; do
    if [ -f "$file" ]; then
        # 使用sed替换文件中的MERCHANT_ID为新的商户ID（适配macOS的sed语法）
        sed -i "" "s/MERCHANT_ID/$MERCHANT_ID/g" "$file"
        # 替换域名
        sed -i "" "s/你的域名/$DOMAIN/g" "$file"
        sed -i "" "s/localhost/$DOMAIN/g" "$file"
        # 替换可能的硬编码IP地址
        # sed -i "" "s/16\.163\.193\.74/$DOMAIN/g" "$file"
        # 替换端口号
        sed -i "" "s/YOUR_PORT/$PORT/g" "$file"
        sed -i "" "s/8080/$PORT/g" "$file"
        # 替换协议为https
        #sed -i "" "s/http:\/\//https:\/\//g" "$file"
    fi

done

# 在data目录中替换商户ID、域名和端口号
for file in "$BASE_DIR/data"/*; do
    if [ -f "$file" ]; then
        # 使用sed替换文件中的MERCHANT_ID为新的商户ID（适配macOS的sed语法）
        sed -i "" "s/MERCHANT_ID/$MERCHANT_ID/g" "$file"
        # 替换域名
        sed -i "" "s/你的域名/$DOMAIN/g" "$file"
        sed -i "" "s/localhost/$DOMAIN/g" "$file"
        # 替换端口号
        sed -i "" "s/YOUR_PORT/$PORT/g" "$file"
        sed -i "" "s/8080/$PORT/g" "$file"
        # 替换协议为https
        sed -i "" "s/http:\/\//https:\/\//g" "$file"
    fi

done

# 清理可能生成的空后缀备份文件
find "$BASE_DIR/html" -name "*''" -delete
find "$BASE_DIR/data" -name "*''" -delete

# 8. 设置权限（适配Nginx容器默认用户UID=101，避免权限拒绝）
echo "4/5：设置目录权限（Nginx可访问）..."
chmod -R 755 "$BASE_DIR"          # 所有者读写，其他只读+执行
chown -R 101:101 "$BASE_DIR"      # 匹配Nginx容器用户ID

# 9. 重载Nginx配置（无需重启容器，立即生效）
echo "5/5：重载Nginx配置..."
if docker compose exec -T nginx nginx -s reload > /dev/null 2>&1; then
    echo "Nginx配置重载成功！"
else
    echo "警告: Nginx重载失败，请检查merchants.conf配置！"
fi

# 输出最终结果
echo -e "\n✅ 商户${MERCHANT_ID}创建完成！"
echo "📁 目录结构："
ls -l "$BASE_DIR"  # 显示子目录（若有tree命令，可替换为tree "$BASE_DIR"）
echo -e "\n🌐 访问地址：https://${DOMAIN}:${PORT}/${MERCHANT_DIR}"  # 完整格式
echo "🌐 访问地址：https://${DOMAIN}:${PORT}/${MERCHANT_ID}"    # 简短格式（如果Nginx配置支持）