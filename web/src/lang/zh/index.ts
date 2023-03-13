import common from './common';
import analyze from './analyze';
import language from './language';
import importPack from './import_pack';
import requestError from './request_error';
import home from './home'

export default {
  common,
  requestError,
  home,
  ...analyze,
  ...language,
  ...importPack,
};
