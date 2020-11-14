/* 这个计算阶乘的程序可以被cy解释器执行
** 语法采用简化后的c
*/
func main() {
	print(jiecheng(5))
	print(jiecheng(10))
	print(jiecheng(20))
}

//计算阶乘
func jiecheng(int a) int {
	if a == 1 {
		return 1
	}
	return a*jiecheng(a-1)
}