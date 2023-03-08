import { Alert, Button } from '@ones-design/core';
import { Checkbox } from 'antd';
import { Modal } from '@ones-design/core';
import { t } from 'i18next';
import { memo, useState } from 'react';
import { useNavigate } from 'react-router';
import styled from 'styled-components';
import TextBox from './text_content';
import { Page } from '@ones-design/icons';
import { useWhyDidYouUpdate } from 'ahooks';

const BorderBox = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
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
    // judge is login
    setIsModalOpen(true);
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
    setIsModalOpen(false);
  };

  const handleCheckBox = () => {
    !isConfirm && setIsShowTips(false);
    setIsConfirm(!isConfirm);
  };

  useWhyDidYouUpdate('modal', {
    isConfirm,
    isModalOpen,
  });

  return (
    <BorderBox>
      <div className="oac-p-4">
        <Alert type="info">{t('environment.tip.message1')}</Alert>
        <TitleBox>
          <TitleNumber>1</TitleNumber>
          {t('common.back')}
        </TitleBox>
        <DescriptionBox>
          <TextBox
            descriptionText={'描述信息哈哈哈'}
            title={'标题'}
            icon={<PageStyled></PageStyled>}
            className={'oac-mr-4'}
          >
            <Button>{'查看'}</Button>
          </TextBox>
          <TextBox
            descriptionText={'描述信息哈哈哈'}
            title={'标题'}
            icon={<PageStyled></PageStyled>}
          >
            <Button>{t('common.back')}</Button>
          </TextBox>
        </DescriptionBox>
        <TitleBox>
          <TitleNumber>2</TitleNumber>
          {'迁移评估'}
        </TitleBox>
        <DescriptionBox>
          <TextBox
            descriptionText={'描述信息哈哈哈'}
            title={'标题'}
            icon={<PageStyled></PageStyled>}
            className={'oac-mr-4'}
          >
            <Button onClick={handleToConfirm}>{'开始评估'}</Button>
            <Modal
              title="Basic Modal"
              visible={isModalOpen}
              onOk={handleOk}
              onCancel={handleCancel}
            >
              <div className={'oac-p-2'}>Some contents...</div>
              <Checkbox onChange={handleCheckBox} checked={isConfirm}>
                {t('environment.tip.message1')}
              </Checkbox>
              {isShowTips ? <div style={{ color: 'red' }}>{t('common.back')}</div> : <br></br>}
            </Modal>
          </TextBox>
        </DescriptionBox>
      </div>
    </BorderBox>
  );
});

export default Home;
