
```go
var i, j, k int                 // int, int, int
var b, f, s = true, 2.3, "four" // bool, float64, string
var boiling float64 = 100       // a float64

// tuple assignment
x, y = y, x%y

var f, err = os.Open(name) // returns a file and an error
f, err := os.Open(name)



/*----------------------------------*/
// [ ] iota




/*----------------------------------*/
/* loop */
for initialization; condition; post { /* ... */ }
// "while" loop
for condition { /* ... */ }
// infinite loop
for { /* ... */ }



/*----------------------------------*/
// range
n := 0
for range "Hello, 世界" {
    n++
}



/*----------------------------------*/
/* strings */
conc := "a" + "b" //concatenates
raw := `raw_string`
fstr := Sprintf("I'm %d years old", 28)
// rune literals are written with single quotes.
// these 3 are equivalent: '世' '\u4e16' '\U00004e16'



/*----------------------------------*/
// Operand `coinflip()` is optional.
// A `switch` doesn't have to have an operand (tagless switch).
// Just list the cases, each of which is a boolean expression (equivalent to switch true).
switch coinflip() {
    case "heads":
        head++ // don't fall through by default
    case "tails":
        tails++
        fallthrough
    default: // optional
        fmt.Println("land on edge!")
}



/*----------------------------------*/
/* pointer */
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)" // shorthand
// equivalent to
(*employeeOfTheMonth).Position +=  " (proactive team player)"


// method for *Point
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
p.ScaleBy(2) // shorthand



/*----------------------------------*/
/* arrays and slices */
// array
var a [3]int // array of 3 integers
var q [3]int = [3]int{1, 2, 3}
var r [3]int = [3]int{1, 2}
s := [...]int{1, 2, 3}

type Currency int
const (
    USD Currency = iota
    EUR
    GBP
    RMB
)
symbol := [...]string{RMB: "￥", USD: "$"}

// slice
// like array but size not given
s := []int{0, 1, 2, 3, 4, 5}
// ...

make(()T, len) // capacity equals the length
make([]T, len, cap) // same as make([]T, cap)[:len]

var x []int
x = append(x, 4, 5, 6)
x = append(x, x...) // append the slice x



/*----------------------------------*/
// maps
ages := make(map[string]int)
ages["a"] = 1
ages["b"] = 2

// map literals
ages := map[string]int{
    "a": 1,
    "b": 2,
}
// a new empty map
m := map[string]int{}

var ages map[string]int
ages == nil // true, zero value is nil
len(ages) == 0 // true
// ages["sb"] = 28 // panic: assignment to entry in nil map
ages = make(map[string]int) // allocate first
ages["sb"] = 28 // ok

// += and ++ also work
// safe even if the element is not in the map
ages["a"] = 11
delete(ages, "a")

// random order
for k, v := range ages { /**/ }
// ordered pattern
names := make([]string, 0, len(ages)) // allocate the required size up front
for name := range ages {
    names = append(names, name)
}
// sort the keys explicitly
sort.Strings(names)
for _, name := range names {
    fmt.Printf("%s\t%d\n", name, ages[name])
}

if age, ok := age["x"]; !ok {
    /* "x" is not a key in this map; age = 0 */
}

// map's only legal comparison is with nil
func equal(x, y map[string]int) bool {
    if len(x) != len(y) {
        return false
    }
    for k, xv := range x {
        if yv, ok := y[k]; !ok || yv != xv {
            return false
        }
    }
    return true
}

// use a map whose keys are slices (not comparable or different definition of equality)
var m = make(map[string]int)
// helper that maps each slice to a string
func k(list []string) { return fmt.Sprintf("%q", list) }

func Add(list []string) { m[k(list)]++ }
func Count(list []string) int { return m[k(list)] }



/*----------------------------------*/
/* structs */
// Field order is significant
type Employee struct {
    ID int
    Name, Address string // combined
    DoB time.Time
    Position string
    Salary int
    ManagerID int
}

type Point struct { X, Y int }
// requires a value be specified for every field in the right ordre
p := Point{1, 2}
// flexible
q := Point{Y: 2}

var a struct{} // type struct{}, empty struct
a = struct{}{} // empty struct literal

// shorthand notation (because structs are so commonly dealt with through pointers)
// `&Point{1, 2}` can be used directly within an expression
pp := &Point{1, 2}
// equivalent to
pp := new(Point) // returns an address
*pp = Point{1, 2}



/*----------------------------------*/
/* JSON */
// func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
data, err := json.MarshalIndent(movies, "", " ")

// The matching process that associates JSON names with Go struct names during unmarshaling is case-insensitive,
// so it's only necessary to use a field tag when there's an underscore in the JSON name but not in the Go name
var titles []struct{ Title string }
if err := json.Unmarshal(data, &titles); err != nil {
    log.Fatal("JSON unmarshaling failed: %s", err)
}
fmt.Println(titles)

// [ ] json.Encoder, json.Decoder ch4/github
```
