func main() {
	int a = 9
	int b = 9
	int i = 0
	while ++i <= a {
		int j = 0
		while ++j <= b {
			print(i, "*", j, "=", i*j)
		}
	}
}