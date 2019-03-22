package template

var BootstrapConfigTmpl = `DisablePathCorrection: false
EnablePathEscape: false
FireMethodNotAllowed: true
DisableBodyConsumptionOnUnmarshal: true
DisableVersionChecker: true
TimeFormat: 2006-01-02 15:04:05
Charset: UTF-8
RemoteAddrHeaders:
  X-Real-Ip: true
  X-Forwarded-For: true
  CF-Connecting-IP: true
Other:
  Debug: true
  SiteTitle: "Bootstrap"
  Copyright: "teamlint"
  ServerPort: 8086
  DBConnectionString: "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

  HashKey: "7f26ee15d8c611e8a6cc18ded76e325f"
  BlockKey: "99250d05d8c611e8a6cc18ded76e325f"
`
