import { memo, useState } from 'react';
import { Layout, Popover } from '@ones-design/core';
import { useUpdateEffect } from 'ahooks';
import styled from 'styled-components';
import { Help, ChevronDown, ChevronUp } from '@ones-design/icons';
import { useTranslation } from 'react-i18next';
import { LANGS, getCurrentLang } from '@/i18n';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';
import { map } from 'lodash-es';

import HelpContentPopover from './help_content_popover';

const HeaderBox = styled(Layout.Header)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff !important;
  height: 48px !important;
  padding: 0 0 0 20px !important;
  margin-bottom: 10px;
  box-shadow: 0px 0px 2px rgba(48, 48, 48, 0.05), 0px 1px 2px rgba(48, 48, 48, 0.2);
`;
const Title = styled.div`
  font-weight: 500;
  font-size: 18px;
  line-height: 26px;
  color: #606060;
  cursor: pointer;
`;
const TopRightBanner = styled.div`
  display: flex;
  align-items: center;
  color: #606060;
  cursor: pointer;
`;
const HelpStyled = styled(Help)`
  font-size: 21px;
  margin: 0 20px;
`;

const GlobalOutlinedStyled = styled(Help)`
  font-size: 14px;
  margin-right: 10px;
`;

const LangBox = styled.div`
  padding-right: 20px;
  border-right: 1px solid #eaeaea;
`;

const LangContentPopoverItem = styled.div`
  cursor: pointer;
  padding: 5px 10px;
  &:hover {
    background: #f5f5f5;
  }
`;

/**
 * top Head
 */
const Header = memo(() => {
  const { t, i18n } = useTranslation();
  const currentLang = getCurrentLang();
  const [selected, setSelected] = useState(currentLang);
  const [open, setOpen] = useState(false);
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const navigate = useNavigate();

  const options = LANGS.map((key) => ({
    value: key,
    label: t(key),
    onClick: () => handleLangSelected(key),
  }));

  useUpdateEffect(() => {
    i18n.changeLanguage(selected);
    dayjs.locale(selected);
    window.location.reload();
  }, [selected]);

  const handleLangSelected = (key) => {
    setSelected(key);
    setOpen(false);
  };

  // click popover content need close popover
  const handleHidePopover = () => {
    setOpen(false);
  };

  const handleOpenChange = (newOpen: boolean) => {
    setOpen(newOpen);
  };

  const handleOpenPopoverChange = (newOpen: boolean) => {
    setIsPopoverOpen(newOpen);
  };

  const handleToHome = () => {
    navigate('/page/home');
  };

  return (
    <HeaderBox style={{}}>
      <Title onClick={handleToHome}>{t('common.jira.title')}</Title>
      <TopRightBanner>
        <LangBox>
          <GlobalOutlinedStyled />
          <Popover
            className={'oac-flex-1'}
            visible={isPopoverOpen}
            placement="bottom"
            trigger="click"
            content={
              <div>
                {map(options, (item) => (
                  <LangContentPopoverItem key={item.label} onClick={item.onClick}>
                    {t(item.value)}
                  </LangContentPopoverItem>
                ))}
              </div>
            }
            onVisibleChange={handleOpenPopoverChange}
          >
            <span style={{ fontWeight: '500' }}>{t(selected)}</span>
            {!isPopoverOpen ? (
              <ChevronDown className={'oac-ml-2'} />
            ) : (
              <ChevronUp className={'oac-ml-2'} />
            )}
          </Popover>
        </LangBox>
        <Popover
          visible={open}
          placement="bottom"
          content={<HelpContentPopover onSelected={handleHidePopover} />}
          trigger="click"
          onVisibleChange={handleOpenChange}
        >
          <HelpStyled />
        </Popover>
      </TopRightBanner>
    </HeaderBox>
  );
});

export default Header;
