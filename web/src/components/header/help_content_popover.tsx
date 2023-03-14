import { memo, useMemo } from 'react';
import type { FC } from 'react';
import { map } from 'lodash-es';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
import { useLocation } from 'react-router';

import { downloadFile } from '@/utils/download';

export interface HelpContentPopoverProps {
  onSelected: () => void;
  label: string;
}

const HelpContentPopoverItem = styled.div`
  cursor: pointer;
  padding: 5px 10px;

  &:hover {
    background: #f5f5f5;
  }
`;

/**
 * content popover of help
 */
const HelpContentPopover: FC<HelpContentPopoverProps> = memo((props) => {
  const { onSelected, label } = props;
  const { t } = useTranslation();
  const en = 'en';
  let controlFlag = false;
  const location = useLocation();
  const HELP_LIST = useMemo(() => {
    if (
      location.pathname === '/page/analyze/environment' ||
      location.pathname === '/page/analyze/pack'
    ) {
      controlFlag = false;
    } else {
      controlFlag = true;
    }
    if (controlFlag) {
      return [
        {
          text: 'common.help.import',
          onClick: () => {
            if (label === en) {
              downloadFile(
                `${window.location.origin}/public/mappingEn.xlsx`,
                'Jira import mapping form'
              );
            } else {
              downloadFile(`${window.location.origin}/public/mappingZh.xlsx`, 'Jira数据导入清单');
            }
          },
        },
        {
          text: 'common.help.document',
          onClick: () => {
            downloadFile(`${window.location.origin}/public/help.pdf`, 'help');
          },
        },
        {
          text: 'common.help.download',
          onClick: () => {
            downloadFile(
              `${
                import.meta.env.VITE_PROXY_DOMAIN_REAL || window.location.origin
              }/import/log/download/all`
            );
          },
        },
      ];
    } else {
      return [
        {
          text: 'common.help.import',
          onClick: () => {
            if (label === en) {
              downloadFile(
                `${window.location.origin}/public/mappingEn.xlsx`,
                'Jira import mapping form'
              );
            } else {
              downloadFile(`${window.location.origin}/public/mappingZh.xlsx`, 'Jira数据导入清单');
            }
          },
        },
        {
          text: 'common.help.document',
          onClick: () => {
            downloadFile(`${window.location.origin}/public/help.pdf`, 'help');
          },
        },
      ];
    }
  }, [location.pathname]);

  return (
    <div onClick={onSelected}>
      {map(HELP_LIST, (item) => (
        <HelpContentPopoverItem key={item.text} onClick={item.onClick}>
          {t(item.text)}
        </HelpContentPopoverItem>
      ))}
    </div>
  );
});

export default HelpContentPopover;
