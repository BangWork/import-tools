const importPack = {
  initPassword: {
    title: 'Set initial password',
    tip: {
      message1:
        '1.An initial password is required for new users imported from Jira. They can also select "Forgot password?" to verify their account emails and reset passwords.',
      message2:
        "2.This password will only be displayed once and can't be changed. Please write it down.",
    },
    form: {
      tip: 'Imported data will be processed on the cloud',
      init: {
        label: 'Initial password',
        placeholder: '8–32 characters, including numbers and letters',
        error: {
          empty: 'Initial password is required',
          rule: 'Password should be 8-32 characters with numbers and letters',
        },
      },
      again: {
        label: 'Confirm password',
        placeholder: 'Please enter the password again',
        error: {
          empty: 'You need to confirm your password',
          diff: "Passwords don't match",
        },
      },
    },
    modal: {
      back: {
        title: 'Return to 「Result」',
        desc: 'After selecting [Back], you will be returned to the [Result」 page and need to provide all the information again. Proceed?',
      },
    },
    startButton: 'Start',
  },
  progressPage: {
    title: 'Jira import',
    logTitle: 'Migration log',
    timeMessage: 'Estimated time: {{totalTime}} mins. Elapsed time: {{leftTime}} mins',
    tip: {
      environment: 'Environment: {{name}}',
      startTime: 'Import time: {{time}}',
      version: 'Jira backup version: {{name}}',
      backUpTime: 'Backup time: {{time}}',
    },
    action: {
      view: 'View import scope',
      stop: 'Pause',
      continue: 'Resume',
      cancel: 'Abort',
      closeLog: 'Close',
      openLog: 'View migration log',
    },
    viewModal: {
      loading: 'Calculating',
      title: 'Import scope',
      currentInfo: {
        title: 'Jira import scope',
        version: 'Version',
        projects: 'Number of projects',
        works: 'Number of issues',
        members: 'Number of member',
        fileSize: 'File size',
        files: 'Number of File',
      },
      unit: '',
    },
    errorModal: {
      miss: {
        title: 'Jira 备份包不存在',
        text: 'Jira 备份包不存在，取消导入！',
      }
    },
    cancelModal: {
      title: 'Abort',
      desc: "Some of the data have already been imported. This action can't be undone. Proceed?",
    },
  },
  importResultPage: {
    downloadText: 'Download migration log',
    info: {
      title: 'Import complete',
      desc: 'New users imported from Jira will need to use initial passwords set by the admin to log into ONES. They can also reset their passwords by selecting "Forgot password?" in the login page.',
    },
    error: {
      title: 'Import failed.',
      desc1: 'Import failed. Jira import tool error.',
      desc2: 'Import failed. Unknown error.',
    },
  },
};

export default importPack;
