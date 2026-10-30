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
    console.log('[mock] /api/v1/keys ->', KEYS.length, 'keys')
    return json(res, 200, { data: { items: KEYS, total: KEYS.length } })
  }

  // mock：模型列表
  if (p === '/v1/models') {
    console.log('[mock] /v1/models ->', MODELS.data.length, 'models')
    return json(res, 200, MODELS)
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
