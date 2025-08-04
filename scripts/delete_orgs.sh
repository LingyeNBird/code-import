#!/bin/bash
ORG_NAME="jacknic/platform"
API_TOKEN="xxx"
API_URL="https://api.cnb.cool"

function delete_repos() {
    local org_name=$1
    # 先获取删除子组织
    repos=$(curl -s -H "accept: application/json" -H "Authorization: Bearer ${API_TOKEN}" "${API_URL}/user/groups/${org_name}?page_size=50")
    echo "${repos}" | jq -r '.[].path' | while read repo; do
        delete_repos "$repo"
        echo "Deleting repository: ${repo}"
        curl -X DELETE -s -H "Authorization: Bearer ${API_TOKEN}" "${API_URL}/${repo}"
    done
    curl -X DELETE -s -H "Authorization: Bearer ${API_TOKEN}" "${API_URL}/${org_name}"
}

count=0
max_attempts=1
while [ $count -lt $max_attempts ]; do
    delete_repos "$ORG_NAME"
    count=$((count + 1))
    echo "-------------------- $count/$max_attempts"
done