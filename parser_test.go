package uci

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	cfg, err := Parse("test", `
# For full documentation of mwan3 configuration:
# https://openwrt.org/docs/guide-user/network/wan/multiwan/mwan3#mwan3_configuration

Config globals 'globals'
	option mmx_mask '0x3F00'

Config interface 'wan'
	option enabled '1'
	list track_ip '8.8.4.4'
	list track_ip '8.8.8.8'
	list track_ip '208.67.222.222'
	list track_ip '208.67.220.220'
	option family 'ipv4'
	option reliability '2'

Config interface 'wan6'
	option enabled '0'
	list track_ip '2001:4860:4860::8844'
	list track_ip '2001:4860:4860::8888'
	list track_ip '2620:0:ccd::2'
	list track_ip '2620:0:ccc::2'
	option family 'ipv6'
	option reliability '2'

Config interface 'wanb'
	option enabled '0'
	list track_ip '8.8.4.4'
	list track_ip '8.8.8.8'
	list track_ip '208.67.222.222'
	list track_ip '208.67.220.220'
	option family 'ipv4'
	option reliability '1'

Config interface 'wanb6'
	option enabled '0'
	list track_ip '2001:4860:4860::8844'
	list track_ip '2001:4860:4860::8888'
	list track_ip '2620:0:ccd::2'
	list track_ip '2620:0:ccc::2'
	option family 'ipv6'
	option reliability '1'

Config member 'wan_m1_w3'
	option interface 'wan'
	option metric '1'
	option weight '3'

Config member 'wan_m2_w3'
	option interface 'wan'
	option metric '2'
	option weight '3'

Config member 'wanb_m1_w2'
	option interface 'wanb'
	option metric '1'
	option weight '2'

Config member 'wanb_m2_w2'
	option interface 'wanb'
	option metric '2'
	option weight '2'

Config member 'wan6_m1_w3'
	option interface 'wan6'
	option metric '1'
	option weight '3'

Config member 'wan6_m2_w3'
	option interface 'wan6'
	option metric '2'
	option weight '3'

Config member 'wanb6_m1_w2'
	option interface 'wanb6'
	option metric '1'
	option weight '2'

Config member 'wanb6_m2_w2'
	option interface 'wanb6'
	option metric '2'
	option weight '2'

Config policy 'wan_only'
	list use_member 'wan_m1_w3'
	list use_member 'wan6_m1_w3'

Config policy 'wanb_only'
	list use_member 'wanb_m1_w2'
	list use_member 'wanb6_m1_w2'

Config policy 'balanced'
	list use_member 'wan_m1_w3'
	list use_member 'wanb_m1_w2'
	list use_member 'wan6_m1_w3'
	list use_member 'wanb6_m1_w2'

Config policy 'wan_wanb'
	list use_member 'wan_m1_w3'
	list use_member 'wanb_m2_w2'
	list use_member 'wan6_m1_w3'
	list use_member 'wanb6_m2_w2'

Config policy 'wanb_wan'
	list use_member 'wan_m2_w3'
	list use_member 'wanb_m1_w2'
	list use_member 'wan6_m2_w3'
	list use_member 'wanb6_m1_w2'

Config rule 'https'
	option sticky '1'
	option dest_port '443'
	option proto 'tcp'
	option use_policy 'balanced'

Config rule 'default_rule_v4'
	option dest_ip '0.0.0.0/0'
	option use_policy 'balanced'
	option family 'ipv4'

Config rule 'default_rule_v6'
	option dest_ip '::/0'
	option use_policy 'balanced'
	option family 'ipv6'
`)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Get("globals").Get("mmx_mask").Values[0] != "0x3F00" {
		t.Fatal()
	}
}

func testParser(t *testing.T, name, input string, expected []token) {
	t.Helper()

	var i int
	ok := scan(name, input).each(func(tok token) bool {
		if dump["token"] {
			fmt.Println(tok) //nolint:forbidigo
		}

		if i >= len(expected) {
			t.Errorf("token %d, unexpected item: %s", i, tok)
			return false
		}
		if ex := expected[i]; tok.typ != ex.typ || !equalItemList(tok.items, ex.items) {
			t.Errorf("token %d\nexpected %s\ngot      %s", i, ex, tok)
			return false
		}

		i++
		return true
	})
	if !ok {
		return
	}

	if l := len(expected); i != l {
		t.Errorf("expected to scan %d token, actually scanned %d", l, i)
	}
}

func equalItemList(a, b []item) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].typ != b[i].typ || a[i].val != b[i].val {
			return false
		}
	}
	return true
}
