import { Button } from 'antd';
import type { ButtonProps } from 'antd';
import { memo, ReactNode, FC } from 'react';
import { t } from 'i18next';
import styled from 'styled-components';
import { useNavigate } from 'react-router';

import { Modal } from '@ones-design/core';
import { loginApi } from '@/api/';
export interface FooterType {
  fun?: (...args: any[]) => void;
  text?: string;
  isLoading?: boolean;
  isDisabled?: boolean;
  type?: ButtonProps['type'];
  htmlType?: ButtonProps['htmlType'];
}
export interface FooterProps {
  handleCancelMigrate?: FooterType;
  handleBack?: FooterType;
  handleNext?: FooterType;
  children?: string | ReactNode;
  width?: string;
  className?: string;
}

const Box = styled.div<Pick<FooterProps, 'width'>>`
  display: flex;
  justify-content: flex-end;
  width: ${(props) => props?.width || 'full'};
`;

const Footer: FC<FooterProps> = memo((props) => {
  const navigate = useNavigate();
  const handleToConfirm = () => {
    const url = window.localStorage.getItem('url') || '';
    const email = window.localStorage.getItem('email') || '';
    const password = window.localStorage.getItem('password') || '';
    loginApi({
      url,
      email,
      password,
    })
      .then((res) => {
        Modal.error({
          title: t('common.cancelMigrate'),
          content: t('common.cancelMigrateTip'),
          okText: t('common.cancelMigrate'),
          onOk: handleHome,
          cancelText: t('common.continueMigrate'),
          onCancel: () => {
            return null;
          },
        });
      })
      .catch((e) => {
        const { err_code } = e;

        if (err_code === 'NotSuperAdministratorError') {
          Modal.warning({
            title: t('common.nonONESTeamAdministratorAccount'),
            content: t('common.nonTeamTips'),
            okText: t('common.ok'),
            onOk: handleHome,
          });
        }
        if (err_code === 'NotOrganizationAdministratorError') {
          Modal.warning({
            title: t('nonONESOrganizationAdministratorAccount'),
            content: t('nonOrganizationTips'),
            okText: t('common.ok'),
            onOk: handleHome,
          });
        }
      });
  };

  const handleHome = () => {
    navigate('/page/home');
  };

  const { handleBack, handleNext, handleCancelMigrate, children, width, className } = props;
  return (
    <Box className={className} width={width}>
      {handleCancelMigrate ? (
        <Button
          type={handleCancelMigrate.type}
          className="oac-mr-2"
          onClick={handleCancelMigrate.fun || handleToConfirm}
          loading={handleCancelMigrate.isLoading}
          disabled={handleCancelMigrate.isDisabled}
        >
          {t(handleCancelMigrate.text || 'common.cancelMigrate')}
        </Button>
      ) : null}
      {handleBack ? (
        <Button
          type={handleBack.type}
          className="oac-mr-2"
          onClick={handleBack.fun}
          loading={handleBack.isLoading}
          disabled={handleBack.isDisabled}
        >
          {t(handleBack.text || 'common.back')}
        </Button>
      ) : null}
      {handleNext ? (
        <Button
          htmlType={handleNext.htmlType}
          type={handleNext.type || 'primary'}
          className="oac-mr-4"
          onClick={handleNext.fun}
          loading={handleNext.isLoading}
          disabled={handleNext.isDisabled}
        >
          {t(handleNext.text || 'common.nextStep')}
        </Button>
      ) : null}
      <div>{children}</div>
    </Box>
  );
});

export default Footer;
