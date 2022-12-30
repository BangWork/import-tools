import { memo } from 'react';
import type { FC, ReactNode } from 'react';
import styled from 'styled-components';

export interface ModalContentProps {
  title: string | ReactNode;
  footer?: string | ReactNode;
  children: string | ReactNode;
  width?: string;
  className?: string;
}

const Box = styled.div<Pick<ModalContentProps, 'width'>>`
  box-shadow: 0px 3px 6px -4px rgba(0, 0, 0, 0.12), 0px 6px 16px rgba(0, 0, 0, 0.08),
    0px 9px 28px 8px rgba(0, 0, 0, 0.05);
  border-radius: 2px;
  width: ${(props) => props?.width || 'auto'};
`;

const Head = styled.div<Pick<ModalContentProps, 'width'>>`
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 16px 24px;
  box-shadow: inset 0px -1px 0px #f0f0f0;
  width: ${(props) => props?.width || 'auto'};
  min-width: ${(props) => props?.width || '600px'};
  font-weight: 600;
  font-size: 16px;
  line-height: 24px;
`;

const Footer = styled.div`
  padding: 10px 16px;

  width: 100%;
  height: 52px;
  box-shadow: inset 0px 1px 0px #f0f0f0;
`;

/**
 * Provides a unified layout container
 */
const ModalContent: FC<ModalContentProps> = memo((props) => {
  const { title, children, footer, width, className } = props;

  return (
    <Box className={className} width={width}>
      <Head width={width}>{title}</Head>
      <div>{children}</div>
      {footer ? <Footer>{footer}</Footer> : null}
    </Box>
  );
});

export default ModalContent;
