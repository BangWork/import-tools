import { memo } from 'react';
import type { FC } from 'react';
import { map } from 'lodash-es';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';

import { downloadFile } from '@/utils/download';

export interface HelpContentPopoverProps {
  onSelected: () => void;
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
  const { onSelected } = props;
  const { t } = useTranslation();

  const HELP_LIST = [
    {
      text: 'common.help.import',
      onClick: () => {
        downloadFile(`${window.location.origin}/public/import_list.csv`, 'import_list');
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
