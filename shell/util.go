package shell

const endDelim = "echox-end\n"

func endswith(b, sub []byte) bool {
	lb := len(b)
	lsub := len(sub)
	if lsub > lb {
		return false
	}
	for i := 1; i < lsub; i++ {
		if b[lb-i] != sub[lsub-i] {
			return false
		}
	}
	return true
}
