#!/bin/bash
# 功能：修复已存在商户的domains.json文件中的域名和端口

# 检查参数是否正确
if [ $# -ne 3 ]; then
    echo "使用方法: $0 <商户ID> <域名> <端口号>"
    exit 1
fi

MERCHANT_ID=$1
DOMAIN=$2
PORT=$3
MERCHANT_DIR="merchant_${MERCHANT_ID}"
BASE_DIR="../www/${MERCHANT_DIR}"

# 验证商户ID合法性
if ! [[ $MERCHANT_ID =~ ^[a-zA-Z0-9_]+$ ]]; then
    echo "错误: 商户ID仅支持字母、数字、下划线！"
    exit 1
fi

# 验证端口号是否为数字
if ! [[ $PORT =~ ^[0-9]+$ ]]; then
    echo "错误: 端口号必须为数字！"
    exit 1
fi

# 检查商户目录是否存在
if [ ! -d "$BASE_DIR" ]; then
    echo "错误: 商户${MERCHANT_ID}不存在（目录：${BASE_DIR}）！"
    exit 1
fi

# 检查domains.json文件是否存在
DOMAINS_FILE="$BASE_DIR/data/domains.json"
if [ ! -f "$DOMAINS_FILE" ]; then
    echo "错误: domains.json文件不存在！"
    exit 1
fi

# 备份文件
BACKUP_FILE="${DOMAINS_FILE}.bak"
cp "$DOMAINS_FILE" "$BACKUP_FILE"
echo "已备份domains.json到${BACKUP_FILE}"

# 替换domains.json文件中的域名和端口
sed -i "" "s/localhost/$DOMAIN/g" "$DOMAINS_FILE"
sed -i "" "s/8080/$PORT/g" "$DOMAINS_FILE"
sed -i "" "s/http:\/\//https:\/\//g" "$DOMAINS_FILE"
echo "已成功替换商户${MERCHANT_ID}的domains.json文件中的域名和端口！"

# 处理HTML文件中的硬编码IP地址和域名
HTML_DIR="$BASE_DIR/html"
for html_file in "$HTML_DIR"/*.html; do
    if [ -f "$html_file" ]; then
        # 备份HTML文件
        html_backup="${html_file}.bak"
        cp "$html_file" "$html_backup"
        
        # 替换硬编码IP地址
        sed -i "" "s/16\.163\.193\.74/$DOMAIN/g" "$html_file"
        # 替换localhost
        sed -i "" "s/localhost/$DOMAIN/g" "$html_file"
        # 替换端口号
        sed -i "" "s/8080/$PORT/g" "$html_file"
        # 替换协议
        sed -i "" "s/http:\/\//https:\/\//g" "$html_file"
        
        echo "已处理HTML文件: $(basename "$html_file")"
    fi
done

# 显示替换后的内容
echo -e "\n替换后的文件内容预览："
head -n 20 "$DOMAINS_FILE"