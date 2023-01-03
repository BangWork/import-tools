import { memo, ReactNode } from 'react';
import type { FC } from 'react';
import { Progress, Typography } from 'antd';
import { isUndefined } from 'lodash-es';
import { FileOutlined } from '@ant-design/icons';
import type { ProgressProps } from 'antd/es/progress';
import styled from 'styled-components';

export interface BusinessProgressProps extends ProgressProps {
  /** top message */
  message?: string | ReactNode;
  bottomMessage: string | ReactNode;
  /** left time of bottom */
  leftTime?: number;
  /** total time of bottom */
  totalTime?: number;
  percent?: number;
}

const FileOutlinedStyled = styled(FileOutlined)`
  font-size: 26px;
  color: #1890ff;
`;

const ProgressStyled = styled(Progress)`
  width: 300px;
  height: 8px;
  line-height: 8px;
  margin: 0;
`;
const Message = styled.div`
  width: 300px;
  font-size: 14px;
  line-height: 20px;
  color: #303030;
`;

const BusinessProgress: FC<BusinessProgressProps> = memo((props) => {
  const { message = '', percent = 0, bottomMessage } = props;

  return (
    <div className="flex items-center justify-center">
      <FileOutlinedStyled className="mr-4" />
      <div className="flex flex-col justify-center">
        <Message>
          <Typography.Text ellipsis={{ tooltip: message }}>{message}</Typography.Text>
        </Message>
        <ProgressStyled percent={percent} showInfo={false} />
        <Message>
          <Typography.Text ellipsis={{ tooltip: bottomMessage }}>{bottomMessage}</Typography.Text>
        </Message>
      </div>
    </div>
  );
});

export default BusinessProgress;
