import styled from 'styled-components';

const ContentBox = styled.div`
  width: 473px;
  height: 195px;
  border: 1px solid #e8e8e8;
  padding: 20px;
  display: flex;
  border-radius: 3px;
`;
const LeftIcon = styled.div`
  width: 48px;
  height: 48px;
  background: #eaf3fc;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 2.88px;
`;
const RightBox = styled.div`
  flex: 1;
  padding-left: 10px;
`;

const RightText = styled.div`
  height: 95px;
  margin-bottom: 15px;
  color: #606060;
`;
const TitleBox = styled.div`
  weight: 500;
  font-size: 16px;
  line-height: 24px;
  margin-bottom: 5px;
  color: #303030;
`;

const TextBox = (props) => {
  const { title, descriptionText1, descriptionText2, icon, children, className } = props;
  return (
    <ContentBox className={className}>
      <LeftIcon>{icon}</LeftIcon>
      <RightBox>
        <RightText>
          <TitleBox>{title}</TitleBox>
          <div>{descriptionText1}</div>
          {descriptionText2 ? <div>{descriptionText2}</div> : null}
        </RightText>
        <>{children}</>
      </RightBox>
    </ContentBox>
  );
};

export default TextBox;
