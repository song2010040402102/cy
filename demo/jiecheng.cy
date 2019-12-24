func main() {
	print(jiecheng(5))
	print(jiecheng(10))
	print(jiecheng(20))
}

func jiecheng(int a) int {
	if a == 1 {
		return 1
	}
	return a*jiecheng(a-1)
}