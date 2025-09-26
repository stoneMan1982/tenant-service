#!/bin/bash

# Nginx日志配置验证脚本
# 此脚本用于检查日志目录权限、验证nginx配置并应用变更

# 设置颜色变量
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # 无颜色

# 日志目录路径
LOG_DIRS=("/var/log/nginx")

# 脚本开始
printf "${YELLOW}=== Nginx日志配置验证与应用 ===${NC}\n\n"

# 检查是否以root权限运行
if [ "$(id -u)" != "0" ]; then
   printf "${RED}错误：此脚本需要以root权限运行${NC}\n"
   printf "请使用: sudo %s\n" "$0"
   exit 1
fi

# 检查并创建日志目录
echo "检查日志目录..."
for dir in "${LOG_DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        echo "创建日志目录: $dir"
        mkdir -p "$dir"
        if [ $? -ne 0 ]; then
            printf "${RED}创建目录失败: $dir${NC}\n"
            exit 1
        fi
    fi
    
    # 检查并设置权限
    current_user=$(stat -c '%U' "$dir")
    current_group=$(stat -c '%G' "$dir")
    
    if [ "$current_user" != "nginx" ] || [ "$current_group" != "nginx" ]; then
        echo "设置目录权限: $dir 给 nginx:nginx"
        chown -R nginx:nginx "$dir"
        chmod -R 755 "$dir"
    fi
    
    echo "✓ 目录 $dir 检查完成"
done

# 验证Nginx配置
echo -e "\n验证Nginx配置..."
nginx -t
if [ $? -eq 0 ]; then
    printf "${GREEN}✓ Nginx配置验证通过${NC}\n"
else
    printf "${RED}✗ Nginx配置验证失败，请检查配置文件${NC}\n"
    exit 1
fi

# 询问是否重启Nginx
echo -e "\n配置验证通过，是否重启Nginx以应用变更？(y/n): "
read -r RESTART

if [[ "$RESTART" == [Yy]* ]]; then
    echo "重启Nginx服务..."
    systemctl restart nginx
    if [ $? -eq 0 ]; then
        printf "${GREEN}✓ Nginx服务已成功重启${NC}\n"
    else
        printf "${RED}✗ Nginx服务重启失败${NC}\n"
        exit 1
    fi
else
    echo "跳过Nginx重启，配置将在下次服务重启时生效"
fi

# 设置logrotate配置
echo -e "\n检查logrotate配置..."
LOGROTATE_CONFIG="nginx/logrotate.conf"
LOGROTATE_DEST="/etc/logrotate.d/nginx"

if [ -f "$LOGROTATE_CONFIG" ]; then
    echo "发现logrotate配置文件，是否复制到系统目录？(y/n): "
    read -r COPY_LOGROTATE
    
    if [[ "$COPY_LOGROTATE" == [Yy]* ]]; then
        cp "$LOGROTATE_CONFIG" "$LOGROTATE_DEST"
        if [ $? -eq 0 ]; then
            printf "${GREEN}✓ logrotate配置已复制到 $LOGROTATE_DEST${NC}\n"
        else
            printf "${RED}✗ 复制logrotate配置失败${NC}\n"
        fi
    fi
fi

# 显示日志查看命令
echo -e "\n${YELLOW}=== 日志查看命令 ===${NC}"
echo "实时查看访问日志: tail -f /var/log/nginx/access.log"
echo "实时查看商户访问日志: tail -f /var/log/nginx/merchant_access.log"
echo "实时查看错误日志: tail -f /var/log/nginx/error.log"
echo "实时查看慢速请求日志: tail -f /var/log/nginx/slow.log"

echo -e "\n${GREEN}✓ 日志配置验证脚本执行完成${NC}\n"