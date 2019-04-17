# Slices
- Slices represent variable-length sequences whose elements all have the same type
- A slice type is written `[]T`, where the elements have type `T`; looks like an array without a size
- A slice is a lightweight data structure that gives access to a subsequence (or perhaps all) of the elements of an array (known as the slice's underlying array)
- A slice has 3 components: a pointer, a length, and a capacity
   1. The **pointer** points to the first elements of the array that's reachable through the slice, which is not necessarily the array's first element
   2. The length is the number of slice elements; can't exceed the capacity
        - `len`
   3. Usually is the number of elements between the start of the slice and the **end** of the underlying array
        - `cap`
- Multiple slices can share the same underlying array and may refer to overlapping parts of that array
- `s[i:j]` (`0<=i<=j<=cap(s)`)
    - `s` may be an array variable, a pointer to an array, or another slice
    - If `i` is omitted, it's 0; if `j` is omitted, it's `len(s)`
- Slicing beyond `cap(s)` causes a panic
- Slicing beyond `len(s)` extends the slice, so the result may be longer than the original
- Since a slice contains a pointer to an element of an array, passing a slice to a function permits the function to **modify** the underlying array elements. In other words, copying a slice creates an alias for the underlying array
- A slice literal looks like an array literal, a sequence of values separated by commas and surrounded by braces, but the size is not given. 
- This implicitly creates an array variable of the right size and yields a slice that points to it

    ```go
	s := []int{0, 1, 2, 3, 4, 5}
    ```

    - As with array literals, slice literals may specify the values in order, or give their indices explicitly, or use a mix of the 2 styles
- Unlike arrays, slices are not comparable
    - So we cannot use `==` to test whether 2 slices contains the same elements
    - The standard library provides `bytes.Equal` for comparing 2 slices of `[]byte`
    - The only legal slice comparison is against `nil`
- The elements of a slice are indirect, making it possible for slice to contain itself
- A fixed slice value may contain different elements at different times as the contents of the underlying array are modified
- A hash table such as Go's map type makes only shallow copies of its keys, it requires that equality for each key remain the same throughout the lifetime of the hash table


- The zero value of a slice type is `nil`
- A nil slice has no underlying array. The nil slice has length and capacity zero
- There're also non-nil slices of length and capacity zero (**`[]int{}`**, `make([]int, 3)[3:]`)
- The nil value of a particular slice can be written using a conversion expression such as `[]int(nil)`
- Use `len(s) == 0`, not `s == nil` to test whether a slice is empty
- Other than comparing equal to `nil`, a nil slice behaves like any other zero-length slice
    - Unless clearly documented to the contrary, Go functions should treat all zero-length slices the same way, whether nil or non-nil
- `make` creates an unnamed array variable and returns a slice of it; the array is accessible only through the returned slice

    ```go
    make(()T, len) // capacity equals the length
    make([]T, len, cap) // same as make([]T, cap)[:len]
    ```

## `append`
- `append` appends items to slices

    ```go
    var runes []rune
    for _, r := range "hello, 世界" {
        runes = append(runes, r)
    }
    fmt.Printf("%q\n", runes)

    fmt.Printf("%q\n", []rune("hello, 世界"))

    var x []int
    x = append(x, 4, 5, 6)
    x = append(x, x...) // append the slice x
    ```

- Usually we don't know whether a given call to `append` will cause a reallocation, so we can't assume that the original slice refers to the same array as the resulting slice, nor that it refers to a different one
- Similarly, we must not assume that operations on elements of the old slice will (or will not) be reflected in the new slice
- As a result, it's usual to assign the result of a call to `append` to the same slice whose value was passed to `append`
- **Updating the slice variable** is required not just when calling `append`, but for any function that may change the **length** or **capacity** of a slice or make it refer to a different **underlying array**
- Although the elements of the underlying array are indirect, the slice's pointer, length, and capacity are not. To update them requires an assignment like `runes = append(runes, r)`
- In this respect, slices are not "pure" reference types but ***resemble an aggregate type*** such as this struct

    ```go
    type IntSlice struct {
        ptr *int
        len, cap int
    }
    ```

## In-place slice techniques
- Slices sharing the same underlying array

    ```go
    func nonempty2(strings []string) {
        out := strings[:0] // zero-length slice of original
    }
    ```