# GLab-MR

Terminal UI tool was written to help to monitor GitLab merge requests.

## How use

Create config `~/.glab_mr.json`
```json
{
  "gitlab_base_url": "string",
  "gitlab_token": "string",
  "gitlab_username": "string",
  "projects": {
    "{{project_name_1}}": {{project_id_1}},
    "{{project_name_2}}": {{project_id_2}},
  }
}
```

Hotkeys:
* j, up - Scroll down
* k, Down - Scroll up
* q - Exit
* o - Open MR in your browser


