const analyze = {
  backupPage: {
    guide: {
      alert: {
        desc: '请在 Jira 系统中获取 Jira 备份包路径，并完成以下信息配置。',
        link: '了解如何获取 Jira 备份包路径',
      },
      step1: {
        title: '获取备份包路径',
        desc1: '1.1 以 Jira 系统管理员身份登录 Jira 应用程序； ',
        desc2: '1.2 在 Jira 顶部导航栏中，点击 <1/> Jira administration > System；',
        desc3: '1.3 在左侧导航中，点击 SYSTEM SUPPORT > System Info。',
      },
      step2: {
        title: '填写 Jira 备份包信息',
        desc: '向下滚动到文件路径部分，将 File Paths > Location of JIRA Local Home 路径复制并粘贴到以下输入框中。',
        link: '查看默认 Jira 本地文件路径',
      },
    },
    form: {
      desc: '输入 Location of JIRA Local Home，点击“获取 Jira 备份包”，路径正确将可以选择导入的 Jira 备份包。',
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
  environment: {
    title: '填写导入的 ONES 环境信息',
    tip: {
      message1: '1、请向相关运维人员获取部署 ONES 服务域名/IP',
      message2: '2、导入者账号需为管理员',
    },
    url: {
      label: 'ONES 服务域名/IP',
      emptyError: '请输入 ONES 服务域名/IP',
      serverError: '请输入正确的 ONES 服务域名/IP',
      placeholder: '例: http://ones.cn 或 https://ones.cn',
    },
    email: {
      label: 'ONES 邮箱（导入者）',
      emptyError: '请输入邮箱',
    },
    password: {
      label: 'ONES 密码（导入者）',
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
    },
    startButton: '开始解析',
  },
  analyzeProgress: {
    title: '解析 Jira 备份包',
    timeMessage: '预计{{totalTime}}分钟，已解析{{leftTime}}分钟',
    tip: {
      environment: '1、导入环境：{{name}}',
      time: '2、开始解析时间：{{time}}',
    },
    cancel: {
      text: '取消解析',
      success: '取消解析成功',
      fail: '取消解析失败',
      loading: '取消解析中...',
      desc: '取消解析将返回「填写 Jira 备份包信息」页面，是否确定？',
    },
    status: {
      doing: '解析中',
    },
    fail: {
      title: '解析失败！',
      normalDesc: '「{{name}}」解析失败！请导入正确的 Jira 数据包。',
      onExistDesc: '「{{name}}」不存在，解析失败！',
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
    title: '选择导入的 ONES 团队',
    form: {
      label: '导入的 ONES 团队',
      placeholder: '选择 ONES 团队',
    },
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
    title: '选择导入的 Jira 项目',
    sourceTitle: 'Jira 项目列表',
    targetTitle: '导入 ONES 项目列表',
    local: {
      searchPlaceholder: '搜索项目名称',
      itemUnit: '个',
    },
    buttonTip: '需选择导入的 Jira 项目',
  },
  issueMap: {
    title: '映射 Jira 问题类型',
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
