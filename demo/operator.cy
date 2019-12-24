func main() {
	int a b
	a += 1<<3
	b -= 8>>1
	print(a|-b, a&-b, a^-b, ~a)
	print(a, b, a%3, a*(b+6)/2)
	print(true&&(false||!true))
	print(a/2>-b?a+b:a-b)
}