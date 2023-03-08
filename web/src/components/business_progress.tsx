import { memo, ReactNode } from 'react';
import type { FC } from 'react';
import { Progress } from '@ones-design/core';
import styled from 'styled-components';

export interface BusinessProgressProps {
  /** top message */
  title?: string | ReactNode;
  statusText?: string | ReactNode;
  status?: 'success' | 'active' | 'fail';
  percentDescription?: string | ReactNode;
  percentTimeText?: string | ReactNode;
  bottomMessage: string | ReactNode;
  /** left time of bottom */
  leftTime?: number;
  /** total time of bottom */
  totalTime?: number;
  percent?: number;
}

const ProgressBoxStyled = styled.div`
  padding: 10px 20px;
  height: 108px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border: 1px solid #dedede;
  border-radius: 3px;
`;

const TitleStyled = styled.div`
  display: flex;
  align-items: center;
  font-size: 18px;
  line-height: 26px;
  font-weight: 500;
`;

const StatusStyled = styled.div`
  padding: 0 5px;
  margin-left: 10px;
  font-weight: 400;
  font-size: 12px;
  height: 20px;
  line-height: 18px;
  border-radius: 10px;
  border: 1px solid #dedede;
`;

const ProgressStyled = styled(Progress)`
  width: 150px;
  height: 8px;
  line-height: 8px;
  margin: 0 20px;
`;

const Message = styled.div`
  font-size: 12px;
  line-height: 20px;
`;
const BusinessProgress: FC<BusinessProgressProps> = memo((props) => {
  const {
    title = '',
    status,
    statusText,
    percent = 0,
    percentDescription,
    percentTimeText,
    bottomMessage,
  } = props;
  const statusMap = {
    success: '#24B47E',
    active: '#F0A100',
    fail: '#E52727',
  };
  return (
    <ProgressBoxStyled>
      <TitleStyled>
        {title}
        <StatusStyled style={{ borderColor: statusMap[status], color: statusMap[status] }}>
          {statusText}
        </StatusStyled>
      </TitleStyled>
      <div className="oac-flex oac-items-center">
        {percentDescription}
        {status !== 'fail' ? (
          <ProgressStyled percent={percent} showInfo={false} />
        ) : (
          <ProgressStyled percent={percent} showInfo={false} status="exception" />
        )}
        {percentTimeText}
      </div>
      <Message>{bottomMessage}</Message>
    </ProgressBoxStyled>
  );
});

export default BusinessProgress;
