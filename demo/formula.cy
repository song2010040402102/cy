func main() {
    int i = 0
    int max = 1000000
    int start = time()
    float a b c x
    a = 100
    b = 300
    c = 125
    while i++ < max {
        x = formula(a, b, c)
        if i == 1 {
            print("a, b, c, x:", a, b, c, x)
        }
    }
    print("cost time:", (time()-start)/1000, "us")
}

func formula(float a, float b, float c) float {
    return (pow(pow(b, 2.0) - 4*a*c > 0? pow(b, 2.0) - 4*a*c : 0, 0.5) - b)/(a*2.0)
}