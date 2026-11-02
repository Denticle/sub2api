# /studio/ —— 无限画布（infinite-canvas）静态站点
#
# 它是 Vite 构建的纯 SPA，没有自己的后端，AI 请求由浏览器直连 /v1/*，
# 所以这里只需要把 dist 当静态目录发出去，再把未命中的路径回落到 index.html。
# 主站的「工作台」用同源 iframe 装它；同源是为了共用 localStorage（登录态、主题、语言）。
#
# 文件放在 /var/www/certbot/studio（nginx 唯一可写挂载的子目录）。
# ACME 挑战只占 /.well-known/acme-challenge/，与 /studio/ 不冲突。
#
# 这个文件不注册 server 块，用 include 挂到 forapi.cloud 的 server 里。
# 后缀是 .inc 不是 .conf：nginx.conf 会 include conf.d/*.conf，放 .conf 会被当成
# 顶层配置再加载一次，location 不能出现在那里。
#
# 缓存策略：
#   assets/  文件名带内容哈希（含 bootstrap），可以当不可变资源长缓存
#   其余     index.html、config.js、插件清单、图标：no-cache，每次带 ETag 回源验证
# 以前 index.html 没有 Cache-Control，浏览器按 Last-Modified 启发式缓存，
# 部署后还会拿旧页面去引用旧 bundle——更新了也不生效。
#
# 每个 location 都要重写一遍安全头：nginx 里 location 只要自己写了 add_header，
# 就不再继承 server 级的 add_header（HSTS 等会悄悄丢掉）。

location = /studio {
  return 301 /studio/;
}

location ^~ /studio/assets/ {
  alias /var/www/certbot/studio/assets/;
  access_log off;
  add_header Cache-Control "public, max-age=31536000, immutable" always;
  add_header Strict-Transport-Security "max-age=31536000" always;
  add_header X-Content-Type-Options "nosniff" always;
  add_header Referrer-Policy "strict-origin-when-cross-origin" always;
  try_files $uri =404;
}

# 运行期配置由容器启动脚本生成，必须不缓存
location = /studio/config.js {
  alias /var/www/certbot/studio/config.js;
  add_header Cache-Control "no-store" always;
  add_header Strict-Transport-Security "max-age=31536000" always;
  add_header X-Content-Type-Options "nosniff" always;
  add_header Referrer-Policy "strict-origin-when-cross-origin" always;
}

location ^~ /studio/ {
  alias /var/www/certbot/studio/;
  index index.html;
  add_header Cache-Control "no-cache" always;
  add_header Strict-Transport-Security "max-age=31536000" always;
  add_header X-Content-Type-Options "nosniff" always;
  add_header Referrer-Policy "strict-origin-when-cross-origin" always;
  # SPA 路由：/studio/canvas/xxx 之类落回 index.html
  try_files $uri $uri/ /studio/index.html;
}
