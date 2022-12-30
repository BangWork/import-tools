# 🌍 国际化指南

> 🎉 贡献翻译是参与开源项目最好的开端，欢迎各位参与。

## 添加新语言
接下来已添加「中文」为例子，进行指引。

1. 打开国际化配置文件
```
src/i18n.ts
```

2. 编写新的枚举类型
```
enum LangEnums {
  en = 'en',
  zh = 'zh', // new enum
}
```

3. 在顶部对应位置导入「Ant Design」的对应语言包,如果「Ant Design」不支持，可以考虑使用英语
```
import zhCN from 'antd/locale/zh_CN';
```

4. 设置「Ant Design」语言包的映射
```
const AntDesignLang = {
  [LangEnums.en]: enUS,
  [LangEnums.zh]: zhCN, // new language
};
```

5. 设置页面语言选择配置
```
// 会影响语言设置排序
export const LANGS = [
  LangEnums.zh, // new language
  LangEnums.en,
  ];
```

6. 前往「src/lang」目录创建翻译文案的文件夹「zh」。
在 「rc/lang/index.ts」导入你的翻译模块
```
import zh from './zh';
import en from './en';

export {
  zh, // new language
  en
};
```

7. 回到「src/i18n.ts」, 导入翻译文案
```
import { zh, en } from '@/lang';
```

8. 找到「resources」配置，载入翻译文案
```
resources: {
  en: {
    translation: en,
  },

  // new config
  zh: {
    translation: zh,
  },
},
```

9. 补充各种语言切换的翻译，在各种语言包「src/lang/**/language.ts」目录下
```
const language = {
  'zh': '中文', // new text
  'en': 'English',
}

export default language;
```
