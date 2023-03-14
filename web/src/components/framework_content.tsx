import { memo } from 'react';
import type { FC, ReactNode } from 'react';
import styled from 'styled-components';
import { Input, Button } from '@ones-design/core';
import { Download, WorkItem } from '@ones-design/icons';
import { t } from 'i18next';
import { debounce } from 'lodash-es';
export interface FrameworkContentProps {
  title: string | ReactNode;
  footer?: string | ReactNode;
  children: string | ReactNode;
  width?: string;
  className?: string;
  search?: { text?: string; fun: (searchValue: any) => void };
  download?: { text?: string; fun: () => void };
  config?: { text?: string; fun: () => void };
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
  padding-left: 20px;
  font-weight: 500;
  font-size: 18px;
  line-height: 26px;
`;
const RightHeadBox = styled.div`
  display: flex;
  align-items: center;
}`;

const ChildrenBox = styled.div`
  padding-left: 20px;
  padding-right: 20px;
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
  const { title, children, footer, width, className, search, download, config } = props;
  const handleSearch = debounce((value) => {
    search && search.fun(value);
  }, 400);
  return (
    <Box className={className}>
      <div className="oac-flex oac-items-center oac-justify-between oac-pb-4">
        <Head>{title}</Head>
        <RightHeadBox>
          {download ? (
            <Button className="oac-mr-1" onClick={download.fun} type="text">
              <Download />
              {t(download.text || 'common.download')}
            </Button>
          ) : null}
          {config ? (
            <Button className="oac-mr-1" onClick={config.fun} type="text">
              <WorkItem />
              {t(config.text)}
            </Button>
          ) : null}
          {search ? (
            <Input.Search
              style={{ width: '280px', marginRight: '20px', height: '32px' }}
              onChange={handleSearch}
              placeholder={t(search.text)}
            />
          ) : null}
        </RightHeadBox>
      </div>
      <ChildrenBox>{children}</ChildrenBox>
      {footer ? <Footer>{footer}</Footer> : null}
    </Box>
  );
});

export default FrameworkContent;
