package transactions

func cleanUp(obj map[string]interface{}) {
	for k, v := range obj {
		switch t := v.(type) {
		case map[string]interface{}:
			cleanUp(t)
			if len(t) == 0 {
				delete(obj, k)
			}
		case nil:
			delete(obj, k)
		default:
			continue
		}
	}
}
