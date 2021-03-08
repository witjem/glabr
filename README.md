# GLab-MR

Terminal UI tool was written to help to monitor GitLab merge requests.

## How use

Create config `~/.glab_mr.json`
```js
{
  "gitlab_base_url": "https://gitlab.com",
  "gitlab_token": "your_token",
  "gitlab_username": "your_username",
  "projects": {
    //project alias, and project id
    "gtw": 123,
    "mqq": 241
  }
}
```

Hotkeys:
* j, up - Scroll down
* k, Down - Scroll up
* q - Exit
* o - Open MR in your browser

