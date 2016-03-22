package conf

import "testing"

func TestGetConf(t *testing.T) {
	Load("testcase")

	cases := []struct {
		path, value string
	}{
		{"notexistcase", ""},
		{"case1", ""},
		{"case2", ""},
		{"/case2/", ""},

		{"case2/top_key1", "top_value1"},
		{"case2/top_key2/", "top_value2"},
		{"/case2/top_key3", "top_value3"},
		{"/case2/top_key4/", "top_value4"},
		{"case2/top_key5", "top_value5"},

		{"case2/error", ""},
		{"c1group2", ""},

		{"case2/c2group1", ""},
		{"case2/c2group1/notexistkey", ""},
		{"case2/c2group1/g1_key1", "g1_value1"},
		{"case2/c2group1/g1_key2", "g1_value2"},
		{"case2/c2group1/g1_key3", "g1_value3"},
		{"case2/c2group1/g1_key4", "g1_value4"},
		{"case2/c2group1/g1_key5", "g1_value5"},
		{"case2/c2group2", ""},
		{"case2/c2group2/g2_key1", "g2_value1"},
		{"case2/c2group2/g2_key2", "g2_value2"},
		{"case2/c2group2/g2_key3", "g2_value3"},
		{"case2/c2group2/g2_key4", "g2_value4"},
		{"case2/c2group2/g2_key5", "g2_value5"},

		{"case3/c2group1/g1_key1", ""},
		{"case3/c3group1", ""},
		{"case3/c3group1/notexistkey", ""},
		{"case3/c3group1/g1_key1", "g1_value1"},
		{"case3/c3group1/g1_key2", "g1_value2"},
		{"case3/c3group1/g1_key3", "g1_value3"},
		{"case3/c3group1/g1_key4", "g1_value4"},
		{"case3/c3group1/g1_key5", "g1_value5"},
		{"case3/c3group2", ""},
		{"case3/c3group2/g2_key1", "g2_value1"},
		{"case3/c3group2/g2_key2", "g2_value2"},
		{"case3/c3group2/g2_key3", "g2_value3"},
		{"case3/c3group2/g2_key4", "g2_value4"},
		{"case3/c3group2/g2_key5", "g2_value5"},
	}
	for _, c := range cases {
		if v := GetConf(c.path); v != c.value {
			t.Errorf("failed [%v]:[%v] , want %v", c.path, v, c.value)
		}
	}
}