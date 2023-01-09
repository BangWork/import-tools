const importPack = {
  initPassword: {
    title: '设置初始密码',
    tip: {
      message1:
        '1.请设置从 Jira 导入的新用户登录系统的初始密码。也可通过登录界面的忘记密码流程，验证邮箱重置密码登录；',
      message2: '2.此密码仅在当前页面设置和显示一次，请记录好你设置的初始登录密码。',
    },
    form: {
      tip: '系统将会在云端对导入数据进行处理。',
      init: {
        label: '初始密码',
        placeholder: '8-32位，且包含数字和字母',
        error: {
          empty: '初始密码不能为空',
          rule: '需包含 8-32 个字符，且包含数字和字母',
        },
      },
      again: {
        label: '确认密码',
        placeholder: '请再一次输入密码',
        error: {
          empty: '确认密码不能为空',
          diff: '两次输入的密码不同',
        },
      },
    },
    modal: {
      back: {
        title: '返回至「解析结果」',
        desc: '点击「返回」将返回至「解析结果」，相关信息需重新配置，是否确定？',
      },
    },
    startButton: '开始导入',
  },
  progressPage: {
    title: 'Jira 导入',
    logTitle: '迁移日志',
    timeMessage: '预计{{totalTime}}分钟，已导入{{leftTime}}分钟',
    tip: {
      environment: '导入环境：{{name}}',
      startTime: '开始导入时间：{{time}}',
      version: 'Jira 备份包版本：{{name}}',
      backUpTime: '备份时间：{{time}}',
    },
    action: {
      view: '查看导入范围',
      stop: '暂停导入',
      continue: '继续导入',
      cancel: '取消导入',
      closeLog: '关闭迁移日志',
      openLog: '查看迁移日志',
    },
    viewModal: {
      loading: '计算中',
      title: '导入范围',
      currentInfo: {
        title: '当前 Jira 导入数据范围',
        version: '版本',
        projects: '项目数',
        works: '工作项数',
        members: '成员数量',
        fileSize: '附件总大小',
        files: '附件数量',
      },
      unit: '个',
    },
    errorModal: {
      miss: {
        title: 'Jira 备份包不存在',
        text: 'Jira 备份包不存在，取消导入！',
      },
    },
    cancelModal: {
      title: '取消导入',
      desc: '部分数据已导入到 ONES 环境中，此操作不可撤销，是否确定取消导入？',
    },
  },
  importResultPage: {
    downloadText: '下载迁移日志',
    info: {
      title: '导入完成',
      desc: '从 Jira 导入到系统的新用户首次登录系统时使用管理员在导入数据时设置的初始密码；也可通过登录界面的忘记密码流程，验证邮箱重置密码登录。',
    },
    error: {
      title: '导入失败',
      desc1: 'Jira 导入工具运行错误，导入失败！',
      desc2: '未知错误，导入失败！',
    },
  },
};

export default importPack;
