/**
 * /studio/ 的接入层。
 *
 * 无限画布本身是个纯前端应用：AI 请求由浏览器直连，API Key 存在自己的
 * localStorage（infinite-canvas:ai_config_store）里，让用户自己粘。
 *
 * 这个脚本在应用 bundle 之前执行，做三件事：
 *   1. 登录守卫：没有 sub2api 的 auth_token 就直接跳回 /login
 *   2. 选 key：列出用户已经创建好的密钥，选一个，而不是让用户自己填
 *   3. 注入配置：把选中的 key + 模型列表写进画布的 store，再把应用加载起来
 *
 * 它必须同步地拦住应用加载——所以 index.html 里应用脚本标了 data-studio-app，
 * 这里先把它摘下来，等配置就绪后自己注入回去。
 */
(function () {
  'use strict'

  var AUTH_TOKEN_KEY = 'auth_token'
  // 构建后应用脚本是 data-studio-app-src；兼容旧版 data-studio-app 标记
  var APP_SCRIPT_SELECTOR = 'script[data-studio-app-src], script[data-studio-app]'
  var CONFIG_STORE_KEY = 'infinite-canvas:ai_config_store'
  var THEME_STORE_KEY = 'infinite-canvas:theme_store'
  var API_BASE = '/api/v1'
  var CHANNEL_ID = 'forapi'
  var CHANNEL_NAME = 'ForAPI'

  // ---- 模型能力推断：与画布 use-config-store.ts 里的 guessCapability 同序同表 ----
  var VIDEO_KEYWORDS = ['video', 'sora', 'veo', 'kling', 'wan', 'hailuo']
  var AUDIO_KEYWORDS = ['audio', 'tts', 'speech', 'voice', 'music', 'sound']
  var IMAGE_KEYWORDS = [
    'seedream', 'gpt-image', 'image', 'dall-e', 'dalle',
    'imagen', 'flux', 'sdxl', 'stable-diffusion', 'midjourney',
  ]

  function guessCapability(name) {
    var v = String(name || '').toLowerCase()
    function has(list) {
      for (var i = 0; i < list.length; i++) if (v.indexOf(list[i]) !== -1) return true
      return false
    }
    if (has(VIDEO_KEYWORDS)) return 'video'
    if (has(AUDIO_KEYWORDS)) return 'audio'
    if (has(IMAGE_KEYWORDS)) return 'image'
    return 'text'
  }

  // ---- 存储读写 ----
  function readJSON(key) {
    try {
      return JSON.parse(localStorage.getItem(key) || 'null')
    } catch (e) {
      return null
    }
  }

  function getToken() {
    try {
      return localStorage.getItem(AUTH_TOKEN_KEY) || ''
    } catch (e) {
      return ''
    }
  }

  /** 画布里当前配置的 key，用于判断「已经接好了，不用再选」 */
  function currentKey() {
    var raw = readJSON(CONFIG_STORE_KEY)
    var channels = (raw && raw.state && raw.state.config && raw.state.config.channels) || []
    for (var i = 0; i < channels.length; i++) {
      if (channels[i] && channels[i].id === CHANNEL_ID) return channels[i].apiKey || ''
    }
    return ''
  }

  function writeConfig(apiKey, models) {
    var modelOptions = models.map(function (m) {
      return CHANNEL_ID + '::' + m.name
    })
    var firstOf = function (cap) {
      for (var i = 0; i < models.length; i++) if (models[i].capability === cap) return CHANNEL_ID + '::' + models[i].name
      return ''
    }
    var image = firstOf('image')
    var text = firstOf('text')
    var video = firstOf('video')
    var audio = firstOf('audio')

    var pick = function (primary, fallback) {
      return primary || fallback || text || (modelOptions[0] || '')
    }

    var config = {
      channelMode: 'local',
      baseUrl: location.origin,
      apiKey: apiKey,
      apiFormat: 'openai',
      channels: [
        {
          id: CHANNEL_ID,
          name: CHANNEL_NAME,
          baseUrl: location.origin,
          apiKey: apiKey,
          apiFormat: 'openai',
          models: models,
        },
      ],
      // 没有 image 能力时留空：画布会判定未就绪并提示，而不是拿文本模型去生图。
      imageModel: image,
      videoModel: video,
      textModel: text,
      audioModel: audio,
      model: pick(image, text),
      models: modelOptions,
    }

    var prev = readJSON(CONFIG_STORE_KEY) || {}
    var prevState = prev.state || {}
    localStorage.setItem(
      CONFIG_STORE_KEY,
      JSON.stringify({
        state: {
          config: config,
          webdav: prevState.webdav || {
            url: '',
            username: '',
            password: '',
            directory: 'infinite-canvas',
            lastSyncedAt: '',
          },
        },
        version: typeof prev.version === 'number' ? prev.version : 0,
      }),
    )
  }

  // ---- 网络 ----
  function api(path, token, opts) {
    var o = opts || {}
    return fetch(API_BASE + path, {
      method: o.method || 'GET',
      headers: Object.assign(
        { Accept: 'application/json', Authorization: 'Bearer ' + token },
        o.body ? { 'Content-Type': 'application/json' } : {},
        o.headers || {},
      ),
      body: o.body ? JSON.stringify(o.body) : undefined,
      credentials: 'same-origin',
    }).then(function (r) {
      if (r.status === 401) {
        var err = new Error('unauthorized')
        err.code = 401
        throw err
      }
      if (!r.ok) throw new Error('HTTP ' + r.status)
      return r.json()
    })
  }

  function listKeys(token) {
    return api('/keys?page=1&page_size=100', token).then(function (d) {
      // 后端是 {data: {items: [...]}} 或 {data: [...]} 两种风格，都兼容
      var payload = d && (d.data || d)
      var items = (payload && payload.items) || (Array.isArray(payload) ? payload : [])
      return items.filter(function (k) {
        return k && k.key && k.status === 'active'
      })
    })
  }

  function listModels(token) {
    return fetch(location.origin + '/v1/models', {
      headers: { Authorization: 'Bearer ' + token },
    })
      .then(function (r) {
        if (!r.ok) throw new Error('HTTP ' + r.status)
        return r.json()
      })
      .then(function (d) {
        var arr = (d && d.data) || []
        return arr
          .map(function (m) {
            return m && m.id
          })
          .filter(Boolean)
          .map(function (id) {
            return { name: id, capability: guessCapability(id) }
          })
      })
  }

  // ---- 加载应用 ----
  // 构建产物里应用脚本被改成了不执行的 <script type="application/json" data-studio-app-src="…">，
  // 这样不依赖「移除 script 元素能取消执行」的浏览器行为——由这里决定什么时候真正跑。
  var appSrc = ''
  function captureAppScript() {
    var el = document.querySelector(APP_SCRIPT_SELECTOR)
    if (!el) return false
    appSrc = el.getAttribute('data-studio-app-src') || el.getAttribute('src') || ''
    el.parentNode.removeChild(el)
    return !!appSrc
  }

  function loadApp() {
    if (!appSrc) return
    var s = document.createElement('script')
    s.type = 'module'
    s.src = appSrc
    s.crossOrigin = 'anonymous'
    document.body.appendChild(s)
  }

  // ---- UI ----
  var STYLE = [
    ':root{--p-500:34 197 94;--p-600:22 163 74;--n-50:250 250 249;--n-100:244 244 243;--n-200:231 231 229;--n-300:212 212 208;--n-500:113 113 122;--n-600:82 82 91;--n-700:63 63 70;--n-900:24 24 27;--r-md:8px;--r-lg:12px}',
    '*{box-sizing:border-box}',
    '#studio-gate{position:fixed;inset:0;z-index:2147483000;display:flex;align-items:center;justify-content:center;padding:24px;',
    'font-family:Inter,"Inter Variable",-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC","Hiragino Sans GB","Microsoft YaHei",sans-serif;',
    'background:rgb(var(--n-50));color:rgb(var(--n-900))}',
    '@media (prefers-color-scheme:dark){#studio-gate{background:#131312;color:#f4f4f3}}',
    '#studio-gate[data-theme=dark]{background:#131312;color:#f4f4f3}',
    '.sg-card{width:100%;max-width:480px;border-radius:var(--r-lg);border:1px solid rgb(var(--n-200));background:#fff;padding:24px}',
    '#studio-gate[data-theme=dark] .sg-card{border-color:#2c2c29;background:#1b1b1a}',
    '@media (prefers-color-scheme:dark){.sg-card{border-color:#2c2c29;background:#1b1b1a}}',
    '.sg-title{margin:0 0 4px;font-size:17px;font-weight:600;letter-spacing:-.01em}',
    '.sg-sub{margin:0 0 18px;font-size:13px;line-height:1.5;color:rgb(var(--n-500))}',
    '.sg-list{display:flex;flex-direction:column;gap:6px;max-height:320px;overflow:auto;margin:0 0 16px;padding:0;list-style:none}',
    '.sg-item{display:flex;align-items:center;gap:10px;width:100%;padding:11px 12px;border:1px solid rgb(var(--n-200));border-radius:var(--r-md);',
    'background:transparent;text-align:left;cursor:pointer;font:inherit;color:inherit;transition:border-color .12s,background-color .12s}',
    '#studio-gate[data-theme=dark] .sg-item{border-color:#2c2c29}',
    '@media (prefers-color-scheme:dark){.sg-item{border-color:#2c2c29}}',
    '.sg-item:hover{border-color:rgb(var(--p-500)/.55);background:rgb(var(--p-500)/.05)}',
    '.sg-item:focus-visible{outline:2px solid rgb(var(--p-500)/.6);outline-offset:1px}',
    '.sg-name{flex:1 1 auto;min-width:0;font-size:14px;font-weight:500;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}',
    '.sg-meta{flex-shrink:0;font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:11px;color:rgb(var(--n-500))}',
    '.sg-empty{font-size:13px;color:rgb(var(--n-500));line-height:1.6}',
    '.sg-empty a{color:rgb(var(--p-600));text-decoration:none;font-weight:500}',
    '.sg-empty a:hover{text-decoration:underline}',
    '.sg-err{margin:12px 0 0;font-size:12px;color:#dc2626}',
    '.sg-spin{display:flex;align-items:center;gap:8px;font-size:13px;color:rgb(var(--n-500))}',
    '.sg-dot{width:6px;height:6px;border-radius:50%;background:rgb(var(--p-500));animation:sg-p 1s ease-in-out infinite}',
    '@keyframes sg-p{0%,100%{opacity:.3}50%{opacity:1}}',
    '@media (prefers-reduced-motion:reduce){.sg-dot{animation:none}}',
  ].join('')

  function maskKey(k) {
    var s = String(k || '')
    if (s.length <= 12) return s
    return s.slice(0, 8) + '…' + s.slice(-4)
  }

  function el(tag, cls, text) {
    var e = document.createElement(tag)
    if (cls) e.className = cls
    if (text != null) e.textContent = text
    return e
  }

  function mountStyle() {
    var s = document.createElement('style')
    s.textContent = STYLE
    document.head.appendChild(s)
  }

  function detectTheme() {
    var raw = readJSON(THEME_STORE_KEY)
    var t = raw && raw.state && raw.state.theme
    if (t === 'light' || t === 'dark') return t
    try {
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    } catch (e) {
      return 'light'
    }
  }

  function showGate(keys, token) {
    mountStyle()
    var gate = el('div', '')
    gate.id = 'studio-gate'
    gate.setAttribute('data-theme', detectTheme())

    var card = el('div', 'sg-card')
    card.appendChild(el('h1', 'sg-title', '选择密钥'))
    card.appendChild(el('p', 'sg-sub', '画布会用它调用模型，产生的费用记在这个密钥上。'))

    if (!keys.length) {
      var empty = el('div', 'sg-empty')
      empty.appendChild(document.createTextNode('还没有可用的密钥。先去 '))
      var a = el('a', '', 'API 密钥')
      a.href = '/keys'
      a.target = '_top'
      empty.appendChild(a)
      empty.appendChild(document.createTextNode(' 创建一个，再回来。'))
      card.appendChild(empty)
    } else {
      var list = el('ul', 'sg-list')
      keys.forEach(function (k) {
        var li = el('li')
        var btn = el('button', 'sg-item')
        btn.type = 'button'
        btn.appendChild(el('span', 'sg-name', k.name || '未命名密钥'))
        btn.appendChild(el('span', 'sg-meta', maskKey(k.key)))
        btn.addEventListener('click', function () {
          chooseKey(token, k.key, gate)
        })
        li.appendChild(btn)
        list.appendChild(li)
      })
      card.appendChild(list)
    }

    gate.appendChild(card)
    document.body.appendChild(gate)
  }

  function chooseKey(token, apiKey, gate) {
    // 换成加载态，避免重复点击
    var card = gate.querySelector('.sg-card')
    card.innerHTML = ''
    card.appendChild(el('h1', 'sg-title', '正在准备画布'))
    var spin = el('div', 'sg-spin')
    spin.appendChild(el('span', 'sg-dot'))
    spin.appendChild(el('span', '', '读取模型列表…'))
    card.appendChild(spin)

    listModels(apiKey)
      .then(function (models) {
        if (!models.length) throw new Error('这个密钥下没有可用模型')
        writeConfig(apiKey, models)
        gate.parentNode.removeChild(gate)
        loadApp()
      })
      .catch(function (err) {
        card.innerHTML = ''
        card.appendChild(el('h1', 'sg-title', '准备失败'))
        card.appendChild(el('p', 'sg-sub', '没法用这个密钥读取模型列表，换一个试试。'))
        card.appendChild(el('p', 'sg-err', err && err.message ? err.message : String(err)))
        var back = el('button', 'sg-item')
        back.type = 'button'
        back.appendChild(el('span', 'sg-name', '返回重新选择'))
        back.addEventListener('click', function () {
          gate.parentNode.removeChild(gate)
          boot()
        })
        card.appendChild(back)
      })
  }

  function showError(message, token) {
    mountStyle()
    var gate = el('div', '')
    gate.id = 'studio-gate'
    gate.setAttribute('data-theme', detectTheme())
    var card = el('div', 'sg-card')
    card.appendChild(el('h1', 'sg-title', '加载失败'))
    card.appendChild(el('p', 'sg-sub', message))
    var retry = el('button', 'sg-item')
    retry.type = 'button'
    retry.appendChild(el('span', 'sg-name', '重试'))
    retry.addEventListener('click', function () {
      gate.parentNode.removeChild(gate)
      boot()
    })
    card.appendChild(retry)
    gate.appendChild(card)
    document.body.appendChild(gate)
  }

  /** 去登录页。嵌在主站 iframe 里时跳顶层窗口，否则登录页会被套进 iframe 里 */
  function toLogin(here) {
    var target = window.top || window
    var back = target === window ? here : target.location.pathname + target.location.search + target.location.hash
    target.location.replace('/login?redirect=' + encodeURIComponent(back))
  }

  // ---- 入口 ----
  function boot() {
    var token = getToken()
    if (!token) {
      var here = location.pathname + location.search + location.hash
      toLogin(here)
      return
    }

    // 已经有配置就直接放行，不再打扰
    if (currentKey()) {
      loadApp()
      return
    }

    var gateShown = false
    listKeys(token)
      .then(function (keys) {
        gateShown = true
        showGate(keys, token)
      })
      .catch(function (err) {
        if (err && err.code === 401) {
          // token 过期：交给登录页续期
          var here = location.pathname + location.search + location.hash
          toLogin(here)
          return
        }
        showError('读取密钥列表失败，请稍后重试。', token)
      })
  }

  if (!captureAppScript()) {
    // 没有应用的 script 标签，说明这页不是画布构建产物，什么都不做
    return
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', boot)
  } else {
    boot()
  }
})()
