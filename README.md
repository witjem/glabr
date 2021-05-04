# Glabr

Terminal UI tool was written to help to monitor GitLab merge requests.

## How use

Create config `~/.glabr.json`
```js
{
  "gitlab_base_url": "https://gitlab.com",
  "gitlab_token": "your_token",
  "gitlab_username": "your_username",
  "projects": {
    //project alias, and project id
    "glb": 123,
    "ain": 241
  }
}
```

Hotkeys:
* k, ↑ - Scroll line up
* j, ↓ - Scroll line down
* q - Exit
* o - Open MR in your browser

## Screenshots

![alt text](./doc/img/screenshot.png)
