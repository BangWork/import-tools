const analyze = {
  environment: {
    title: '登录 ONES',
    tip: {
      message1: '1. 为了确保你的正式环境下的 ONES 环境不会受到非必要的影响，建议你准备一个 ONES 测试环境，在此测试环境下完成一次数据迁移测试。验收通过后，在正式环境下进行正式的迁移。',
      message2: '2. 请向相关运维人员获取部署 ONES 服务域名/IP，仅 ONES 管理员可以执行 Jira 数据迁移。',
    },
    url: {
      label: 'ONES 服务域名/IP',
      emptyError: '请输入 ONES 服务域名/IP',
      serverError: '请输入正确的 ONES 服务域名/IP',
      placeholder: '例：http://ones.com 或 https://ones.com',
    },
    email: {
      label: 'ONES 邮箱',
      emptyError: '请输入邮箱',
    },
    password: {
      label: 'ONES 密码',
      emptyError: '密码不能为空',
    },
    serverError: {
      count: {
        title: '请输入正确的帐号或密码',
        desc: '请到 ONES 环境下验证正确的帐号和密码，10分钟后再重新登录此工具',
      },
      account: '请输入正确的帐号或密码',
      team: '此 ONES 帐号非团队管理员，请重新填写',
      organize: '此 ONES 帐号非组织管理员，请重新填写',

      version: {
        title: '请更新 ONES 环境版本',
        desc1: 'Jira 迁移工具不支持当前的 ONES 版本，请联系 ONES 迁移团队，协助安装合适的 ONES 版本。',
        desc2:'了解 Jira 迁移工具适用范围',
      }
    },
    accountError:'帐号或密码错误，你还有 {{count}} 次尝试机会',
    startButton: '开始解析',
    isLogin: {
      title: 'Log in to ONES',
      profile: '头像',
      ip: 'ONES domain/IP',
      email:'ONES account email',
    }
  },
  backupPage: {
    title:'选择 Jira 备份包',
    guide: {
      alert: {
        desc1: '1. 请确保你在 Jira 系统中已备份 Jira 数据',
        link1: '了解如何获取 Jira 备份包路径',
        desc2: '2. 请在 Jira 系统中获取 Jira 备份包路径，并完成以下信息配置。',
        link2: 'Find the location of the Jira home directory',
      },
      step1: {
        title: '获取备份包路径',
        desc1: '1.1 以 Jira 系统管理员身份登录 Jira 应用程序； ',
        desc2: '1.2 在 Jira 顶部导航栏中，点击 <1/> Jira administration > System；',
        desc3: '1.3 在左侧导航中，点击 SYSTEM SUPPORT > System Info。',
      },
      step2: {
        title: '填写 Jira 备份包信息',
        desc: '1.1 向下滚动到文件路径部分，将 File Paths > Location of JIRA Local Home 路径复制并粘贴到以下输入框中。',
        link: '查看默认 Jira 本地文件路径',
      },
    },
    form: {
      desc: '1.2 点击“获取 Jira 备份包”，路径正确将可以选择导入的 Jira 备份包。',
      localHome: {
        label: 'Location of Jira Local Home',
        get: '获取 Jira 备份包',
        emptyError: '必填项不能为空',
        serverError: '路径错误，请重新填写',
      },
      backup: {
        label: 'Jira 备份包名称',
        placeholder: '请选择导入的 Jira 备份包',
        emptyError: '请选择 Jira 备份包',
      },
      tip: '需选择导入的 Jira 备份包',
    },
  },

  analyzeProgress: {
    title: '解析 Jira 备份包',
    timeMessage: '预计{{totalTime}}分钟，剩余{{leftTime}}分钟',
    tip: {
      environment: '解析数据包将占用大量 Jira、ONES 服务器资源，建议你在低峰期开始解析。',

    },
    backupMessage: {
      title: 'Jira 备份包信息',
      status: {
        active: '解析中',
        success: '解析完成',
        fail: '解析失败',
      },
      analyzeProgress: '解析进度',
      analyzeBackupName: 'Jira backup name:',
      analyzeTime: '解析时间:',
      analyzeEnvironment: 'ONES environment:',
      analyzeFail:'解析失败，Jira 备份包数据格式错误，请重新上传'


    },
    analyzeResult: {
      title: '解析结果',
      jiraBackupResult: '2.1 Jirab备份包解析结果',
      onesTeamResult: '2.2 ONES团队信息解析结果',
      localStorage: '本地磁盘存储，磁盘容量大小: {{memory}} GB',
      localStorageSupport: 'ONES 服务器磁盘容量支持迁移',
      localStorageNotSupport: ' ONES 服务器磁盘容量不足全量导入，请扩充容量',
      localStorageRule: '  了解迁移磁盘规则',
    },
    tableTitle: {
      version: 'Jira 版本',
      projects: '项目数量',
      works: '工作项数量',
      members: '成员数量',
      fileSize: '附件总大小',
      files: '附件数量',
      id: 'Jira 服务器 ID',
      team: 'ONES 团队信息',
      status: '迁移状态',
      time: '迁移时间',

    },

  },
  analyzeResult: {
    title: '解析结果',
    current: {
      title: '当前 Jira 备份包信息',
      version: '版本',
      projects: '项目数量',
      works: '工作项数量',
      members: '成员数量',
      fileSize: '附件总大小',
      files: '附件数量',
      unit: '个',
    },
    environment: {
      title: 'ONES 环境信息',
      history: '{{time}} 导入 Jira {{version}}',
      empty: '未导入 Jira 数据',
    },
    modal: {
      back: {
        title: '返回至「填写 Jira 备份包信息」',
        desc: '点击「返回」将返回至第一步「填写 Jira 备份包信息」，相关信息需重新填写，是否确定？',
      },
    },
  },
  teamPage: {
    title: '选择 ONES 团队',
    desc:'建议选择未迁移过 Jira 数据的团队进行迁移。不建议重复迁移，可能会导致部分 Jira 数据未能正常迁移。',
    table: {
      teamName: 'ONES 团队名称',
      migrateStatus: '迁移状态',
      jiraBackupName: 'Jira 备份包名称',
      jiraVersion: 'Jira 版本',
      jiraId: 'Jira 服务器 ID',
      migrateTime: '迁移时间',
    },
    selectZero: '已选 0 个',
    selectTeam: '已选「{{teamName}}」团队',
    toSelectTeam: '请选择一个 ONES 团队',
    search: '搜索 ONES 团队',
    error: {
      packDiff: {
        title: 'Jira 备份包来源不同',
        desc: '此 ONES 团队于 {{time}} 导入 Jira {{version}} ，该 Jira {{packVersion}} 备份包为不同来源，禁止导入。请填写新的 Jira 备份包信息！',
      },
      importDiff: {
        title: 'ONES 团队已导入 Jira 数据',
        desc: '此 ONES团队于 {{time}} 导入 Jira {{version}} ，是否确定导入该 Jira {{packVersion}} 备份包？',
      },
      importVersionDiff: {
        title: 'ONES 团队已导入 Jira 数据',
        desc: '此 ONES 团队于 {{time}} 日导入 Jira {{version}} ，该数据包为 Jira {{packVersion}} ，可能部分数据未能正常导入，是否确定导入此 Jira {{packVersion}} 备份包？',
      },
      versionDiff: {
        title: 'Jira 版本不兼容',
        desc: '此 Jira 备份包为 Jira {{version}} ，可能部分数据未能正常导入，是否确定导入此 Jira {{packVersion}} 备份包？',
      },
      warning: {
        cancel: '继续导入',
        ok1: '重新选择团队',
        ok2: '重新选择 Jira 备份包',
      }
    },
    buttonTip: '需选择导入的 ONES 团队',

  },
  importProject: {
    title: '选择Jira 项目',
    search: '搜索项目名称、Key、负责人',
    noSupportMigrateProject: '不支持自助迁移的项目',
    desc1: '1. Jira 迁移工具仅支持迁移 Jira software 项目，即项目类型为“software”“business”的 Jira 项目。',
    desc2: '2. 不支持自助迁移的项目及其业务数据将不予迁移，如果你需要迁移此类数据，请咨询 ONES 迁移团队。',
    link: '联系我们',
    table: {
      projectName: '项目名称',
      projectKey: 'Key',
      leader: '负责人',
      projectClassification: '项目分类',
      issueCount: '问题数量',
    },
    selectZero: '已选 0 个',
    selectProject: '已选「{{projectName}}」项目',
    toSelectProject: '请选择一个 Jira 项目',

  },
  issueMap: {
    title: 'Jira 工作项映射为 ONES 工作项',
    tip: {
      message1: '「ONES」列支持选择如何映射 Jira 问题类型。选择映射关系后，对应问题工作流、属性将会对应创建。',
      message2: '选择「ONES 自定义工作项」则将该 Jira 问题类型创建为 ONES 新的自定义工作项类型。',
      message3: '选择需求、任务、缺陷、子任务，则将该 Jira 问题类型映射为 ONES 已有的系统工作项类型。',
      message4: '问题映射关系创建后将无法更改，请您谨慎选择！',
    },
    table: {
      columns: {
        jira: 'Jira',
        issueID: 'Jira 问题类型 ID',
        onesIssue: 'ONES',
      },
      disabledTip: '此映射关系已确定，不可更改',
      placeholder: 'ONES 自定义工作项类型',
    },
  },
};

export default analyze;
