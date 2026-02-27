## 使用步骤


1. Fork 本仓库到你的组织下任意位置，迁移任务将把源平台仓库批量迁移至当前仓库所属的根组织名下。
![img.png](../img/web_trigger_1.png)


2. 点击仓库上方`CODING 迁移至 CNB`按钮，如果是其他平台，点击按钮旁边`...`，选择对应`源平台`，根据提示填写必填参数。  

   如果选择新增的 `同步Sealantern` 选项，只需要输入两个参数：`PLUGIN_SOURCE_TOKEN`（GitHub Token）和 `PLUGIN_CNB_TOKEN`（CNB Token），其余关键迁移参数已在流水线里硬编码。

    Tips:有提供默认参数的一般可以直接使用默认参数，不需要修改。

    >从 CODING 迁移至 CNB，只需填写 源平台 Token、CNB平台 Token，如未提前准备，请参考[迁移前准备](./ready.md)
![img.png](../img/web_trigger_2.png)

3. 点击`左下方 橙色 按钮`，启动自定义事件，开始运行迁移任务


4. 点击弹窗中的超链接，查看任务运行日志
![img.png](../img/web_trigger_4.png)


5. 耐心等待迁移任务执行完成，点击 code-import 这一步，查看日志最终输出结果，检查 `迁移失败` 和 `忽略迁移` 数量是否为 0，并在 CNB 侧确认仓库是否全部完成迁移
![img_2.png](../img/web_trigger_5.png)
![img_8.png](../img/web_trigger_6.png)
