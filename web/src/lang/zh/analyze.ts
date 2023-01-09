const analyze = {
  backupPage: {
    guide: {
      title: '获取 Jira 备份包路径',
      desc: '在 Jira 系统中找到 File Path（文件路径）以获取 Jira 备份包以便进行解析。',
      step1: {
        title: '点击 Jira 菜单栏 <1/> Administration/Settings （JIRA 管理）> System （系统）',
      },
      step2: {
        title: '点击 SYSTEM SUPPORT （系统支持） > System Info （系统信息）',
      },
      step3: {
        title:
          '在 System Info （系统信息）右侧下滑至 File Paths（文件路径） > Location of JIRA Local Home（JIRA 本地文件路径）， 复制该路径',
      },
    },
    form: {
      title: '填写 Jira 备份包信息',
      desc: '输入 Location of JIRA Local Home，点击“获取 Jira 备份包”，路径正确将可以选择导入的 Jira 备份包。',
      localHome: {
        label: 'Location of Jira Local Home（路径）',
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
      label: '部署 ONES 服务域名/IP',
      emptyError: '请输入域名/IP',
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
      network: {
        title: '访问失败！',
        desc: '无法访问 ONES 服务器，请重新填写 ONES 环境信息',
      },
      count: {
        title: '账号或密码错误！',
        desc: '请到 ONES 环境下验证正确的账号、密码后重新登陆该工具',
      },
      account: '账号或密码错误，请重新填写',
      team: '该 ONES 账号非团队管理员，请重新填写',
      organize: '该 ONES 账号非组织管理员，请重新填写',
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
      normalDesc: '「{{name}}」解析失败！请导入正确的Jira数据包。',
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
        desc: '该 ONES 团队于 {{time}} 导入 Jira {{version}} ，该 Jira {{packVersion}} 备份包为不同来源，禁止导入。请填写新的 Jira 备份包信息！',
      },
      importDiff: {
        title: 'ONES 团队已导入 Jira 数据',
        desc: '该ONES团队于 {{time}} 导入 Jira {{version}} ，是否确定导入该 Jira {{packVersion}} 备份包？',
      },
      importVersionDiff: {
        title: 'ONES 团队已导入 Jira 数据',
        desc: '该 ONES 团队于 {{time}} 日导入 Jira {{version}} ，该数据包为 Jira {{packVersion}} ，可能部分数据未能正常导入，是否确定？',
      },
      versionDiff: {
        title: 'Jira 版本不兼容',
        desc: '该 Jira 备份包为 Jira {{packVersion}} ，可能部分数据未能正常导入，是否确定？',
      },
    },
    backButton: '重新选择 Jira 备份包',
    buttonTip: '需选择导入的 ONES 团队',
  },
  importProject: {
    title: '选择导入的 Jira 项目',
    sourceTitle: 'Jira 项目列表',
    targetTitle: '导入 ONES 项目列表',
    local: {
      searchPlaceholder: '搜索项目名称',
      itemUnit: '项目',
    },
    buttonTip: '需选择导入的 Jira 项目',
  },
  issueMap: {
    title: 'Jira 工作项映射为 ONES 工作项',
    tip: {
      message1: '1、「 Jira 」列展示导入的 Jira 工作项类型',
      message2: '2、「 ONES 」列展示 Jira 工作项类型在 ONES 中的映射',
      message3: '3、单个映射关系只能建立一次，无法更改，请谨慎选择！',
    },
    table: {
      columns: {
        jira: 'Jira',
        issueID: 'Jira 工作项类型ID',
        onesIssue: 'ONES',
      },
      disabledTip: '该映射关系已确定，不可更改',
      placeholder: 'ONES 自定义工作项类型',
    },
  },
  shareDisk: {
    title: '是否使用共享磁盘导入数据',
    tip: {
      message1: '1、当使用共享磁盘方式后，Jira 附件会直接拷贝到指定位置，无需走 HTTP 方式导入',
      message2: '2、附件上传速度预计提升5倍，点击「使用」填写共享磁盘目录',
      message3:
        '3、注意：共享磁盘需要由运维人员将 ONES 文件存储目录(local_file_root)挂载到 Jira 服务器',
    },
    form: {
      label: '共享磁盘目录',
      emptyError: '共享磁盘目录不能为空',
      serverError: '共享磁盘目录错误，请重新输入',
    },
  },
};

export default analyze;
