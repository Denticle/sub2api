export default {
  batchImageGuide: {
    title: '图片批量生成',
    description: '一次提交多条提示词，任务完成后可统一下载图片结果'
  },
  // Home Page
  home: {
    viewOnGithub: '在 GitHub 上查看',
    viewDocs: '查看文档',
    docs: '文档',
    switchToLight: '切换到浅色模式',
    switchToDark: '切换到深色模式',
    dashboard: '控制台',
    login: '登录',
    getStarted: '立即开始',
    goToDashboard: '进入控制台',
    // 新增：面向用户的价值主张
    heroSubtitle: '一个密钥，畅用多个 AI 模型',
    heroDescription: '无需管理多个订阅账号，一站式接入 Claude、GPT、Gemini 等主流 AI 服务',
    tags: {
      subscriptionToApi: '订阅转 API',
      stickySession: '会话保持',
      realtimeBilling: '按量计费'
    },
    // 用户痛点区块
    painPoints: {
      title: '你是否也遇到这些问题？',
      items: {
        expensive: {
          title: '订阅费用高',
          desc: '每个 AI 服务都要单独订阅，每月支出越来越多'
        },
        complex: {
          title: '多账号难管理',
          desc: '不同平台的账号、密钥分散各处，管理起来很麻烦'
        },
        unstable: {
          title: '服务不稳定',
          desc: '单一账号容易触发限制，影响正常使用'
        },
        noControl: {
          title: '用量无法控制',
          desc: '不知道钱花在哪了，也无法限制团队成员的使用'
        }
      }
    },
    // 解决方案区块
    solutions: {
      title: '我们帮你解决',
      subtitle: '简单三步，开始省心使用 AI'
    },
    features: {
      title: '平台能力',
      subtitle: '把多账号调度、计费与可观测性放在同一层',
      unifiedGateway: '一键接入',
      unifiedGatewayDesc: '获取一个 API 密钥，即可调用所有已接入的 AI 模型，无需分别申请。',
      multiAccount: '稳定可靠',
      multiAccountDesc: '智能调度多个上游账号，自动切换和负载均衡，告别频繁报错。',
      balanceQuota: '用多少付多少',
      balanceQuotaDesc: '按实际使用量计费，支持设置配额上限，团队用量一目了然。'
    },
    // 优势对比
    comparison: {
      title: '为什么选择我们？',
      headers: {
        feature: '对比项',
        official: '官方订阅',
        us: '本平台'
      },
      items: {
        pricing: {
          feature: '付费方式',
          official: '固定月费，用不完也付',
          us: '按量付费，用多少付多少'
        },
        models: {
          feature: '模型选择',
          official: '单一服务商',
          us: '多模型随意切换'
        },
        management: {
          feature: '账号管理',
          official: '每个服务单独管理',
          us: '统一密钥，一站管理'
        },
        stability: {
          feature: '服务稳定性',
          official: '单账号易触发限制',
          us: '多账号池，自动切换'
        },
        control: {
          feature: '用量控制',
          official: '无法限制',
          us: '可设配额、查明细'
        }
      }
    },
    hero: {
      eyebrow: '一个密钥，接入全部模型',
      description: '兼容 OpenAI / Anthropic / Gemini 协议，只需替换 Base URL，现有代码不用改。',
      compatible: '兼容 OpenAI SDK、Anthropic SDK、LangChain 与 Vercel AI SDK'
    },
    stats: {
      providers: '支持平台',
      protocols: '兼容协议',
      billing: '账单延迟',
      uptime: '服务可用性'
    },
    providers: {
      title: '支持的模型',
      description: '一个密钥调用多个上游，无需逐个申请账号',
      viewAll: '查看全部',
      supported: '已支持',
      soon: '即将推出',
      claude: 'Claude',
      gpt: 'GPT',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: '更多',
      claudeDesc: '长上下文推理与代码生成，适合复杂任务',
      gptDesc: '通用对话与工具调用，生态最完整',
      geminiDesc: '多模态理解，图片与长文档处理见长',
      antigravityDesc: '面向 Agent 场景的调用能力',
      deepseekDesc: '性价比高的推理与代码模型',
      grokDesc: '实时信息与推理能力',
      kimiDesc: '超长上下文，适合整库阅读',
      glmDesc: '中文理解与工具调用表现稳定',
      tagChat: '对话',
      tagMultimodal: '多模态',
      tagAgent: 'Agent',
      tagLongContext: '长上下文'
    },
    steps: {
      title: '三步接入',
      subtitle: '从注册到发出第一个请求，只需要几分钟',
      signup: '注册账号',
      signupDesc: '用邮箱注册即可获得初始额度，无需绑定信用卡。',
      credits: '充值额度',
      creditsDesc: '按需购买额度，用多少扣多少，余额随时可查。',
      key: '创建 API Key',
      keyDesc: '在控制台生成密钥，替换 Base URL 后即可直接调用。'
    },
    // CTA 区块
    cta: {
      title: '准备好开始了吗？',
      description: '注册即可获得免费试用额度，体验一站式 AI 服务',
      button: '免费注册'
    },

    // ── 重新设计的落地页（配色沿用 forapi token，结构参照 apipod.ai） ──
    lp: {
      status: '运行中',
      statusMeta: '{models} 个模型 · {providers} 个上游',
      headlineOne: '一个接口，接上',
      headlineTwo: '你需要的每个模型。',
      lede: '用同一个密钥调用 Claude、GPT 与 Gemini，按量计费，没有月费。改一行 base URL，现有代码原样跑通。',
      ledeNote: '上游账号在池内轮转，某个账号被限流时自动改道，请求不会因此失败。',
      ctaPrimary: '获取 API Key',
      ctaSecondary: '查看可用模型',
      traceRouted: '路由到',
      traceLimited: '被限流',
      traceFailover: '已改道',
      traceMs: '毫秒',
      statModels: '可用模型',
      statProviders: '上游平台',
      statRequests: '已处理请求',
      statBilling: '实时',
      statBillingValue: '<1s',
      statBillingLabel: '用量结算',
      statProtocols: '兼容 SDK',
      providersTitle: '一个密钥，接上全部上游',
      providersSub: '不必逐个申请账号。目录里的每个模型都已经接通，按平台分布如下。',
      modelsTitle: '模型与费率',
      modelsSub: '费率是相对上游官方的倍数，数字越小越省。',
      colModel: '模型',
      colPlatform: '上游',
      colRate: '费率',
      modelsAll: '查看全部模型',
      modelsEmpty: '模型目录暂不可用',
      capTitle: '这一层替你做四件事',
      capRouteTitle: '按负载选通道',
      capRouteDesc: '按并发、延迟与权重，把请求交给池里最合适的账号，而不是永远压在同一台上。',
      capFailTitle: '故障自己收敛',
      capFailDesc: '连续失败即熔断，冷却后再探活。坏掉的账号被隔在外面，流量改走健康账号。',
      capObserveTitle: '每次请求可追溯',
      capObserveDesc: '模型、状态、token、费用逐条落库，出问题时能直接定位到具体上游账号。',
      capBillTitle: '按 token 结算',
      capBillDesc: '没有订阅、没有最低消费。每一笔扣费都能对回具体那条请求。',
      compatTitle: '改一行就接入',
      compatSub: '协议与官方保持一致，SDK 不用换。',
      quickTitle: '三分钟拿到第一个响应',
      quickStep1: '注册并登录',
      quickStep1Desc: '邮箱注册即可，不必绑卡。',
      quickStep2: '创建 API Key',
      quickStep2Desc: '在控制台生成，随时可以吊销。',
      quickStep3: '替换 base URL',
      quickStep3Desc: '把请求指向本域名的 /v1 即可。',
      quickCopy: '复制',
      quickCopied: '已复制',
      quickDocs: '阅读文档',
      ctaTitle: '先跑通，再决定。',
      ctaDesc: '注册后在控制台创建密钥，把第一行请求发出来。',
      ctaButton: '获取 API Key',

      wbKicker: '工作台',
      wbTitle: '不写代码，也能直接用上这些模型',
      wbSub: '同一个账号里还有一套创作工具：在画布上把提示词、图片和视频串起来，或者单独生图、做视频。用你已经创建好的密钥，费用记在同一个余额里。',
      wbCta: '打开工作台',
      wbOpen: '打开',
      wbNodePrompt: '提示词',
      wbNodeImage: '图片',
      wbNodeVideo: '视频',
      wbMockPrompt: '一只博美犬，柔光棚拍，浅景深，暖色调',
      wbMockGenerating: '生成中…',
      wbCanvasTitle: '我的画布',
      wbCanvasDesc: '无限画布上放提示词、图片、视频节点，用连线把上一步的结果当作下一步的参考，一条链路从头做到尾。',
      wbImageTitle: '图片创作',
      wbImageDesc: '文生图与参考图编辑，一次出多张。每次生成都留记录，结果能下载，也能收进素材库。',
      wbVideoTitle: '视频创作',
      wbVideoDesc: '用文字或图片生成视频。任务在后台跑，切到别的页面也不会中断。',
      wbLibraryTitle: '提示词与素材',
      wbLibraryDesc: '提示词库里有现成的写法可以直接套用；生成的图片和视频收进「我的资产」，下次接着用。',
      wbFactKey: '用已有的密钥',
      wbFactKeyDesc: '在工作台里选一个你创建好的密钥即可开始，不用另外注册，也不用填接口地址。',
      wbFactBalance: '同一个余额',
      wbFactBalanceDesc: '工作台的调用和 API 调用走同一套计费，费率一致，都能在使用记录里查到。',
      wbFactLocal: '作品留在本地',
      wbFactLocalDesc: '画布、素材和生成记录保存在你的浏览器里；需要多台设备同步时，可以接上自己的 WebDAV。'
    },

    footer: {
      allRightsReserved: '保留所有权利。',
      keyUsage: '用量查询'
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key 用量查询',
    subtitle: '输入您的 API Key 以查看实时消费金额与使用状态',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: '查询',
    querying: '查询中...',
    privacyNote: '您的 Key 仅在浏览器本地处理，不会被存储',
    dateRange: '统计范围:',
    dateRangeToday: '今日',
    dateRange7d: '7 天',
    dateRange30d: '30 天',
    dateRange90d: '90 天',
    dateRangeCustom: '自定义',
    apply: '应用',
    used: '已使用',
    detailInfo: '详细信息',
    tokenStats: 'Token 统计',
    dailyDetail: '按日明细',
    modelStats: '模型用量统计',
    // Table headers
    date: '日期',
    model: '模型',
    requests: '请求数',
    inputTokens: '输入 Tokens',
    outputTokens: '输出 Tokens',
    cacheCreationTokens: '缓存创建',
    cacheReadTokens: '缓存读取',
    cacheWriteTokens: '缓存写入',
    totalTokens: '总 Tokens',
    cost: '费用',
    // Status
    quotaMode: 'Key 限额模式',
    walletBalance: '钱包余额',
    // Ring card titles
    totalQuota: '总额度',
    limit5h: '5 小时限额',
    limitDaily: '日限额',
    limit7d: '7 天限额',
    limitWeekly: '周限额',
    limitMonthly: '月限额',
    // Detail rows
    remainingQuota: '剩余额度',
    expiresAt: '过期时间',
    todayExpires: '(今日到期)',
    daysLeft: '({days} 天)',
    usedQuota: '已用额度',
    resetNow: '即将重置',
    subscriptionType: '订阅类型',
    billingType: '计费方式',
    subscriptionExpires: '订阅到期',
    // Usage stat cells
    todayRequests: '今日请求',
    todayInputTokens: '今日输入',
    todayOutputTokens: '今日输出',
    todayTokens: '今日 Tokens',
    todayCacheCreation: '今日缓存创建',
    todayCacheRead: '今日缓存读取',
    todayCost: '今日费用',
    rpmTpm: 'RPM / TPM',
    totalRequests: '累计请求',
    totalInputTokens: '累计输入',
    totalOutputTokens: '累计输出',
    totalTokensLabel: '累计 Tokens',
    totalCacheCreation: '累计缓存创建',
    totalCacheRead: '累计缓存读取',
    totalCost: '累计费用',
    avgDuration: '平均耗时',
    // Messages
    enterApiKey: '请输入 API Key',
    querySuccess: '查询成功',
    queryFailed: '查询失败',
    queryFailedRetry: '查询失败，请稍后重试',
    noDailyUsage: '暂无按日用量数据',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API 安装向导',
    description: '配置您的 Sub2API 实例',
    database: {
      title: '数据库配置',
      description: '连接到您的 PostgreSQL 数据库',
      host: '主机',
      port: '端口',
      username: '用户名',
      password: '密码',
      databaseName: '数据库名称',
      sslMode: 'SSL 模式',
      passwordPlaceholder: '密码',
      ssl: {
        disable: '禁用',
        require: '要求',
        verifyCa: '验证 CA',
        verifyFull: '完全验证'
      }
    },
    redis: {
      title: 'Redis 配置',
      description: '连接到您的 Redis 服务器',
      host: '主机',
      port: '端口',
      username: '用户名（可选）',
      password: '密码（可选）',
      database: '数据库',
      usernamePlaceholder: '默认用户留空',
      passwordPlaceholder: '密码',
      enableTls: '启用 TLS',
      enableTlsHint: '连接 Redis 时使用 TLS（公共 CA 证书）'
    },
    admin: {
      title: '管理员账户',
      description: '创建您的管理员账户',
      email: '邮箱',
      password: '密码',
      confirmPassword: '确认密码',
      passwordPlaceholder: '至少 8 个字符',
      confirmPasswordPlaceholder: '确认密码',
      passwordMismatch: '密码不匹配'
    },
    ready: {
      title: '准备安装',
      description: '检查您的配置并完成安装',
      database: '数据库',
      redis: 'Redis',
      adminEmail: '管理员邮箱'
    },
    status: {
      testing: '测试中...',
      success: '连接成功',
      testConnection: '测试连接',
      installing: '安装中...',
      completeInstallation: '完成安装',
      completed: '安装完成！',
      redirecting: '正在跳转到登录页面...',
      restarting: '服务正在重启，请稍候...',
      timeout: '服务重启时间超出预期，请手动刷新页面。'
    }
  },

  // Common
}
