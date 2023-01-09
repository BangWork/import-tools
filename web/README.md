# import-tools-web

[简体中文](./documents/zh-Hans/README.md)

This is the front-end library for importing Jira data

## 🖥 Environment Support

- Modern browsers

| [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/edge/edge_48x48.png" alt="Edge" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Edge | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/firefox/firefox_48x48.png" alt="Firefox" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Firefox | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/chrome/chrome_48x48.png" alt="Chrome" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Chrome | [<img src="https://raw.githubusercontent.com/alrra/browser-logos/master/src/safari/safari_48x48.png" alt="Safari" width="24px" height="24px" />](http://godban.github.io/browsers-support-badges/)<br>Safari |
| --- | --- | --- | --- |
| >=88 | >=78 | >=87 | >=14 |

## 🚀 Start
1. install dependence
```bash
npm install
```
2. build local test
```bash
npm run dev
```

## 📦 Deploy
> After the modification, execute the packaging command to ensure that the changes can be packaged normally
1. install dependence
```bash
npm install
```
2. build
```bash
npm run build
```
3. export path
```
/dist
```

## 🌍 Internationalization
> 👏 Contributing to translation is the best start to participate in open source projects, and everyone is welcome to participate.

How to achieve internationalization, see [i18n](./documents/i18n/i18n-en.md).

## ⚙️ Routing management
> ✨ The route of this project has been automatically collected, and the specific implementation can be viewed [source code](./src/router/index.tsx)

use hash route,route path belong project file path to transform.Just like this example.

### Example
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
Collect routing pages for the corresponding directory through the file suffix .page

### redirect path setting
you can in this [redirect route config](./src//router/routes.ts) file to set path mapping.

## 🔔 Caveats
This project integrates the following frameworks, which can be used directly in accordance with standard usage.

| name | version |
| - | - |
| Vite | v4 |
| React | v18 |
| styled-components | 5.3.6 |
| Ant Design | v5 |
| react-router-dom | 6.4.3 |
| lodash-es | 4.17.21 |
| ahooks | 3.7.2 |
| typescript| 4.6.4 |
| tailwindCss | v3 |

For other detailed versions, please check [package.json](./package.json)
