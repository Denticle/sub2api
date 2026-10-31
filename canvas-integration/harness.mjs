// 本地验证 /studio/ 的接入层：静态托管 dist + mock 掉两个后端接口。
// 只用于验证 bootstrap 的登录守卫、选 key、写配置、加载应用这条链路。
import http from 'node:http'
import { readFile } from 'node:fs/promises'
import path from 'node:path'

const ROOT = process.env.STUDIO_DIST || '/tmp/studio-dist'
const PORT = Number(process.env.PORT || 4180)

const MIME = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'application/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.json': 'application/json; charset=utf-8',
  '.png': 'image/png',
  '.ico': 'image/x-icon',
  '.woff2': 'font/woff2',
}

const KEYS = [
  { id: 1, key: 'sk-forapi-aaaaaaaaaaaaaaaa1111', name: '生产密钥', status: 'active' },
  { id: 2, key: 'sk-forapi-bbbbbbbbbbbbbbbb2222', name: '测试密钥', status: 'active' },
  { id: 3, key: 'sk-forapi-cccccccccccccccc3333', name: '已停用', status: 'inactive' },
]

const MODELS = { data: [{ id: 'gpt-image-2' }, { id: 'deepseek-v4.1-flash' }, { id: 'grok-imagine-video' }] }

function json(res, code, body) {
  res.writeHead(code, { 'content-type': 'application/json; charset=utf-8' })
  res.end(JSON.stringify(body))
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, `http://127.0.0.1:${PORT}`)
  const p = url.pathname

  // mock：密钥列表（需要 Bearer）
  if (p === '/api/v1/keys') {
    const auth = req.headers.authorization || ''
    if (!auth.startsWith('Bearer ')) return json(res, 401, { code: 401, message: 'unauthorized' })
    // 用于验证 token 过期分支
    if (auth === 'Bearer expired-token') return json(res, 401, { code: 401, message: 'token expired' })
    console.log('[mock] /api/v1/keys ->', KEYS.length, 'keys')
    return json(res, 200, { data: { items: KEYS, total: KEYS.length } })
  }

  // mock：模型列表
  if (p === '/v1/models') {
    console.log('[mock] /v1/models ->', MODELS.data.length, 'models')
    return json(res, 200, MODELS)
  }

  // 模拟主站工作台：侧栏五项 + iframe，消息协议和 frontend/src/views/user/WorkbenchView.vue 一致
  if (p === '/host') {
    res.writeHead(200, { 'content-type': 'text/html; charset=utf-8' })
    return res.end(`<!doctype html><meta charset="utf-8"><style>
      body{margin:0;display:flex;height:100vh;font:13px system-ui;background:#141413;color:#eee}
      nav{width:150px;padding:12px;display:flex;flex-direction:column;gap:4px}
      nav button{all:unset;padding:8px 10px;border-radius:8px;cursor:pointer}
      nav button.on{background:#22c55e22;color:#4ade80}
      iframe{flex:1;border:0}</style>
      <nav id="nav"></nav><iframe id="f" src="/studio/canvas"></iframe>
      <script>
        const TOOLS=['canvas','image','video','prompts','assets'];let cur='canvas',ready=false,childTool='';
        window.__routes=[];
        const nav=document.getElementById('nav'),f=document.getElementById('f');
        const draw=()=>{nav.innerHTML='';TOOLS.forEach(t=>{const b=document.createElement('button');b.textContent=t;b.dataset.tool=t;if(t===cur)b.className='on';b.onclick=()=>go(t);nav.appendChild(b)})};
        const send=t=>f.contentWindow.postMessage({type:'studio:navigate',path:'/'+t},location.origin);
        function go(t){const same=t===cur;cur=t;draw();if(ready&&(same||t!==childTool))send(t)}
        addEventListener('message',e=>{if(e.origin!==location.origin||e.source!==f.contentWindow)return;const d=e.data;if(!d||d.type!=='studio:route')return;
          window.__routes.push(d.path);const first=d.path.replace(/^\\/+/,'').split('/')[0];const t=TOOLS.includes(first)?first:'';
          if(!ready){ready=true;childTool=t;if(t!==cur)send(cur);return}
          childTool=t;if(t&&t!==cur){cur=t;draw()}});
        draw();
      </script>`)
  }

  // mock：登录页（守卫跳回来时要有东西接）
  if (p === '/login') {
    res.writeHead(200, { 'content-type': 'text/html; charset=utf-8' })
    return res.end('<h1 id="login-stub">LOGIN PAGE</h1>')
  }

  let file
  if (p === '/' || p === '/studio' || p === '/studio/') file = 'index.html'
  else if (p.startsWith('/studio/')) file = p.slice('/studio/'.length)
  else file = p.slice(1)

  const full = path.join(ROOT, file)
  try {
    const buf = await readFile(full)
    res.writeHead(200, { 'content-type': MIME[path.extname(full)] || 'application/octet-stream' })
    res.end(buf)
  } catch {
    // SPA 回落
    try {
      const buf = await readFile(path.join(ROOT, 'index.html'))
      res.writeHead(200, { 'content-type': 'text/html; charset=utf-8' })
      res.end(buf)
    } catch {
      res.writeHead(404)
      res.end('not found')
    }
  }
})

server.listen(PORT, '127.0.0.1', () => console.log(`studio harness on http://127.0.0.1:${PORT}/studio/`))
