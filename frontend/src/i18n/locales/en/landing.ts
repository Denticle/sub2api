export default {
  batchImageGuide: {
    title: 'Batch Image Generation',
    description: 'Submit multiple prompts in one job and download the generated images when complete'
  },
  // Home Page
  home: {
    viewOnGithub: 'View on GitHub',
    viewDocs: 'View Documentation',
    docs: 'Docs',
    switchToLight: 'Switch to Light Mode',
    switchToDark: 'Switch to Dark Mode',
    dashboard: 'Dashboard',
    login: 'Login',
    getStarted: 'Get Started',
    goToDashboard: 'Go to Dashboard',
    // User-focused value proposition
    heroSubtitle: 'One Key, All AI Models',
    heroDescription: 'No need to manage multiple subscriptions. Access Claude, GPT, Gemini and more with a single API key',
    tags: {
      subscriptionToApi: 'Subscription to API',
      stickySession: 'Session Persistence',
      realtimeBilling: 'Pay As You Go'
    },
    // Pain points section
    painPoints: {
      title: 'Sound Familiar?',
      items: {
        expensive: {
          title: 'High Subscription Costs',
          desc: 'Paying for multiple AI subscriptions that add up every month'
        },
        complex: {
          title: 'Account Chaos',
          desc: 'Managing scattered accounts and API keys across different platforms'
        },
        unstable: {
          title: 'Service Interruptions',
          desc: 'Single accounts hitting rate limits and disrupting your workflow'
        },
        noControl: {
          title: 'No Usage Control',
          desc: "Can't track where your money goes or limit team member usage"
        }
      }
    },
    // Solutions section
    solutions: {
      title: 'We Solve These Problems',
      subtitle: 'Three simple steps to stress-free AI access'
    },
    features: {
      title: 'Platform capabilities',
      subtitle: 'Routing, billing and observability in one layer',
      unifiedGateway: 'One-Click Access',
      unifiedGatewayDesc: 'Get a single API key to call all connected AI models. No separate applications needed.',
      multiAccount: 'Always Reliable',
      multiAccountDesc: 'Smart routing across multiple upstream accounts with automatic failover. Say goodbye to errors.',
      balanceQuota: 'Pay What You Use',
      balanceQuotaDesc: 'Usage-based billing with quota limits. Full visibility into team consumption.'
    },
    // Comparison section
    comparison: {
      title: 'Why Choose Us?',
      headers: {
        feature: 'Comparison',
        official: 'Official Subscriptions',
        us: 'Our Platform'
      },
      items: {
        pricing: {
          feature: 'Pricing',
          official: 'Fixed monthly fee, pay even if unused',
          us: 'Pay only for what you use'
        },
        models: {
          feature: 'Model Selection',
          official: 'Single provider only',
          us: 'Switch between models freely'
        },
        management: {
          feature: 'Account Management',
          official: 'Manage each service separately',
          us: 'Unified key, one dashboard'
        },
        stability: {
          feature: 'Stability',
          official: 'Single account rate limits',
          us: 'Multi-account pool, auto-failover'
        },
        control: {
          feature: 'Usage Control',
          official: 'Not available',
          us: 'Quotas & detailed analytics'
        }
      }
    },
    hero: {
      eyebrow: 'One key for every model',
      description: 'Compatible with OpenAI, Anthropic and Gemini protocols. Swap the base URL, keep your code.',
      compatible: 'Works with OpenAI SDK, Anthropic SDK, LangChain and Vercel AI SDK'
    },
    stats: {
      providers: 'Providers',
      protocols: 'Protocols',
      billing: 'Billing lag',
      uptime: 'Uptime'
    },
    providers: {
      title: 'Supported Models',
      description: 'Call several upstreams with one key, no separate accounts',
      viewAll: 'View all',
      supported: 'Supported',
      soon: 'Soon',
      claude: 'Claude',
      gpt: 'GPT',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: 'More',
      claudeDesc: 'Long-context reasoning and code generation',
      gptDesc: 'General chat and tool calling, broadest ecosystem',
      geminiDesc: 'Multimodal understanding for images and long docs',
      antigravityDesc: 'Agent-oriented invocation capabilities',
      deepseekDesc: 'Cost-effective reasoning and code models',
      grokDesc: 'Real-time information and reasoning',
      kimiDesc: 'Very long context for whole-repo reading',
      glmDesc: 'Solid Chinese understanding and tool calling',
      tagChat: 'Chat',
      tagMultimodal: 'Multimodal',
      tagAgent: 'Agent',
      tagLongContext: 'Long context'
    },
    steps: {
      title: 'How it works',
      subtitle: 'From signup to your first request in a few minutes',
      signup: 'Sign up',
      signupDesc: 'Create an account with your email and get starting credits. No card required.',
      credits: 'Add credits',
      creditsDesc: 'Buy credits as you need them. Pay only for what you use, check balance anytime.',
      key: 'Create an API key',
      keyDesc: 'Generate a key in the console, swap the base URL, and start calling.'
    },
    // CTA section
    cta: {
      title: 'Ready to Get Started?',
      description: 'Sign up now and get free trial credits to experience seamless AI access',
      button: 'Sign Up Free'
    },

    // ── Redesigned landing page (forapi tokens, structure after apipod.ai) ──
    lp: {
      announceTag: 'New',
      announceMore: 'Read more',
      announceClose: 'Dismiss announcement',
      navProduct: 'Product',
      navPlaza: 'Models',
      navGallery: 'Showcase',
      navDocs: 'Docs',
      navLogin: 'Sign in',
      navConsole: 'Console',
      navStart: 'Get started',
      menuStudio: 'Creative studio',
      menuDev: 'Developers',
      menuCanvas: 'Infinite canvas',
      menuCanvasDesc: 'Wire nodes into one pipeline',
      menuAgent: 'Agent',
      menuAgentDesc: 'Describe it, the agent plans the steps',
      menuDrama: 'Shorts & storyboards',
      menuDramaDesc: 'Script to shots, keyframes and video',
      menuImage: 'Images',
      menuImageDesc: 'Text-to-image, reference edits, batches',
      menuSkills: 'Skills',
      menuSkillsDesc: '{n} ready-made workflows',
      menuApi: 'API',
      menuApiDesc: 'OpenAI / Anthropic compatible',
      menuPlaza: 'Model catalog',
      menuPlazaDesc: 'Prices and capabilities of {n} models',
      menuDocs: 'Docs',
      menuDocsDesc: 'Base URL, keys and code samples',
      menuUsage: 'Usage & billing',
      menuUsageDesc: 'Per-token billing, every call logged',
      menuFeature: 'See what others made',
      menuFeatureDesc: 'Showcase pieces with their prompts',
      status: '{models} models · {providers} providers · per-token billing',
      headlineOne: 'Prompts, images',
      headlineTwo: 'and storyboards,',
      headlineAccent: 'in one line.',
      ledeA: 'No nodes to learn. One sentence gets you images and shots; when you want finer control, ',
      ledeHl: 'every result on the canvas feeds the next step',
      ledeB: '.',
      ctaPrimary: 'Start creating',
      ctaSecondary: 'See what it makes',
      note: 'No card required. Try the studio first, top up when you like it.',
      cvDramaTitle: 'Short drama',
      cvDramaBody: 'A city romance in five episodes',
      cvDramaLabel: 'AI script · 5 episodes',
      cvBoards: 'Storyboard · 79',
      cvDramaOut: 'Full drama · 15min',
      cvAd: 'TVC ad · 30s',
      cvCodeTitle: 'Write a React component',
      cvCodeBody: 'A responsive card grid',
      cvCodeOut: 'React component · JSX',
      cvBadge: '3 scenarios · many pipelines',
      statModels: 'Models',
      statProviders: 'Providers',
      statSkills: 'Built-in skills',
      statBilling: 'Billing',
      billingValue: 'Per',
      wsTitle: 'Say what you want. The agent builds it on the canvas.',
      wsSub: 'Shot list, keyframes, video, music — every step lands on the canvas. Redo just the shot you don\'t like.',
      wsCrumbRoot: 'My canvas',
      wsCrumbFile: 'Xianxia short · Ep. 3',
      wsRunning: 'Agent running',
      wsUserMsg: 'Make a 30s xianxia short: a sword duel in a temple, a collapsing bridge, flying out over the sea, then the final fight with the beast. Cinematic, 16:9.',
      wsReply: 'Got it — four shots, in order:',
      wsStep1: 'Write shot list',
      wsStep1Meta: '4 shots · 30s',
      wsStep2: 'Generate keyframes',
      wsStep2Meta: 'gpt-image · 4',
      wsStep3: 'Keyframes to video',
      wsStep3Meta: 'Seedance 2.5 · 4',
      wsStep4: 'Music & sound',
      wsStep5: 'Export',
      wsInput: 'Make the bridge in shot 2 collapse harder…',
      wsShot: 'Shot {n} · {name}',
      wsShot1: 'Temple duel',
      wsShot2: 'Bridge escape',
      wsShot3: 'Out to sea',
      wsShot4: 'The beast',
      wsPreview: 'Preview',
      wsNote: 'This pipeline is itself a skill. The studio ships {n} of them — shorts, ads, storyboards, 3D previs.',
      waysKicker: 'Two ways in',
      waysTitle: 'Whether you want a finished piece or an API',
      waysSub: 'One account, one balance. Make things in the studio, wire up code with the docs.',
      waysMakeTag: 'Create',
      waysMakeTitle: 'Creative studio',
      waysMakeDesc: 'For creators, small teams and brands. No code — pick a model, write a prompt, get images and video.',
      waysMake1: 'Infinite canvas — wire prompts, images and shots into a pipeline',
      waysMake2: 'Images — single, batch, reference control',
      waysMake3: 'Video — prompt to video, with shot continuation',
      waysMake4: 'Prompt library — finished prompts, ready to tweak',
      waysMake5: 'Assets — work, media and history in one place',
      waysMakeCta: 'Open the studio',
      waysApiTag: 'Integrate',
      waysApiTitle: 'API',
      waysApiDesc: 'For developers and teams. Change one base URL, keep your SDK, your code runs as is.',
      waysApi1: 'OpenAI / Anthropic / Gemini compatible',
      waysApi2: 'Account pool rotation, automatic rerouting on rate limits',
      waysApi3: 'Every request traceable, billed per token',
      waysApi4: 'One key across groups, routed by model',
      waysApi5: 'Protocol details in the docs',
      waysApiCta: 'Read the docs',
      galKicker: 'Showcase',
      galTitle: 'All of this was made here',
      galSub: 'Pick any piece to open it in the studio with its prompt.',
      galAll: 'See all {n}',
      galFilterAll: 'All',
      galFilterImage: 'Images',
      galFilterVideo: 'Video',
      galEmpty: 'The showcase is loading — check back soon.',
      clTitle: 'Make your first piece today',
      clDesc: 'No card to sign up. Studio, canvas and API share one account and one balance, billed per token, no subscription.',
      clPrimary: 'Start creating free',
      clSecondary: 'Contact sales',
      clPoint1: 'Studio access on sign-up',
      clPoint2: 'Pay only for what you use',
      ftTagline: 'AI creative studio and model API, one account.',
      ftStatus: 'All systems operational',
      ftProduct: 'Product',
      ftDev: 'Developers',
      ftCompany: 'Company',
      ftLegal: 'Legal',
      ftStudio: 'Creative studio',
      ftCanvas: 'Infinite canvas',
      ftAgent: 'Agent',
      ftSkills: 'Skills',
      ftGallery: 'Showcase',
      ftApi: 'API',
      ftPlaza: 'Model catalog',
      ftDocs: 'Docs',
      ftUsage: 'Usage & billing',
      ftKeyUsage: 'Key usage lookup',
      ftAbout: 'About',
      ftContact: 'Contact sales',
      ftAnnouncements: 'Announcements',
      ftChangelog: 'Changelog',
      themeToLight: 'Switch to light mode',
      themeToDark: 'Switch to dark mode',
      mainNav: 'Main navigation',
      footerNav: 'Footer navigation'
    },
    footer: {
      allRightsReserved: 'All rights reserved.',
      keyUsage: 'Key usage'
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key Usage',
    subtitle: 'Enter your API Key to view real-time spending and usage status',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: 'Query',
    querying: 'Querying...',
    privacyNote: 'Your Key is processed locally in the browser and will not be stored',
    dateRange: 'Date Range:',
    dateRangeToday: 'Today',
    dateRange7d: '7 Days',
    dateRange30d: '30 Days',
    dateRange90d: '90 Days',
    dateRangeCustom: 'Custom',
    apply: 'Apply',
    used: 'Used',
    detailInfo: 'Detail Information',
    tokenStats: 'Token Statistics',
    dailyDetail: 'Daily Detail',
    modelStats: 'Model Usage Statistics',
    // Table headers
    date: 'Date',
    model: 'Model',
    requests: 'Requests',
    inputTokens: 'Input Tokens',
    outputTokens: 'Output Tokens',
    cacheCreationTokens: 'Cache Creation',
    cacheReadTokens: 'Cache Read',
    cacheWriteTokens: 'Cache Write',
    totalTokens: 'Total Tokens',
    cost: 'Cost',
    // Status
    quotaMode: 'Key Quota Mode',
    walletBalance: 'Wallet Balance',
    // Ring card titles
    totalQuota: 'Total Quota',
    limit5h: '5-Hour Limit',
    limitDaily: 'Daily Limit',
    limit7d: '7-Day Limit',
    limitWeekly: 'Weekly Limit',
    limitMonthly: 'Monthly Limit',
    // Detail rows
    remainingQuota: 'Remaining Quota',
    expiresAt: 'Expires At',
    todayExpires: '(expires today)',
    daysLeft: '({days} days)',
    usedQuota: 'Used Quota',
    resetNow: 'Resetting soon',
    subscriptionType: 'Subscription Type',
    billingType: 'Billing Type',
    subscriptionExpires: 'Subscription Expires',
    // Usage stat cells
    todayRequests: 'Today Requests',
    todayInputTokens: 'Today Input',
    todayOutputTokens: 'Today Output',
    todayTokens: 'Today Tokens',
    todayCacheCreation: 'Today Cache Creation',
    todayCacheRead: 'Today Cache Read',
    todayCost: 'Today Cost',
    rpmTpm: 'RPM / TPM',
    totalRequests: 'Total Requests',
    totalInputTokens: 'Total Input',
    totalOutputTokens: 'Total Output',
    totalTokensLabel: 'Total Tokens',
    totalCacheCreation: 'Total Cache Creation',
    totalCacheRead: 'Total Cache Read',
    totalCost: 'Total Cost',
    avgDuration: 'Avg Duration',
    // Messages
    enterApiKey: 'Please enter an API Key',
    querySuccess: 'Query successful',
    queryFailed: 'Query failed',
    queryFailedRetry: 'Query failed, please try again later',
    noDailyUsage: 'No daily usage data',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API Setup',
    description: 'Configure your Sub2API instance',
    database: {
      title: 'Database Configuration',
      description: 'Connect to your PostgreSQL database',
      host: 'Host',
      port: 'Port',
      username: 'Username',
      password: 'Password',
      databaseName: 'Database Name',
      sslMode: 'SSL Mode',
      passwordPlaceholder: 'Password',
      ssl: {
        disable: 'Disable',
        require: 'Require',
        verifyCa: 'Verify CA',
        verifyFull: 'Verify Full'
      }
    },
    redis: {
      title: 'Redis Configuration',
      description: 'Connect to your Redis server',
      host: 'Host',
      port: 'Port',
      username: 'Username (optional)',
      password: 'Password (optional)',
      database: 'Database',
      usernamePlaceholder: 'Leave empty for default user',
      passwordPlaceholder: 'Password',
      enableTls: 'Enable TLS',
      enableTlsHint: 'Use TLS when connecting to Redis (public CA certs)'
    },
    admin: {
      title: 'Admin Account',
      description: 'Create your administrator account',
      email: 'Email',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      passwordPlaceholder: 'Min 8 characters',
      confirmPasswordPlaceholder: 'Confirm password',
      passwordMismatch: 'Passwords do not match'
    },
    ready: {
      title: 'Ready to Install',
      description: 'Review your configuration and complete setup',
      database: 'Database',
      redis: 'Redis',
      adminEmail: 'Admin Email'
    },
    status: {
      testing: 'Testing...',
      success: 'Connection Successful',
      testConnection: 'Test Connection',
      installing: 'Installing...',
      completeInstallation: 'Complete Installation',
      completed: 'Installation completed!',
      redirecting: 'Redirecting to login page...',
      restarting: 'Service is restarting, please wait...',
      timeout: 'Service restart is taking longer than expected. Please refresh the page manually.'
    }
  },

  // Common
}
