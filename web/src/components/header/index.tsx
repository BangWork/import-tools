import { memo, useEffect, useState } from 'react';
import { Layout, Select, Popover } from 'antd';
import { useUpdateEffect } from 'ahooks';
import styled from 'styled-components';
import { QuestionCircleOutlined, GlobalOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { LANGS, getCurrentLang } from '@/i18n';
import dayjs from 'dayjs';

import HelpContentPopover from './help_content_popover';

const HeaderBox = styled(Layout.Header)`
  margin: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff !important;
  box-shadow: 0px 0px 2px rgba(48, 48, 48, 0.05), 0px 1px 2px rgba(48, 48, 48, 0.2);
`;
const Title = styled.div`
  font-weight: 500;
  font-size: 18px;
  line-height: 26px;
  color: #606060;
  padding-right: 20px;
  border-right: 1px solid #c7c7c7;
`;
const TopRightBanner = styled.div`
  display: flex;
  align-items: center;
  color: #606060;
`;
const QuestionCircleOutlinedStyled = styled(QuestionCircleOutlined)`
  font-size: 21px;
  margin: 0 20px;
`;

const GlobalOutlinedStyled = styled(GlobalOutlined)`
  font-size: 21px;
`;

const LangBox = styled.div`
  display: flex;
  align-items: center;
`;
const LangSelect = styled(Select)`
  width: 100px;
`;

/**
 * top Head
 */
const Header = memo((props) => {
  const { t, i18n } = useTranslation();
  const currentLang = getCurrentLang();
  const [selected, setSelected] = useState(currentLang);
  const [open, setOpen] = useState(false);

  const options = LANGS.map((key) => ({
    value: key,
    label: t(key),
  }));

  useUpdateEffect(() => {
    i18n.changeLanguage(selected);
    dayjs.locale(selected);
    window.location.reload();
  }, [selected]);

  const handleLangSelected = (key) => {
    setSelected(key);
  };

  // click popover content need close popover
  const handleHidePopover = () => {
    setOpen(false);
  };

  const handleOpenChange = (newOpen: boolean) => {
    setOpen(newOpen);
  };

  return (
    <HeaderBox style={{}}>
      <Title>{t('common.jira.title')}</Title>
      <TopRightBanner>
        <LangBox>
          <GlobalOutlinedStyled />
          <LangSelect
            value={selected}
            bordered={false}
            onSelect={handleLangSelected}
            options={options}
          />
        </LangBox>
        <Popover
          open={open}
          placement='bottom'
          content={<HelpContentPopover onSelected={handleHidePopover} />}
          trigger='click'
          onOpenChange={handleOpenChange}
        >
          <QuestionCircleOutlinedStyled />
        </Popover>
      </TopRightBanner>
    </HeaderBox>
  );
});

export default Header;
