package quick

import (
	"net/http"
	"strings"
)

// defined []http.Cookie alias Cookies
type Cookies = []*http.Cookie

// You should init it by using NewCookiesWithString like this:
// 	cookies := quick.NewCookiesWithString(
//		"key1=value1; key2=value2; key3=value3"
// 	)
// Note: param is cookie string
func NewCookiesWithString(rawstr string) Cookies {
	if len(rawstr) == 0 {
		return nil
	}
	strs := strings.Split(rawstr, ";")
	cookies := make(Cookies, 0, len(strs))
	for i := 0; i < len(strs); i++ {
		cookie := strings.Split(strs[i], "=")
		if len(cookie) != 2 {
			continue
		}
		cookies = append(cookies, &http.Cookie{
			Name:  strings.TrimSpace(cookie[0]),
			Value: strings.TrimSpace(cookie[1]),
		})
	}
	return cookies
}
