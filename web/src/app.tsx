import { createRoot } from 'react-dom/client';
import { HashRouter } from 'react-router-dom';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
import { Layout, ConfigProvider } from '@ones-design/core';

import AutoRouter from '@/router';
import { redirectRoutes } from '@/router/routes';

import Header from '@/components/header';
import LeftSteps from '@/components/steps';
import Loading from '@/components/loading';
import NoMatch from '@/components/no_match';
import GlobalStyle from './global';

import './index.css';
import { getCurrentLang, getAntDesignLang } from '@/i18n';

const FollowBox = styled.div`
  height: calc(100% - 58px);
  display: flex;
`;
const Content = styled.div`
  width: 100%;
  overflow-y: auto;
  padding: 0 10px 10px 0;
`;

/** default layout routes */
const NormalLayoutRoutes = () => (
  <Layout className="oac-h-full oac-bg-white" style={{ background: '#EAEAEA' }}>
    <Header />
    <FollowBox>
      <LeftSteps />
      <Content>
        <AutoRouter
          Loading={Loading}
          NoMatch={NoMatch}
          redirectRoutes={redirectRoutes}
        ></AutoRouter>
      </Content>
    </FollowBox>
  </Layout>
);

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
