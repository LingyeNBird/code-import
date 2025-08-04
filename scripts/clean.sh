#!/bin/bash
ORG_NAME="xxx" # 替换为你需要清理的组织名称
API_TOKEN="xxx" # token所需权限 group-delete:rw,repo-delete:rw,group-resource:rw
API_URL="https://api.cnb.cool"

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
