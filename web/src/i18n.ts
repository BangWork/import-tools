import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import { includes } from 'lodash-es';
import LanguageDetector from 'i18next-browser-languagedetector';

/** Ant Design language pack */
import enUS from 'antd/locale/en_US';
import zhCN from 'antd/locale/zh_CN';

import { zh, en } from '@/lang';

enum LangEnums {
  zh = 'zh',
  en = 'en',
}

// langs map of AntDesign
const AntDesignLang = {
  [LangEnums.zh]: zhCN,
  [LangEnums.en]: enUS,
};


// language setting config, can effect the view at topLeft;
export const LANGS = [LangEnums.zh, LangEnums.en];

i18n
  // detect user language
  // learn more: https://github.com/i18next/i18next-browser-languageDetector
  .use(LanguageDetector)
  // pass the i18n instance to react-i18next.
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    debug: true,
    fallbackLng: LANGS,
    interpolation: {
      escapeValue: false, // not needed for react as it escapes by default
    },
    resources: {
      zh: {
        translation: zh,
      },
      en: {
        translation: en,
      },
    },
  });

export const getCurrentLang = (currentI18n = i18n) => {
  return includes(LANGS, currentI18n.resolvedLanguage) ? currentI18n.resolvedLanguage : LANGS[0];
};

export const getAntDesignLang = (lang: LangEnums | string) => {
  return AntDesignLang[lang];
};

export default i18n;
