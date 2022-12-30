import { createGlobalStyle } from 'styled-components';

/** global style */
export default createGlobalStyle`
    *{
        margin: 0;
        padding: 0;
        outline:0;
        box-sizing:border-box;
        font-style: normal;
        &::before, &::after {
          box-sizing: border-box;
        }
    }

    #root{
        height: 100%;
    }

    html,body {
      font-size: 14px;
      font-weight: 400;
      -webkit-font-smoothing: antialiased;
      -moz-osx-font-smoothing: grayscale;
      -webkit-overflow-scrolling: touch;
      line-height: 1.33;
      height: 100%;
      min-width: 1028px;
      width: 100%;
    }

    input,button,select,textarea {
      font-family: inherit;
      font-size: inherit;
      line-height: inherit;
    }

    a {
      text-decoration: none;
    }

    figure {
      margin: 0;
    }

    img {
      vertical-align: middle;
    }

    label {
      font-weight: normal;
    }
  `;
