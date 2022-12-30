# import-tools-web

[English](../../README.md)

这个是导入 Jira 数据的前端库


## 🖥 兼容环境

- 现代浏览器。

| [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/edge/edge_48x48.png" alt="Edge" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Edge | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/firefox/firefox_48x48.png" alt="Firefox" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Firefox | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/chrome/chrome_48x48.png" alt="Chrome" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Chrome | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/safari/safari_48x48.png" alt="Safari" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Safari |
| --- | --- | --- | --- |
| >=88 | >=78 | >=87 | >=14 |


## 🚀 运行测试
1. 先安装依赖
```bash
npm install
```
2. 运行测试环境
```bash
npm run dev
```

## 📦 打包发布
> 修改后，执行打包指令，确保改动能打包正常
1. 先安装依赖
```bash
npm install
```
2. 执行编译
```bash
npm run build
```
3. 编译文件路径
```
/dist
```

## 🌍 国际化
国际化能力，请参考 [国际化文档]()。

## ⚙️ 路由管理
> 本项目路由已实现自动收集，具体实现可查看 [source code](./src//router/index.tsx)

使用的是 hash 路由，路径由 page 目录直接构成，举个例子。

### 例子
file path
```
1. /src/page/demo/index.page.tsx
2. /src/page/demo/home/index.page.tsx
3. /src/page/demo/test_demo/index.page.tsx
```
url path
```
https://XXX.com/#/page/demo
https://XXX.com/#/page/demo/home
https://XXX.com/#/page/demo/test_demo
```
通过文件后缀 .page 收集为对应目录路由页面

### 义重定向设置
自定义重定向，可以直接在 [redirect route config](./src//router/routes.ts) 配置

## 🔔 注意事项
本项目集成以下框架，可直接按照标准用法使用

| 框架名称 | 版本 |
| - | - |
| Vite | v4 |
| React | v18 |
| styled-components | 5.3.6 |
| Ant Design | v5 |
| react-router-dom | 6.4.3 |
| Lodash | 4.17.21 |
| ahooks | 3.7.2 |
| Typescript | 4.6.4 |
| TailwindCss | v3 |

其他细节版本请查看 [package.json](./package.json)
