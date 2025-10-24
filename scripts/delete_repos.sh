#!/bin/bash
# 功能：删除指定组织下所有仓库，慎用！！！
# 请先在组织设置-组织管控-危险操作-允许通过 Open API 删除组织下资源
ORG_NAME="xxx" # 替换为你需要清理的组织名称
API_TOKEN="xxx" # token所需权限 repo-delete:rw,group-resource:r
API_URL="https://api.cnb.cool"

# 检查并安装 jq 依赖
function check_and_install_jq() {
    if command -v jq &> /dev/null; then
        echo "jq is already installed"
        return 0
    fi
    
    echo "jq is not installed. Installing..."
    
    # 检测操作系统类型
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install jq
        else
            echo "Error: Homebrew is not installed. Please install Homebrew first or install jq manually."
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            case "$ID" in
                ubuntu|debian)
                    sudo apt-get update && sudo apt-get install -y jq
                    ;;
                centos|rhel|fedora)
                    if command -v dnf &> /dev/null; then
                        sudo dnf install -y jq
                    else
                        sudo yum install -y jq
                    fi
                    ;;
                alpine)
                    sudo apk add --no-cache jq
                    ;;
                arch|manjaro)
                    sudo pacman -S --noconfirm jq
                    ;;
                opensuse*|sles)
                    sudo zypper install -y jq
                    ;;
                *)
                    echo "Error: Unsupported Linux distribution: $ID"
                    echo "Please install jq manually: https://stedolan.github.io/jq/download/"
                    exit 1
                    ;;
            esac
        else
            echo "Error: Cannot detect Linux distribution"
            echo "Please install jq manually: https://stedolan.github.io/jq/download/"
            exit 1
        fi
    else
        echo "Error: Unsupported operating system: $OSTYPE"
        echo "Please install jq manually: https://stedolan.github.io/jq/download/"
        exit 1
    fi
    
    # 验证安装是否成功
    if command -v jq &> /dev/null; then
        echo "jq installed successfully"
    else
        echo "Error: jq installation failed"
        exit 1
    fi
}

# 执行 jq 依赖检查
check_and_install_jq

function delete_repos() {
    repos=$(curl -s -H "accept: application/json" -H "Authorization: Bearer ${API_TOKEN}" "${API_URL}/${ORG_NAME}/-/repos?page_size=50")
    echo "${repos}" | jq -r '.[].path' | while read repo; do
        echo "Deleting repository: ${repo}"
        curl -X DELETE -s -H "Authorization: Bearer ${API_TOKEN}" "${API_URL}/${repo}"
    done
}



count=0
max_attempts=50
while [ $count -lt $max_attempts ]; do
    delete_repos
    count=$((count + 1))
    echo "-------------------- $count/$max_attempts"
done
