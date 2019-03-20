package certassets

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type CertAssets struct {
	RootDomainList []string            //待测试的域名[必选]。数组格式，如果已知一个公司包含多个顶级域名，可以直接指定这些域名。如["baidu.com"],["taobao.com","alibaba.com"]
	ShowLog        bool                //是否显示日志
	SubDomainList  []string            //用于测试添加的子域名字典（默认为["www"]）[可选]
	TopDomainList  []string            //顶级域名列表，用于初始化TopDomainMap。(默认为["com", "net", "org", "gov", "edu", "co"])[可选]
	queueDomains   []string            //待查询队列
	topDomainMap   map[string]struct{} //顶级域名map，用于判断域名对应的根域名时使用。如"www.test.com.cn"的根域名是"test.com.cn"而不是"com.cn"
	domainMap      map[string]struct{} //queueDomains 对应的map，用于去重，避免重复添加到队列中
	rootDomainMap  map[string]struct{} //证书中*.xxx.com 类型的DNS Names
	subDomainMap   map[string]struct{} //证书中xxx.xxx.com 类型的DNS Names
}

func (ca *CertAssets) init() {
	if len(ca.SubDomainList) == 0 {
		ca.SubDomainList = []string{"www"}
	}
	ca.queueDomains = ca.RootDomainList
	ca.domainMap = map[string]struct{}{}
	for _, rootDomain := range ca.RootDomainList {
		ca.domainMap[rootDomain] = struct{}{}
		for _, subDomain := range ca.SubDomainList {
			if !strings.HasPrefix(rootDomain, subDomain+".") {
				ca.queueDomains = append(ca.queueDomains, subDomain+"."+rootDomain)
				ca.domainMap[subDomain+"."+rootDomain] = struct{}{}
			}
		}
	}
	ca.rootDomainMap = map[string]struct{}{}
	ca.subDomainMap = map[string]struct{}{}
	if len(ca.TopDomainList) == 0 {
		ca.TopDomainList = []string{"com", "net", "org", "gov", "edu", "co"}
	}
	ca.topDomainMap = map[string]struct{}{}
	for _, domain := range ca.TopDomainList {
		ca.topDomainMap[domain] = struct{}{}
	}
}

func (ca *CertAssets) getDnsNames(domain string) ([]string, error) {
	cfg := tls.Config{}
	dialer := &net.Dialer{
		Timeout: time.Second * 3,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", domain+":443", &cfg)
	if err != nil {
		return nil, err
	}
	if len(conn.ConnectionState().PeerCertificates) < 1 {
		return nil, errors.New("unknown cert error")
	}
	return conn.ConnectionState().PeerCertificates[0].DNSNames, nil
}

func (ca *CertAssets) getRootDomain(domain string) string {
	domainSlice := strings.Split(domain, ".")
	length := len(domainSlice)
	if length <= 2 {
		return domain
	}
	if _, ok := ca.topDomainMap[domainSlice[length-2]]; ok {
		return domainSlice[length-3] + "." + domainSlice[length-2] + "." + domainSlice[length-1]
	} else {
		return domainSlice[length-2] + "." + domainSlice[length-1]
	}

}

func (ca *CertAssets) Run() {
	ca.init()
	index := -1
	for {
		index++
		if index >= len(ca.queueDomains) {
			break
		}
		if ca.ShowLog {
			fmt.Printf("total: %d\tleft: %d\tcurrent: %s\n", len(ca.queueDomains), len(ca.queueDomains)-index-1, ca.queueDomains[index])
		}
		newDnsNames, err := ca.getDnsNames(ca.queueDomains[index])
		if err != nil {
			continue
		}
		for _, newDnsName := range newDnsNames {
			if len(newDnsName) < 4 {
				continue
			}
			newDomainList := []string{}
			if newDnsName[0] == '*' {
				ca.rootDomainMap[newDnsName[2:]] = struct{}{}
				newDomainList = append(newDomainList, newDnsName[2:])
				for _, subDomain := range ca.SubDomainList {
					if !strings.HasPrefix(newDnsName, subDomain+".") {
						newDomainList = append(newDomainList, subDomain+newDnsName[1:])
					}
				}
			} else {
				ca.subDomainMap[newDnsName] = struct{}{}
				newDomainList = append(newDomainList, newDnsName)
			}
			for _, newDomain := range newDomainList {
				if _, ok := ca.domainMap[newDomain]; !ok {
					ca.queueDomains = append(ca.queueDomains, newDomain)
					ca.domainMap[newDomain] = struct{}{}
				}
			}
		}
	}
}

func (ca *CertAssets) GetRootDomains() []string {
	rootDomainMap := map[string]struct{}{}
	for rootDomain, _ := range ca.rootDomainMap {
		if rootDomain == ca.getRootDomain(rootDomain) {
			rootDomainMap[rootDomain] = struct{}{}
		}
	}
	for rootDomain, _ := range ca.rootDomainMap {
		realRootDomain := ca.getRootDomain(rootDomain)
		if _, ok := rootDomainMap[realRootDomain]; !ok {
			rootDomainMap[rootDomain] = struct{}{}
		}
	}
	result := []string{}
	for domain, _ := range rootDomainMap {
		result = append(result, domain)
	}
	return result
}

func (ca *CertAssets) GetSubDomains() []string {
	result := []string{}
	for domain, _ := range ca.subDomainMap {
		result = append(result, domain)
	}
	return result
}
