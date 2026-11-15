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
      announceTag: '新',
      announceMore: '查看公告',
      announceClose: '关闭公告',
      navProduct: '产品',
      navPlaza: '模型广场',
      navGallery: '作品',
      navDocs: '文档',
      navLogin: '登录',
      navConsole: '控制台',
      navStart: '免费开始',
      menuStudio: '创作工作台',
      menuDev: '开发者',
      menuCanvas: '无限画布',
      menuCanvasDesc: '节点连线，一条链路做到底',
      menuAgent: 'Agent 创作',
      menuAgentDesc: '说一句要什么，自动拆步执行',
      menuDrama: '短剧与分镜',
      menuDramaDesc: '脚本拆镜头，出关键帧和视频',
      menuImage: '图片创作',
      menuImageDesc: '文生图、参考图编辑、批量出图',
      menuSkills: '技能库',
      menuSkillsDesc: '{n} 个现成工作流，选了就用',
      menuApi: 'API 接入',
      menuApiDesc: '兼容 OpenAI / Anthropic 协议',
      menuPlaza: '模型广场',
      menuPlazaDesc: '{n} 个模型的价格和能力',
      menuDocs: '接入文档',
      menuDocsDesc: 'base URL、密钥、示例代码',
      menuUsage: '用量与计费',
      menuUsageDesc: '按 token 结算，每笔可查',
      menuFeature: '看看别人做了什么',
      menuFeatureDesc: '灵感库里的作品，带提示词',
      status: '{models} 个模型 · {providers} 个上游平台 · 按 token 计费',
      headlineOne: '把提示词、图片',
      headlineTwo: '和分镜',
      headlineAccent: '连成一条线。',
      ledeA: '不用学节点。一句话就能出图、排镜；等你想要更精细的链路，画布上的',
      ledeHl: '每一步结果都是下一步的输入',
      ledeB: '。',
      ctaPrimary: '开始创作',
      ctaSecondary: '看看能做什么',
      note: '不用绑卡。注册就能进工作台试，满意了再充值。',
      cvDramaTitle: '短剧脚本',
      cvDramaBody: '霸总与灰姑娘的都市爱情',
      cvDramaLabel: 'AI 生成脚本 · 5 集',
      cvBoards: '分镜 · 79 张',
      cvDramaOut: '完整短剧 · 15min',
      cvAd: 'TVC 广告 · 30s',
      cvCodeTitle: '写个 React 组件',
      cvCodeBody: '响应式的卡片布局网格',
      cvCodeOut: 'React 组件 · JSX',
      cvBadge: '3 种场景 · 多条链路',
      statModels: '可调用模型',
      statProviders: '上游平台',
      statSkills: '内置技能',
      statBilling: '计价方式',
      billingValue: '按',
      wsTitle: '说一句要什么，Agent 在画布上把片子做出来',
      wsSub: '拆镜头、出关键帧、生成视频、配乐，每一步都落在画布上。哪一镜不满意，单独重做那一镜。',
      wsCrumbRoot: '我的画布',
      wsCrumbFile: '仙侠短片 · 第 3 集',
      wsRunning: 'Agent 运行中',
      wsUserMsg: '做一支 30 秒仙侠短片：神殿里对剑，被逼上断桥，御剑出海，最后和穷奇决战。电影感，16:9。',
      wsReply: '好，拆成 4 个镜头，按顺序做：',
      wsStep1: '拆分镜头脚本',
      wsStep1Meta: '4 镜 · 30s',
      wsStep2: '生成关键帧',
      wsStep2Meta: 'gpt-image · 4 张',
      wsStep3: '关键帧转视频',
      wsStep3Meta: 'Seedance 2.5 · 4 段',
      wsStep4: '配乐与音效',
      wsStep5: '合成导出',
      wsInput: '镜头 2 的桥再塌得猛一点…',
      wsShot: '镜头 {n} · {name}',
      wsShot1: '神殿对剑',
      wsShot2: '断桥突围',
      wsShot3: '御剑出海',
      wsShot4: '大战穷奇',
      wsPreview: '成片预览',
      wsNote: '这条链路本身就是一个技能。工作台里有 {n} 个这样的现成技能，短剧、广告、分镜、3D 预演，选一个就能开始。',
      waysKicker: '两种用法',
      waysTitle: '不管你要的是作品，还是接口',
      waysSub: '同一个账号，同一套账。想做东西从工作台进，想接程序从文档进。',
      waysMakeTag: '做作品',
      waysMakeTitle: '创作工作台',
      waysMakeDesc: '适合内容创作者、小团队、品牌方。不用写代码，选模型、填提示词，出图出片。',
      waysMake1: '无限画布 —— 把提示词、图片、分镜连成流水线',
      waysMake2: '图片创作 —— 单图、批量、参考图控制',
      waysMake3: '视频创作 —— 提示词直出视频，含分镜续写',
      waysMake4: '提示词库 —— 现成的成品提示词，改改就能用',
      waysMake5: '我的资产 —— 作品、素材、历史统一管理',
      waysMakeCta: '进入工作台',
      waysApiTag: '接接口',
      waysApiTitle: 'API 接入',
      waysApiDesc: '适合开发者与团队。改一行 base URL，SDK 不用换，现有代码原样跑通。',
      waysApi1: '兼容 OpenAI / Anthropic / Gemini 协议',
      waysApi2: '账号池自动轮转，单账号被限流自动改道',
      waysApi3: '每笔请求可追溯，token 用量逐条结算',
      waysApi4: '一把密钥绑定多个分组，按模型自动选组',
      waysApi5: '协议细节见接入文档',
      waysApiCta: '查看接入文档',
      galKicker: '作品与模型',
      galTitle: '这些都是在这里做出来的',
      galSub: '点任意一张，直接带着它的提示词进工作台改。',
      galAll: '看全部 {n} 件',
      galFilterAll: '全部',
      galFilterImage: '图片',
      galFilterVideo: '视频',
      galEmpty: '作品加载中，稍后再来看看。',
      clTitle: '今天就把第一支片子做出来',
      clDesc: '注册不用绑卡。工作台、画布和 API 用同一个账号、同一笔余额，按 token 计费，没有月费。',
      clPrimary: '免费开始创作',
      clSecondary: '联系销售',
      clPoint1: '注册即可进工作台',
      clPoint2: '余额按量扣，不用不扣',
      ftTagline: 'AI 创作工作台与模型 API，一个账号全部搞定。',
      ftStatus: '所有服务运行正常',
      ftProduct: '产品',
      ftDev: '开发者',
      ftCompany: '公司',
      ftLegal: '条款',
      ftStudio: '创作工作台',
      ftCanvas: '无限画布',
      ftAgent: 'Agent 创作',
      ftSkills: '技能库',
      ftGallery: '作品灵感',
      ftApi: 'API 接入',
      ftPlaza: '模型广场',
      ftDocs: '接入文档',
      ftUsage: '用量与计费',
      ftKeyUsage: '用量查询',
      ftAbout: '关于我们',
      ftContact: '联系销售',
      ftAnnouncements: '公告',
      ftChangelog: '更新日志',
      themeToLight: '切换到浅色模式',
      themeToDark: '切换到深色模式',
      mainNav: '主导航',
      footerNav: '页脚导航'
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
