# 🌍 Internationalization Guidelines

> 🎉 Contributing to translation is the best start to participate in open source projects, and everyone is welcome to participate.

## add new language
Next, "Chinese" has been added as an example for guidance.

1. Open the internationalization configuration file
```
src/i18n.ts
```

2. Write a new enum type
```
enum LangEnums {
  en = 'en',
  zh = 'zh', // new enum
}
```

3. Import the corresponding language pack of "Ant Design" at the corresponding position at the top. If "Ant Design" does not support it, you can consider using English
```
import zhCN from 'antd/locale/zh_CN';
```

4. Set the mapping for the "Ant Design" language pack
```
const AntDesignLang = {
  [LangEnums.en]: enUS,
  [LangEnums.zh]: zhCN, // new language
};
```

5. Settings page language selection configuration
```
// Affects language setting sorting
export const LANGS = [
  LangEnums.zh, // new language
  LangEnums.en,
  ];
```

6. Go to the "src/lang" directory to create a folder "zh" for translations.
Import your translation module in "rc/lang/index.ts"
```
import zh from './zh';
import en from './en';

export {
  zh, // new language
  en
};
```

7. Go back to "src/i18n.ts", import translation copy
```
import { zh, en } from '@/lang';
```

8. Find the "resources" configuration and load the translation text
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

9. Supplement the translation of various language switching, in the directory of various language packs "src/lang/**/language.ts"
```
const language = {
  'zh': '中文', // new text
  'en': 'English',
}

export default language;
```
