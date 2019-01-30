# Floating-point
- Go provides 2 sizes floating-point numbers, `float32` and `float64`
    - Their arithmetic properties are governed by IEEE 754 standard implemented by all modern CPUs
- `math.MaxFloat32` (about 3.4e38), `math.MaxFloat32` (about 1.8e308), `math.SmallestNonzeroFloat32` (1.4e-45), `math.SmallestNonzeroFloat64` (4.9e-324)
    - > https://golang.org/pkg/math/#pkg-constants
- A `float32` provides approximately **6** decimal digits of precision
- A `float64` provides about **15** digits
- `float64` should be preferred for most purposes because `float32` computations accumulate error rapidly unless one is quite careful, and the largest positive integer that can be exactly represented as a `float32` is not large

    ```go
    var f float32 = 16777216 // 1<<24
    fmt.Println(f == f+1) // true
    ```

- Digits may be omitted before the decimal point or after it
- Very small or very large numbers are better written in scientific notation, with the letter `e` or `E` preceding the decimal exponent
- Special values
    - The positive and negative infinities represent numbers of excessive magnitude and the result of division by zero
    - `NaN` is the result of such mathematically dubious operations as `0/0` or `Sqrt(-1)`
- `math.IsNaN` tests whether its argument is a not-a-number value; `math.NaN` returns such a value
- Any comparison with `NaN` always yields `false` (except `!=`, which is always the negation of `==`)
    - So don't use `NaN` as a sentinel value in a numeric computation
- If a function that returns floating-point result might fail, it's better to report the failure separately

    ```go
    // Practice
    func compute() (value float64, ok bool) {
        if failed {
            return 0, false
        }
        return result, true
    }
    ```