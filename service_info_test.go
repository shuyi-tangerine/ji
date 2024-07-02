package ji

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetServiceInfoByName(t *testing.T) {
	serviceInfo, err := GetServiceInfoByName("tangerine/csdn")
	if !assert.Nil(t, err) {
		return
	}

	fmt.Println(serviceInfo)
}

func Test_MarshalServiceInfos(t *testing.T) {
	serviceInfos := []*ServiceInfo{
		{
			Name:       "tangerine/csdn",
			Protocol:   "compact",
			IsBuffered: true,
			IsFramed:   true,
			IP:         "",
			Port:       23460,
			UseSecure:  false,

			WebIP:   "",
			WebPort: 23461,
		},
		{
			Name:       "tangerine/money",
			Protocol:   "json",
			IsBuffered: true,
			IsFramed:   false,
			IP:         "",
			Port:       21080,
			UseSecure:  false,

			WebIP:   "",
			WebPort: 31081,
		},
		{
			Name:       "tangerine/little_book_booker",
			Protocol:   "binary",
			IsBuffered: false,
			IsFramed:   true,
			IP:         "",
			Port:       23458,
			UseSecure:  false,

			WebIP:   "",
			WebPort: 23459,
		},
		{
			Name:       "tangerine/comment",
			Protocol:   "simplejson",
			IsBuffered: false,
			IsFramed:   true,
			IP:         "",
			Port:       23462,
			UseSecure:  false,

			WebIP:   "",
			WebPort: 23463,
		},
	}

	bts, err := json.Marshal(serviceInfos)
	if !assert.Nil(t, err) {
		return
	}

	fmt.Println("==> ", string(bts))
}
