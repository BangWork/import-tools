
import styled from "styled-components"


const ContentBox = styled.div`
  width:473px;
  height:261px;
  border:1px solid #E8E8E8;
  padding:20px;
  display:flex;
  border-radius:3px;
`
const LeftIcon = styled.div`
  width:48px;
  height:48px;
  background: #EAF3FC;
  display:flex;
  justify-content:center;
  align-items:center;
  border-radius:2.88px

`
const RightBox = styled.div`
  flex:1;
  padding-left:10px;
`

const RightText = styled.div`
  height:161px;
  margin-bottom:15px;
`
const TitleBox = styled.div`
  weight:500;
  font-size:16px;
  line-height:24px;
`


const TextBox = (props) => {
  const {title,descriptionText,icon,children,className}  = props
  return <ContentBox className={className}>
    <LeftIcon>{ icon}</LeftIcon>
    <RightBox>
      <RightText>
        <TitleBox>{title }</TitleBox>
        <div>{ descriptionText}</div>
      </RightText>
      <>{ children}</>
    </RightBox>

  </ContentBox>
}

export default TextBox
