package main

func rotateRunes(s string, shift int) string {
	list := []rune(s)
	n := len(list)
	if n == 0 {return ""}
	shift %= n 
	if shift < 0 {shift += n}
	result := make([]rune, n)
	for i := range result{
		result[i] = list[(i+shift) % n]
	}
	return string(result)
}
