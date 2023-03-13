import { Alert, Button } from '@ones-design/core';
import { Checkbox } from 'antd';
import { Modal } from '@ones-design/core';
import { Edit, Launch } from '@ones-design/icons';
import { t } from 'i18next';
import { memo, useState } from 'react';
import { useNavigate } from 'react-router';
import styled from 'styled-components';
import TextBox from './text_content';
import { Page } from '@ones-design/icons';
import { useWhyDidYouUpdate } from 'ahooks';
import { isHasCookie } from '@/utils/getCookie';

const BorderBox = styled.div`
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: start;
  box-shadow: 0px 0px 2px rgba(48, 48, 48, 0.05), 0px 1px 2px rgba(48, 48, 48, 0.2);
  border-radius: 3px;
  background: #ffffff;
  margin: 0 0 10px 10px;
`;
const PageStyled = styled(Page)`
  fontsize: 38;
  color: #338fe5;
`;
const TitleBox = styled.div`
  font-weight: 500;
  line-height: 24px;
  font-size: 16px;
  display: flex;
  padding-top: 20px;
  padding-bottom: 10px;
`;
const TitleNumber = styled.div`
  width: 24px;
  height: 24px;
  border-radius: 12px;
  background: #e8e8e8;
  font-size: 14px;
  line-height: 22px;
  text-align: center;
  margin-right: 10px;
`;
const DescriptionBox = styled.div`
  padding-left: 34px;
  display: flex;
`;

const Home = memo(() => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isConfirm, setIsConfirm] = useState(false);
  const [isShowTips, setIsShowTips] = useState(false);
  const navigate = useNavigate();

  const handleToConfirm = () => {
    if (isHasCookie()) {
      navigate('/page/analyze/environment');
    } else {
      setIsModalOpen(true);
    }
  };

  const handleOk = () => {
    if (!isConfirm) {
      setIsShowTips(true);
    } else {
      setIsModalOpen(false);
      navigate('/page/analyze/environment');
    }
  };

  const handleCancel = () => {
    setIsShowTips(false);
    setIsModalOpen(false);
  };

  const handleCheckBox = () => {
    setIsShowTips(false);
    setIsConfirm(!isConfirm);
  };

  useWhyDidYouUpdate('modal', {
    isConfirm,
    isModalOpen,
  });

  return (
    <BorderBox>
      <div className="oac-p-4">
        <Alert type="info">{t('home.alert')}</Alert>
        <TitleBox>
          <TitleNumber>1</TitleNumber>
          {t('home.preparation')}
        </TitleBox>
        <DescriptionBox>
          <TextBox
            descriptionText1={t('home.guideText1')}
            descriptionText2={t('home.guideText2')}
            title={t('home.toolUserGuide')}
            icon={<PageStyled></PageStyled>}
            className={'oac-mr-4'}
          >
            <Button>{t('home.viewUserGuide')}</Button>
          </TextBox>
          <TextBox
            descriptionText1={t('home.effectText1')}
            title={t('home.evaluateMigrationEffectiveness')}
            icon={<PageStyled></PageStyled>}
          >
            <Button type="primary">
              <Edit />
              {t('home.startEvaluation')}
            </Button>
          </TextBox>
        </DescriptionBox>
        <TitleBox>
          <TitleNumber>2</TitleNumber>
          {t('home.migrateData')}
        </TitleBox>
        <DescriptionBox>
          <TextBox
            descriptionText1={t('home.dataText1')}
            descriptionText2={t('home.dataText2')}
            title={t('home.migrateJiraData')}
            icon={<PageStyled></PageStyled>}
            className={'oac-mr-4'}
          >
            <Button onClick={handleToConfirm} type="primary">
              {t('home.startMigration')}
            </Button>
            <Modal
              style={{ position: 'relative' }}
              title={t('home.modal.title')}
              visible={isModalOpen}
              onOk={handleOk}
              onCancel={handleCancel}
            >
              <div className={'oac-p-2'}>{t('home.modal.text')}</div>
              <Checkbox onChange={handleCheckBox} checked={isConfirm}>
                {t('home.modal.agreeText1')}
                <a target="_blank" rel="noopener noreferrer" href="www.baidu.com">
                  {t('home.modal.term')}
                  <Launch style={{ margin: '0 5px' }} />
                </a>
                {t('home.modal.agreeText2')}
              </Checkbox>
              {isShowTips ? (
                <div style={{ color: 'red', position: 'absolute', bottom: '50px' }}>
                  {t('home.modal.warning')}
                </div>
              ) : null}
            </Modal>
          </TextBox>
        </DescriptionBox>
      </div>
    </BorderBox>
  );
});

export default Home;
