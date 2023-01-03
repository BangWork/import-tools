import { Spin } from 'antd';
import type { FC } from 'react';
import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import type { SpinProps } from 'antd/es/spin';


/**
 * Spin need more tip to user
 */
const Loading: FC<SpinProps> = memo((props) => {
  const { t } = useTranslation();

  return <Spin tip={t('common.loading')} {...props} />;
});

export default Loading;
