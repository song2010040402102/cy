/* 这个计算阶乘的程序可以被cy解释器执行
** 语法采用简化后的c
*/
func main() {
	int i = 0
	int max = 1000000
	int start = time()
	while i++ < max {
		int a = jiecheng(5)
		int b = jiecheng(10)
		int c = jiecheng(20)
		if i == 1 {
			print("a, b, c:", a, b, c)
		}
	}
	print("cost time:", (time()-start)/1000, "us")
}

//计算阶乘
func jiecheng(int a) int {
	if a == 1 {
		return 1
	}
	return a*jiecheng(a-1)
}