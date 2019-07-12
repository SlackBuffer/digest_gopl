
```go
var i, j, k int                 // int, int, int
var b, f, s = true, 2.3, "four" // bool, float64, string
var boiling float64 = 100       // a float64

// tuple assignment
x, y = y, x%y

var f, err = os.Open(name) // returns a file and an error
f, err := os.Open(name)



/*----------------------------------*/
// range
n := 0
for range "Hello, 世界" {
    n++
}



/*----------------------------------*/
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)"
// equivalent to
(*employeeOfTheMonth).Position +=  " (proactive team player)"


func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}

// pp is pointer type
pp := &Point{1, 2}
// equivalent to
pp := new(Point) // returns an address
*pp = Point{1, 2}

p := Point{1, 2}
(&p).ScaleBy(2)
// shorthand
p.ScaleBy(2)
```
