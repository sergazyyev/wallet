package deflt

func String(origin, def string) string {
	if origin != "" {
		return origin
	}
	return def
}

func StringAny(origin interface{}, def string) string {
	str, ok := origin.(string)
	if ok {
		return str
	}
	return def
}
