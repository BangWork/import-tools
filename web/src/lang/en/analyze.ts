const analyze = {
  backupPage: {
    guide: {
      title: 'Find Jira backup path',
      desc: 'Jira file paths are required to import and analyze backups from Jira.',
      step1: {
        title: 'Click the <1/> icon in the top menu bar, then click System',
      },
      step2: {
        title: 'Navigate to SYSTEM SUPPORT -> System Info',
      },
      step3: {
        title:
          'Scroll down and navigate to File Path -> Location of JIRA Local Home. Copy the path shown on the right.',
      },
    },
    form: {
      title: 'Provide Jira backup info',
      desc: 'Input the path you copied from Location of JIRA Local Home, then click "Get Jira backups" and select the backups you want to import.',
      localHome: {
        label: 'Location of Jira Local Home',
        get: 'Get Jira backup',
        emptyError: "Required fields can't be empty",
        serverError: 'Invalid directory. Please try again',
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
      message1: '1. Get ONES domain / IP from DevOPS',
      message2: '2. Only admin can perform this action',
    },
    url: {
      label: 'ONES domain / IP',
      emptyError: 'Please enter domain/IP',
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
      network: {
        title: 'Access denied',
        desc: 'Unable to access server. Please make sure you provided the correct ONES environment info',
      },
      count: {
        title: 'Wrong account or password.',
        desc: 'Please log into ONES with a valid account and password',
      },
      account: 'Wrong account or password. Please try again',
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
      doing: 'Analyzing'
    },
    fail: {
      title: 'Analysis failed',
      normalDesc: 'Failed to analysis [{{name}}]. Please make sure you provided the correct Jira backup info.',
      onExistDesc: 'Analysis failed. [{{name}}] does not exist',
    },
  },
  analyzeResult: {
    title: 'Results',
    current: {
      title: 'Jira backup info',
      version: 'Version',
      projects: 'Number of projects',
      works: 'Number of issues',
      members: 'Number of members',
      fileSize: 'File size',
      files: 'Number of files',
      unit: '',
    },
    environment: {
      title: 'Target ONES environment info',
      history: 'Imported data from Jira {{version}} on {{time}}',
      empty: 'No data imported from Jira'
    },
    modal: {
      back: {
        title: 'Return to [Provide Jira backup info]',
        desc: 'After selecting [Back], you will be returned to the [Provide Jira backup info] page and need to provide all the information again. Proceed?',
      }
    }
  },
  teamPage: {
    title: 'Select and import data into a ONES team',
    form: {
      label: 'Import data into a ONES team',
      placeholder: 'Select a ONES team',
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
        desc: "Jira {{packVersion}} will be imported into the selected ONES team. Some data may not be properly imported. Proceed? This team has already selected and imported data from Jira {{version}} on {{time}}.",
      },
      versionDiff: {
        title: 'Incompatible Jira version',
        desc: 'Jira {{packVersion}} will be imported into the selected ONES team. Some data may not be properly imported. Proceed?',
      },
    },
    backButton: 'Select Jira backup again',
    buttonTip: 'Please select a ONES team',
  },
  importProject: {
    title: 'Select which Jira project you want to import data from',
    sourceTitle: 'Jira projects',
    targetTitle: 'Projects imported to ONES',
    local: {
      searchPlaceholder: 'Search project',
      itemUnit: 'Project',
    },
    buttonTip: 'You need to select Jira projects to import data',
  },
  issueMap: {
    title: 'Jira issue mapping',
    tip: {
      message1: '1、Jira issue types will be displayed in the [Jira] section',
      message2: '2、The [ONES] section displays how issue types from Jira will be mapped on ONES',
      message3: "3、Each mapping relationship can only be established once and can't be changed. Please proceed with caution.",
    },
    table: {
      columns: {
        jira: 'Jira',
        issueID: 'Jira issue type ID',
        onesIssue: 'ONES',
      },
      disabledTip: "This mapping relationship is already established and can't be changed",
      placeholder: 'ONES custom issue type',
    },
  },
  shareDisk: {
    title: 'Do you want to import data via shared disks?',
    tip: {
      message1:
        '1、When importing data via shared disks, Jira data and attachments will be copied to a designated location. HTTP method will no longer be required',
      message2: '2、Import speed of shared disks is 5 times faster. If you select [Yes], you will need to provide shared disk paths',
      message3:
        '3、Note: In order to import data via shared disks, DevOPS team needs to mount ONES file storage path (local_file_root) on to the Jira server',
    },
    form: {
      label: 'Shared disk path',
      emptyError: 'Shared disk path is required',
      serverError: 'Wrong path. Please try again',
    },
  },
};

export default analyze;
