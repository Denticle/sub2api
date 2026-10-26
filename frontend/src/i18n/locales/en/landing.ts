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
      status: 'Operational',
      statusMeta: '{models} models · {providers} providers',
      headlineOne: 'One endpoint for',
      headlineTwo: 'every model you need.',
      lede: 'Call Claude, GPT and Gemini with a single key. Pay for what you use, with no monthly fee. Change one base URL and your existing code keeps running.',
      ledeNote: 'Upstream accounts rotate inside a pool. When one gets rate limited, requests are rerouted instead of dropped.',
      ctaPrimary: 'Get an API key',
      ctaSecondary: 'Browse models',
      traceRouted: 'routed to',
      traceLimited: 'rate limited',
      traceFailover: 'failed over',
      traceMs: 'ms',
      statModels: 'Models',
      statProviders: 'Providers',
      statRequests: 'Requests served',
      statBilling: 'Real time',
      statBillingValue: '<1s',
      statBillingLabel: 'Usage billing',
      statProtocols: 'Compatible SDKs',
      providersTitle: 'One key, every upstream',
      providersSub: 'Every model in the catalog is already connected. Here is how they break down by provider.',
      modelsTitle: 'Models and rates',
      modelsSub: 'Rates are a multiple of the upstream list price. Lower is cheaper.',
      colModel: 'Model',
      colPlatform: 'Provider',
      colRate: 'Rate',
      modelsAll: 'See all models',
      modelsEmpty: 'Model catalog unavailable',
      capTitle: 'Four things this layer handles',
      capRouteTitle: 'Route by load',
      capRouteDesc: 'Concurrency, latency and weight decide which account in the pool takes the request, instead of pinning everything to one machine.',
      capFailTitle: 'Failures contain themselves',
      capFailDesc: 'Consecutive failures trip a breaker. After a cooldown the account gets probed again. Bad accounts are isolated, traffic goes to healthy ones.',
      capObserveTitle: 'Every request is traceable',
      capObserveDesc: 'Model, status, tokens and cost land in the database per request, so an incident traces back to the exact upstream account.',
      capBillTitle: 'Billed by token',
      capBillDesc: 'No subscription, no minimum. Every charge reconciles to the individual request behind it.',
      compatTitle: 'Change one line',
      compatSub: 'The protocol matches upstream, so your SDK stays as it is.',
      quickTitle: 'First response in three minutes',
      quickStep1: 'Sign up',
      quickStep1Desc: 'Email is enough. No card required.',
      quickStep2: 'Create a key',
      quickStep2Desc: 'Generate it in the console and revoke it any time.',
      quickStep3: 'Swap the base URL',
      quickStep3Desc: 'Point requests at this domain under /v1.',
      quickCopy: 'Copy',
      quickCopied: 'Copied',
      quickDocs: 'Read the docs',
      ctaTitle: 'Make it work, then decide.',
      ctaDesc: 'Sign up, create a key in the console, and send the first request.',
      ctaButton: 'Get an API key'
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
