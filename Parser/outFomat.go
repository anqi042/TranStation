package Parser
/*
{"host":"6CU401A1VN","reg_test.test":"whatareudoing?","test.TAG":"1","topic":"test"}

topic:host,key,value...
*/

func FormatOutput(m map[string]string) string  {
	const TOPICKEY="topic"
	const HOSTKEY  ="host"
	const TAGPOSTFIX = ".TAG"
	tagKey := m[TOPICKEY] + TAGPOSTFIX
	kv := ""
	for k,v := range m{
		if k != HOSTKEY && k != tagKey && k != TOPICKEY{
			//normal key value
			kv = kv + k + "," + v +","
		}
	}
	if elem,ok := m[tagKey];ok{
		kv += QueryTag(m[TOPICKEY])+","+elem+","
	}

	raw := m[TOPICKEY] + ":" + m[HOSTKEY] + "," + kv
	return raw[:len(raw)-1]
}
