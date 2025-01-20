#!/bin/bash

# 打印当前参数
echo "当前参数: $@"

# 提示用户输入 yes 以继续
read -p "请输入 'yes' 以继续执行: " user_input

# 检查用户输入是否为 yes
if [[ $user_input == "yes" ]]; then
    echo "用户确认，继续执行..."
    # 在这里添加你需要执行的命令
    echo "执行命令: $@"
    # 示例：执行传入的命令
    "$@"
else
    echo "用户取消，退出脚本。"
    exit 1
fi