import React, { memo, useEffect } from 'react';
import { createRoot } from 'react-dom/client';
import { HashRouter } from 'react-router-dom';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
import { Layout, ConfigProvider } from 'antd';

import AutoRouter from '@/router';
import { redirectRoutes } from '@/router/routes';

import Header from '@/components/header';
import Loading from '@/components/loading';
import NoMatch from '@/components/no_match';
import GlobalStyle from './global';

import './index.css';
import { getCurrentLang, getAntDesignLang } from '@/i18n';

const LayoutBox = styled(Layout)`
  background: #fff;
  height: 100%;
`;

const Content = styled.div`
  min-height: calc(100% - 104px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 50px;
  overflow-y: auto;
`;

/** default layout routes */
const NormalLayoutRoutes = memo(() => (
  <LayoutBox>
    <Header />
    <Content>
      <AutoRouter Loading={Loading} NoMatch={NoMatch} redirectRoutes={redirectRoutes} />
    </Content>
  </LayoutBox>
));

const Main = () => {
  const { i18n } = useTranslation();
  const currentLang = getCurrentLang(i18n);
  const local = getAntDesignLang(currentLang);

  return (
    <ConfigProvider locale={local}>
      <GlobalStyle />
      <HashRouter>
        <NormalLayoutRoutes />
      </HashRouter>
    </ConfigProvider>
  );
};

const rootElement = document.getElementById('root');
const root = createRoot(rootElement);

root.render(<Main />);
