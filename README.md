# file-matrix

## 流程
1. client - server
- reg:
  - req: name key
  - resp: operate
- push:
  - req: name key file
2. browser - server
- index:
  - resp: html
- names:
  - req: key
  - resp: name list
- dir:
  - req: name key user dir
  - resp: file list
- file:
  - req: name key user file
  - resp: file