#!/bin/bash
# 功能：将没有下划线的商户目录（如merchant1001）重命名为标准格式（merchant_1001）

BASE_DIR="../www"

# 检查www目录是否存在
if [ ! -d "$BASE_DIR" ]; then
    echo "错误: www目录不存在！"
    exit 1
fi

# 查找所有没有下划线但以merchant开头的目录
MERCHANT_DIRS=$(find "$BASE_DIR" -maxdepth 1 -type d -name "merchant[0-9]*" -not -name "merchant_*")

if [ -z "$MERCHANT_DIRS" ]; then
    echo "没有需要重命名的商户目录。"
    exit 0
fi

echo "发现以下需要重命名的商户目录："
echo "$MERCHANT_DIRS"
echo

# 询问用户是否继续
read -p "是否继续重命名操作？(y/n): " CONFIRM
if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo "操作已取消。"
    exit 0
fi

# 逐个重命名目录
for DIR in $MERCHANT_DIRS; do
    # 提取商户ID（去除merchant前缀）
    MERCHANT_ID=${DIR##*/merchant}
    
    # 构建新的目录名（带下划线）
    NEW_DIR="$BASE_DIR/merchant_${MERCHANT_ID}"
    
    echo "正在重命名 $DIR -> $NEW_DIR"
    
    # 执行重命名
    if mv "$DIR" "$NEW_DIR"; then
        echo "重命名成功！"
        
        # 查找并重命名index.html文件（如果存在）
        OLD_INDEX_FILE="$NEW_DIR/html/merchant${MERCHANT_ID}_index.html"
        NEW_INDEX_FILE="$NEW_DIR/html/merchant_${MERCHANT_ID}_index.html"
        
        if [ -f "$OLD_INDEX_FILE" ]; then
            echo "  正在重命名index.html文件..."
            mv "$OLD_INDEX_FILE" "$NEW_INDEX_FILE"
        fi
        
        # 替换文件内容中的商户ID引用（从merchant1001改为merchant_1001）
        echo "  正在更新文件内容中的引用..."
        
        # 处理HTML文件
        for HTML_FILE in "$NEW_DIR/html"/*.html; do
            if [ -f "$HTML_FILE" ]; then
                sed -i "" "s/merchant${MERCHANT_ID}/merchant_${MERCHANT_ID}/g" "$HTML_FILE"
                sed -i "" "s/MERCHANT_ID/${MERCHANT_ID}/g" "$HTML_FILE"
            fi
        done
        
        # 处理data目录下的文件
        for DATA_FILE in "$NEW_DIR/data"/*; do
            if [ -f "$DATA_FILE" ]; then
                sed -i "" "s/merchant${MERCHANT_ID}/merchant_${MERCHANT_ID}/g" "$DATA_FILE"
                sed -i "" "s/MERCHANT_ID/${MERCHANT_ID}/g" "$DATA_FILE"
            fi
        done
        
        echo "  商户${MERCHANT_ID}的所有引用已更新！"
    else
        echo "错误: 重命名失败！"
    fi
    
echo

done

echo "所有操作完成！"
echo "请运行以下命令重载Nginx配置："
echo "docker compose exec -T nginx nginx -s reload"