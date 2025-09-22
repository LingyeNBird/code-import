## Usage Steps
⚠️ Note: This method will start a pipeline in the repository to execute migration, logs will be public and contain `repository names`.

1. Fork this repository to any location under your organization
![img.png](../img/web_trigger_1.png)

2. Click `Execute` button above repository, select corresponding `source platform`, fill in configuration parameters as prompted.
![img.png](../img/web_trigger_2.png)

3. Click `orange button at bottom left` to trigger custom event and start migration task

4. Click hyperlink in popup to view task execution logs
![img.png](../img/web_trigger_4.png)

5. Click code-import step to view detailed logs
![img_2.png](../img/web_trigger_5.png)

6. Wait for migration task to complete, check final logs for `Failed migrations` and `Skipped migrations` counts (should be 0), verify all repositories migrated successfully
![img_8.png](../img/web_trigger_6.png)