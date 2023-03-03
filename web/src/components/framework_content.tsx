import { memo } from 'react';
import type { FC, ReactNode } from 'react';
import styled from 'styled-components';

export interface FrameworkContentProps {
  title: string | ReactNode;
  footer?: string | ReactNode;
  children: string | ReactNode;
  width?: string;
  className?: string;
}

const Box = styled.div<Pick<FrameworkContentProps, 'width'>>`
  border-radius: 2px;
  padding: 20px 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  box-shadow: 0px 0px 2px rgba(48, 48, 48, 0.05), 0px 1px 2px rgba(48, 48, 48, 0.2);
  border-radius: 3px;
  background: #ffffff;
`;

const Head = styled.div`
  padding: 0 0 20px 20px;
  width: 100%;
  font-weight: 500;
  font-size: 18px;
  line-height: 26px;
`;

const ChildrenBox = styled.div`
  padding-left: 20px;
  flex: 1;
`;
const Footer = styled.div`
  border-top: 1px solid #eaeaea;
  padding-top: 20px;
  width: 100%;
`;

/**
 * Provides a unified layout container
 */
const FrameworkContent: FC<FrameworkContentProps> = memo((props) => {
  const { title, children, footer, width, className } = props;

  return (
    <Box className={className}>
      <Head>{title}</Head>
      <ChildrenBox>{children}</ChildrenBox>
      {footer ? <Footer>{footer}</Footer> : null}
    </Box>
  );
});

export default FrameworkContent;
