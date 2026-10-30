"""已废弃：主站不再消费画布的 studio:route 消息。

画布内部路由留在 iframe 里自己走，不同步到主站地址栏——同步收益小、
出错面大（history 路由 + 子路径没有服务端回落）。

保留这个文件只是记录当时的决策，不要执行。
"""
raise SystemExit("patch-route-sync 已废弃，不要执行；原因见文件头。")
