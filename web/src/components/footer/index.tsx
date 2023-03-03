import { Button } from 'antd';
import type { ButtonProps } from 'antd';
import { memo, ReactNode, FC } from 'react';
import { t } from 'i18next';
import styled from 'styled-components';

export interface FooterType {
  fun?: (...args: any[]) => void;
  text?: string;
  isLoading?: boolean;
  isDisabled?: boolean;
  type?: ButtonProps['type'];
  htmlType?: ButtonProps['htmlType'];
}
export interface FooterProps {
  handleCancleMigrate?: FooterType;
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
  const { handleBack, handleNext, handleCancleMigrate, children, width, className } = props;
  return (
    <Box className={className} width={width}>
      {handleCancleMigrate ? (
        <Button
          type={handleCancleMigrate.type}
          className="oac-mr-2"
          onClick={handleCancleMigrate.fun}
          loading={handleCancleMigrate.isLoading}
          disabled={handleCancleMigrate.isDisabled}
        >
          {t(handleCancleMigrate.text || 'common.canclemigrate')}
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
