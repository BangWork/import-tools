const analyze = {
  backupPage: {
    guide: {
      alert: {
        desc: 'Please get the Location of Jira Local Home in the Jira system and enter the following information.',
        link: 'Find the location of the Jira home directory',
      },
      step1: {
        title: 'Find Jira backup path',
        desc1: '1.1 Log into Jira as Jira admin；',
        desc2: '1.2 In the top navigation bar of Jira, click <1/> Jira administration > System；',
        desc3: '1.3 In the left navigation, click SYSTEM SUPPORT > System Info。',
      },
      step2: {
        title: 'Provide Jira backup info',
        desc: 'Scroll down to the section File Paths, copy and paste the File Paths>Location of JIRA Local Home path into the following input box.',
        link: 'View the default Jira local file path',
      },
    },
    form: {
      desc: 'Input the path you copied from Location of JIRA Local Home, then click "Get Jira backups" and select the backups you want to import.',
      localHome: {
        label: 'Location of Jira Local Home',
        get: 'Get Jira backup',
        emptyError: 'Please input Location of Jira Local Home.',
        serverError: 'Please input correct path.',
      },
      backup: {
        label: 'Jira backup name',
        placeholder: 'Please select which backup to import',
        emptyError: 'Please select Jira backup',
      },
      tip: 'Select which backup to import',
    },
  },
  environment: {
    title: 'Provide target ONES environment info',
    tip: {
      message1: '1. Get ONES domain / IP from operation personnel',
      message2: '2. Only admin can perform this action',
    },
    url: {
      label: 'ONES domain / IP',
      emptyError: 'Please enter domain/IP',
      serverError: 'Please input correct domain / IP.',
      placeholder: 'eg：http://ones.com OR https://ones.com',
    },
    email: {
      label: "Importer's ONES account email",
      emptyError: 'Please enter email',
    },
    password: {
      label: 'Password',
      emptyError: "Password can't be empty",
    },
    serverError: {
      count: {
        title: 'Please enter correct account and password',
        desc: 'Please log into ONES to verify the correct account and password, then log into this tool 10 minutes later.',
      },
      account: 'Please enter correct account and password',
      team: 'This user is not a team admin. Please try again',
      organize: 'This user is not an organization admin. Please try again',
    },
    startButton: 'Analyze',
  },
  analyzeProgress: {
    title: 'Analyze Jira backups',
    timeMessage: 'Estimated time: {{totalTime}} mins. Elapsed time: {{leftTime}} mins',
    tip: {
      environment: '1、Environment: {{name}}',
      time: '2、Analysis started at: {{time}}',
    },
    cancel: {
      text: 'Abort',
      success: 'Analysis aborted',
      fail: 'Failed to abort analysis',
      loading: 'Aborting analysis...',
      desc: 'Aborting the process will return you to the [Provide Jira backup info] page. Proceed?',
    },
    status: {
      doing: 'Analyzing',
    },
    fail: {
      title: 'Analysis failed',
      normalDesc:
        'Failed to analysis [{{name}}]. Please make sure you provided the correct Jira backup info.',
      onExistDesc: 'Analysis failed. [{{name}}] does not exist',
    },
  },
  analyzeResult: {
    title: 'Results',
    current: {
      title: 'Jira backup info',
      version: 'Version',
      projects: 'Projects',
      works: 'Issues',
      members: 'Members',
      fileSize: 'File size',
      files: 'Number of files',
      unit: '',
    },
    environment: {
      title: 'Target ONES environment info',
      history: 'Imported data from Jira {{version}} on {{time}}',
      empty: 'No data imported from Jira',
    },
    modal: {
      back: {
        title: 'Return to [Provide Jira backup info]',
        desc: 'After selecting [Back], you will be returned to the [Provide Jira backup info] page and need to provide all the information again. Proceed?',
      },
    },
  },
  teamPage: {
    title: 'Select a targeted ONES team',
    form: {
      label: 'Import data into a ONES team',
      placeholder: 'Select a targeted ONES team',
    },
    error: {
      packDiff: {
        title: "Domains don't match",
        desc: "Jira {{packVersion}} can't be imported, because this team has already selected and imported data from Jira {{version}} which belongs to a different domain on {{time}}. Please import from another Jira backup.",
      },
      importDiff: {
        title: 'Jira backup data imported into ONES team',
        desc: 'Confirm data import from Jira {{packVersion}}? This team has already selected and imported data from Jira {{version}} on {{time}}.',
      },
      importVersionDiff: {
        title: 'Jira backup data imported into ONES team',
        desc: 'Jira {{packVersion}} will be imported into the selected ONES team. Some data may not be properly imported. Proceed? This team has already selected and imported data from Jira {{version}} on {{time}}.',
      },
      versionDiff: {
        title: 'Incompatible Jira version',
        desc: 'Jira {{packVersion}} will be imported into the selected ONES team. Some data may not be properly imported. Proceed?',
      },
      warning: {
        cancel: 'Resume',
        ok1: 'Select a ONES team again',
        ok2: 'Select Jira backup again',
      },
    },
    buttonTip: 'Please select a ONES team',
  },
  importProject: {
    title: 'Select Jira projects',
    sourceTitle: 'Jira projects',
    targetTitle: 'Projects imported to ONES',
    local: {
      searchPlaceholder: 'Search project',
      itemUnit: '',
    },
    buttonTip: 'You need to select Jira projects to import data',
  },
  issueMap: {
    title: 'Jira issue mapping',
    tip: {
      message1:
        'Choose how to map JIRA issue types at 「ONES」 column. When the mapping is selected, the workflows and fields will be created accordingly.',
      message2:
        'Select“ONES custom issue type” to map the JIRA issuetype as a new ONES custom issue type.',
      message3:
        'Select Requirements, Tasks, Bug, Subtask to map the JIRA issue type as ONES existing system issue type.',
      message4:
        "JIRA issue type will only be mapped once in ONES and can't be changed afterwards. Please choose carefully.",
    },
    table: {
      columns: {
        jira: 'Jira',
        issueID: 'Jira issue type ID',
        onesIssue: 'ONES',
      },
      disabledTip: "Mapping relationship already established",
      placeholder: 'ONES custom issue type',
    },
  },
};

export default analyze;
