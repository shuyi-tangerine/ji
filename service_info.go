package ji

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

var (
	_defaultServiceInfos map[string]*ServiceInfo
)

func init() {
	serviceInfos, err := GetServiceInfosByEnvKey(context.Background(), "SERVICE_INFOS_PATH")
	if err != nil {
		panic("github.com/shuyi-tangerine/ji need env SERVICE_INFOS_PATH")
	}
	mapSI := make(map[string]*ServiceInfo, len(serviceInfos))
	for _, serviceInfo := range serviceInfos {
		mapSI[serviceInfo.Name] = serviceInfo
	}
	_defaultServiceInfos = mapSI
}

func GetServiceInfoByName(name string) (serviceInfo *ServiceInfo, err error) {
	serviceInfo = _defaultServiceInfos[name]
	if serviceInfo == nil {
		return nil, fmt.Errorf("no service name:%s", name)
	}
	return serviceInfo, nil
}

func GetServiceInfosByEnvKey(ctx context.Context, envKey string) (serviceInfos []*ServiceInfo, err error) {
	path := os.Getenv(envKey)
	if path == "" {
		path = "service_infos.json"
	}

	return GetServiceInfosByPath(ctx, path)
}

func GetServiceInfosByPath(ctx context.Context, path string) (serviceInfos []*ServiceInfo, err error) {
	bts, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &serviceInfos)
	if err != nil {
		return
	}

	return
}

type ServiceInfo struct {
	Name string `json:"name"`

	Protocol   string `json:"protocol"`
	IsBuffered bool   `json:"is_buffered"`
	IsFramed   bool   `json:"is_framed"`
	IP         string `json:"ip"`
	Port       int64  `json:"port"`
	UseSecure  bool   `json:"use_secure"`

	WebIP   string `json:"web_ip"`
	WebPort int64  `json:"web_port"`
}

func (m *ServiceInfo) GetAddr() string {
	return fmt.Sprintf("%s:%d", m.IP, m.Port)
}

func (m *ServiceInfo) GetWebAddr() string {
	return fmt.Sprintf("%s:%d", m.WebIP, m.WebPort)
}
