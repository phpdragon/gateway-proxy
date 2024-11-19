package net

import (
	"strings"
	"testing"
)

var urls = []string{
	"https://www.test.com",
	"https://www-sit.test.com",
	"https://www-SIT.test.com",
	"https://www.te-st.com",
	"https://www.te_st.com",
	"https://test.com",
	"https://www.test.me",
	"https://w-w_w.te-s_t.com",
	"https://w-w_w.te-s_t.com:8080/abc",
	"https://w-w_w.te-s_t.c-o_m:8080",
	"https://w-w_w.te-s_t.c-o_m:8080/",
	"https://w-w_w.te-s_t.c-o_m:8080/abc",
}

func TestIsDomainPortUrl(t *testing.T) {
	for _, url := range urls {
		if !IsDomainPortUrl(url) {
			t.Logf("test fail , url: %s", url)
			t.Fail()
		}

		urlTmp := strings.Replace(url, "https://", "http://", 1)
		if !IsDomainPortUrl(urlTmp) {
			t.Logf("test fail , url: %s", urlTmp)
			t.Fail()
		}
	}
}
