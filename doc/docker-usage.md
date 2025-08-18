# 在Docker上使用

## 注意事项
1. 示例参数值为 `xxx`，需替换为自己执行参数值。[详细参数介绍](parameters.md)  


2. **[云原生开发](https://docs.cnb.cool/zh/workspaces/intro.html) 环境默认集成 docker cli，且内网运行更快速** 。


3. Windows Powershell 环境 换行符 ` \ ` 请替换为反引号 ` ， $(pwd) 替换为 ${PWD}

## 使用方法
1. 在 CNB 创建 1 个空仓库


2. 点击 `云原生开发`
![img_4.png](../img/docker_usage_1.png)


3. 使用 `WebIDE` 打开


4. 根据实际情况，复制下方迁移命令到终端，`xxx` 记得替换，然后回车执行
![img_5.png](../img/docker_usage_2.png)


5. 等待迁移完成，确认最终迁移结果


## 示例
* <details>
    <summary> 从 CODING 迁移 </summary>

    迁移之后的效果：原 CODING 项目会在 CNB 中创建一个同名的子组织，并将原项目下的仓库迁移至该子组织下面
    
    ### 迁移团队下所有仓库
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
    
    ### 迁移指定项目仓库
    PLUGIN_SOURCE_PROJECT 字段根据需要自行替换，详见[参数介绍](parameters.md)
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_SOURCE_PROJECT="project1,project2" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -e PLUGIN_MIGRATE_TYPE="project" \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
    
    ### 迁移指定仓库
    PLUGIN_SOURCE_REPO 字段根据需要自行替换，详见[参数介绍](parameters.md)
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_SOURCE_REPO="<TEAM-NAME>/<PROJECT-NAME>/<REPO-NAME>,test-team/project1/repoA,test-team/project2/repoB" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -e PLUGIN_MIGRATE_TYPE="repo" \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
  </details>

* <details>
    <summary> 从 GitHub 迁移 </summary>

    迁移之后的效果：原 GitHub 账号下有权限的所有组织，会在 CNB 中创建同名的子组织，并将原组织下的仓库迁移至该子组织下面
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_SOURCE_URL="https://github.com" \
      -e PLUGIN_SOURCE_PLATFORM="github" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
  </details>

* <details>
    <summary> 从 GitLab 迁移 </summary>
  
    迁移之后的效果：原 GitLab 账号下有权限的所有 group，会在 CNB 中创建同名的子组织，并将原group下的仓库迁移至该子组织下面（如果 gitlab 是多级的group，迁移至 CNB 子组织和仓库仍会保留原有的多层级结构）
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_SOURCE_URL="https://gitlab.com" \
      -e PLUGIN_SOURCE_PLATFORM="gitlab" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
  </details>

* <details>
    <summary> 从 Gitee 迁移 </summary>

    迁移之后的效果：原 Gitee 账号下有权限的所有组织，会在 CNB 中创建同名的子组织，并将原组织下的仓库迁移至该子组织下面（如果 Gitee 是多级的组织/仓库组，迁移至 CNB 子组织和仓库仍会保留原有的多层级结构）
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_SOURCE_URL="https://gitee.com" \
      -e PLUGIN_SOURCE_PLATFORM="gitee" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
  </details>

* <details>
    <summary> 从 Codeup (云效)迁移 </summary>

    迁移之后的效果：原云效账号下有权限的所有组织，会在 CNB 中创建同名的子组织，并将原仓库组下的仓库迁移至该子组织下面（如果云效是多级的仓库组，迁移至 CNB 子组织和仓库仍会保留原有的多层级结构）
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_SOURCE_PLATFORM="aliyun" \
      -e PLUGIN_SOURCE_ORGANIZATIONID="xxx" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
  </details>

* <details>
    <summary> 从通用第三方代码平台迁移 </summary>

    PLUGIN_SOURCE_REPO 字段中，group字段会映射为子组织，如果 cnb 的根组织下，没有该命名的子组织将会自动创建，如果有该命名的子组织会将仓库创建在已有的同名子组织下面
    ### http协议
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_USERNAME="xxx"  \
      -e PLUGIN_SOURCE_PASSWORD="xxx"  \
      -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
      -e PLUGIN_SOURCE_URL="https://common.example.com" \
      -e PLUGIN_SOURCE_PLATFORM="common" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
    
    ### ssh协议
    ⚠️使用ssh协议时请在当前工作目录确保有对应的私钥文件，文件名固定为`ssh.key`
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_REPO="group1/repo1,group1/repo2,group2/repo3" \
      -e PLUGIN_SOURCE_URL="https://common.example.com" \
      -e PLUGIN_SOURCE_PLATFORM="common" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -e GIT_SSH_COMMAND='ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null' \  
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
  </details>

* <details>
    <summary> 从工蜂迁移 </summary>

  ```shell
    docker run --rm  \
    -e PLUGIN_SOURCE_TOKEN="xxx"  \
    -e PLUGIN_SOURCE_URL="https://git.woa.com" \
    -e PLUGIN_SOURCE_PLATFORM="gongfeng" \
    -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
    -e PLUGIN_CNB_TOKEN="xxx"  \
    -v $(pwd):$(pwd) -w $(pwd) \
    cnbcool/code-import
  ```
  </details>

* <details>
    <summary> 从 CNB 迁移 </summary>

    迁移指定根组织下所有仓库
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_GROUP="xxx" \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_SOURCE_URL="https://cnb.example1.com" \
      -e PLUGIN_SOURCE_PLATFORM="cnb" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -e PLUGIN_CNB_URL="https://cnb.example2.com" \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
  </details>

* <details>
    <summary> 只迁移部分仓库（仓库选择功能） </summary>
    
    首次运行，生成仓库列表文件 `repo-path.txt`
    
    这里以CODING为例，其他平台只需在原有迁移命令基础上增加`-e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \`参数即可：
    
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
    
    **首次运行后，工具会在当前目录生成 `repo-path.txt`，请手动编辑该文件，仅保留需要迁移的仓库路径。**
    
    编辑完成后，再次运行同样的命令即可只迁移你选择的仓库：
    
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx"  \
      -e PLUGIN_MIGRATE_ALLOW_SELECT_REPOS="true" \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```

    **如需重新选择仓库，只需删除 `repo-path.txt`，重新运行上述命令即可。**

  </details>

* <details>
    <summary> 仅下载仓库模式 </summary>

    如果只需要将仓库下载到本地而不推送到 CNB 平台，可以使用仅下载模式。该模式下会跳过所有 CNB 相关操作，仅执行仓库克隆，无需提供 CNB 相关配置信息。
    
    这里以 CODING 为例，其他平台只需在原有迁移命令基础上增加 `-e PLUGIN_MIGRATE_DOWNLOAD_ONLY="true"` 参数即可：
    
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_TOKEN="xxx"  \
      -e PLUGIN_MIGRATE_DOWNLOAD_ONLY="true" \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```
    **仅下载模式特点：**
    - 仅执行仓库克隆操作，不推送到 CNB 平台
    - 无需提供 CNB 相关配置信息（如 CNB_TOKEN、CNB_ROOT_ORGANIZATION 等）
    - 跳过所有 CNB 相关操作（如创建子组织、创建仓库等）
    - 下载完成后保留工作目录（格式：`source_git_dir_时间戳`）

 </details>

* <details>
    <summary> 从本地仓库迁移（Local）</summary>

    本地模式会扫描工作目录下 `source_git_dir` 目录中的所有裸仓库，并将它们迁移至 CNB。裸仓库需满足目录下包含 `HEAD` 文件与 `objects/` 目录的结构。
    
    ### 准备本地裸仓库
    - 将待迁移的裸仓库放入当前工作目录的 `source_git_dir/` 下，支持多级子目录；
    - 例如：
      - `source_git_dir/group1/repo1/`
      - `source_git_dir/group2/sub/repo2/`
    
    ### 执行迁移
    ```shell
    docker run --rm  \
      -e PLUGIN_SOURCE_PLATFORM="local" \
      -e PLUGIN_CNB_ROOT_ORGANIZATION="xxx" \
      -e PLUGIN_CNB_TOKEN="xxx" \
      -v $(pwd):$(pwd) -w $(pwd) \
      cnbcool/code-import
    ```

    ### 说明
    - 本地模式不需要配置 `PLUGIN_SOURCE_URL`、`PLUGIN_SOURCE_TOKEN`；
    - 裸仓库的相对路径会映射为 CNB 的子组织层级，例如 `source_git_dir/group1/repo1` 最终路径为 `/<CNB根组织>/group1/repo1`；

  </details>


## 迁移完成后，增量更新原平台最新内容
清空原工作目录下的 successful.log

效果：重新同步原平台的所有仓库，已迁移至 CNB 的仓库，如在原平台有更新，会将内容增量同步至 CNB 平台。 