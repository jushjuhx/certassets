# 横向资产搜集工具

通过证书递归查找域名相关资产的利器，包括根域名和子域名。
## 安装

``` bash
go get github.com/mounsurf/certassets
```

## 使用方式


```go
package main

import (
	"fmt"
	"github.com/mounsurf/certassets"
)

func main() {
	ca := &certassets.CertAssets{
		RootDomainList: []string{"baidu.com"},
		//ShowLog:        true,
		//SubDomainList:  []string{"m", "www"},
	}
	ca.Run()
	rootDomainList := ca.GetRootDomains()
	subDomainList := ca.GetSubDomains()
	
	fmt.Println("baidu.com相关根域名：")
	for _, domain := range  rootDomainList{
		fmt.Println(domain)
	}
	fmt.Println("\nbaidu.com相关子域名：")
	for _, domain := range  subDomainList{
		fmt.Println(domain)
	}
}
```
结果：
```
baidu.com相关根域名：
baidubce.com
hao123.com
nuomi.com
apollo.auto
baidupcs.com
bdtjrcv.com
bcehost.com
dlnel.org
mipcdn.com
baidu.com
bdstatic.com
aipage.cn
aipage.com
bdimg.com
baifae.com
hao222.com
chuanke.com
91.com
baidustatic.com
ssl2.duapps.com
baifubao.com
trustgo.com
smartapps.cn
baiducontent.com
dlnel.com

baidu.com相关子域名：
www.baidu.com.cn
www.baidu.net.ph
baidupcs.com
sni.cloudflaressl.com
baidu.cn
wwww.baidu.com.cn
wn.pos.baidu.com
www.baidu.net.au
dwz.cn
baidu.com.cn
ww.baidu.com
www.baidu.com.hk
click.hm.baidu.com
baifae.com
update.pan.baidu.com
baidu.com
w.baidu.com
baifubao.com
log.hm.baidu.com
www.baidu.hk
wwww.baidu.com
www.baidu.cn
apollo.auto
su.baidu.com
www.baidu.net.tw
www.baidu.net.vn
mct.y.nuomi.com
cm.pos.baidu.com
```

