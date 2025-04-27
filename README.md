
### 功能:
1. 切换代理
2. 管理员启动

配置说明:
proxy_type: 启动时的代理方式, 1->系统代理  2->tun 0->关闭
github_proxy: 下载资源(mihomo/GeoSite/GeoIp)时使用的代理
http_proxy: 订阅代理
core: 启动的核心时参数,主要是ip
gui_port: 没用
profiles: 配置订阅, 每个filters 会加一个代理组 如果不是 "🐟 漏网之鱼", "🚀 节点选择", "DIRECT",  会自动添加代理组
prepend_rules: 追加rule, 
cover_config: 覆盖配置,其中的rule 优先于 prepend_rules
