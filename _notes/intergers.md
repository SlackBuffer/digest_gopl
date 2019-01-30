# Integers
- 4 sizes
   1. `int8`, `int16`, `int32`, `int64`
   2.`uint8`, `uint16`, `uint32`, `uint64`
- `int` and `uint` are the natural or most efficient size for signed and unsigned integers on a particular platform (32 or 64 bits)
    - `int` is by far the most widely used numeric type
    - Different compilers may make different choices even on identical hardware
- The type ***`rune`*** is a synonym for `int32` and conveniently indicates that a value is a **Unicode code point**. The 2 names may be used interchangeably
    - Rune literals are written as a character within **single quotes**
        - The simplest examples is an ASCII character like 'a', but it's possible to write any Unicode code point either directly or with numeric escapes
    - Printed with `%c`, or with `%q` if quoting is desired, or with `%d` for numeric value
- The type `byte` is a synonym for `uint8`, and emphasizes that the value is a piece of **raw data** rather than s small numeric quantity
- Unsigned integer type `uintptr`
    - Its width is not specified but is sufficient to hold all the bits of a pointer value
    - Only used for low-level programming, such as at the boundary of a Go program with a C library or an operating system
- Regardless of their size, `int`, `uint`, `uintptr` are different types from their explicitly sized sibling
    - Thus `int` is not the same as `int32`, even if the natural size of integers is 32 bits. An explicit conversions is required
- Signed numbers are represented in 2's-complement form, in which the high-order bit is reserved for the sign of the number and the range of values of an n-bit number is from **-2<sup>n-1</sup>** to **2<sup>n-1</sup>-1**
- Unsigned integers use the full range of bits for non-negative values and thus have the range **0** to **2<sup>n</sup>-1**
- The built-in `len` returns a signed `int`

    ```go
    medals := []string{"gold", "silver", "bronze"}
    for i := len(medals) - 1; i >= 0; i-- {
        fmt.Println(medals[i]) // b.., s.., g..
    }
    ```

    - We **tend to** use the signed `int` form even for quantities that can't be negative
    - If `len` returned an unsigned number, then `i` would be a `uint`, and the condition `i >= 0` would always be true by definition
    - After the third iteration, in which `i==0`, the `i--` statement would cause `i` to become not -1, but the maximum `uint` value, and the evaluation of `medals[i]` would fail at run time, or panic, by attempting to access an element outside the bounds of the slice
- For this reason, unsigned numbers **tend to** be used only when their bitwise operators or peculiar arithmetic operators are required, as when implementing bit sets, parsing binary file formats, or for hashing and cryptography. Typically they are not used for merely non-negative quantities
- In general, an explicit conversion is required to convert a value from one type to another, and binary operators for arithmetic and logic (**expect shifts**) must have operands of the same type
    - > **Converting everything to a common type**
- Many integer-to-integer conversions do not entail any change in value; they just tell the compiler how to interpret a value
- But a conversion that narrows a big integer into a smaller one, or a conversion from integer to floating-point or vice versa, may change the value or lose precision
    - Float to integer conversion discards any fractional part, truncating toward zero
    - **Avoid** conversions in which the operand is out of range for the target type, because the behavior depends on the implementation
- Integer literals of any size and type can be written as ordinary decimal numbers, octal numbers, or hexadecimal
- Octal numbers begin with `0`
    - Used for file permissions on POSIX systems
- Hexadecimal numbers begin with `0x` or `0X`
    - Emphasize the bit pattern of a number over its numeric value
> ## Background knowledge
- 机器数：一个数在计算机中的二进制表示；带符号
- 真值：机器数对应的真正的数值
- 原码 = 符号位 + 真值绝对值
    - 原码若用于有减法参与的运算会得到结果错误

        ```
         00010000   (16)                10001000    (-8)
        +10001000   (-8)               +10001000    (-8)
        ---------                     ---------
         10011000   (-24, incorrect)    00010000    (16, incorrect)
        ```

    - 表明正常的加法规则不适用于有负数参与的加法，意味着不得不制定两套运算规则，并为加法运算设计两种电路
- 反码
    - 正数反码即正数本身
    - 负数反码是**其原码**的**符号位以外的各位**逐位**取反**；或该负数的绝对值（即**对应的正数**）的原码的每一个二进制位都取反

    ```
     00010000   (16 的反码)            11110111    (-8 的反码)
    +11110111   (-8 的反码)           +11110111    (-8 的反码)
    ---------                        ---------
     00000111   (7, incorrect)        11111000    (-7, incorrect)  
    ```

- 补码
    - 正数补码即正数本身
    - 负数补码是该负数的绝对值（即**对应的正数**）的原码的**每一个**二进制位都取反，再加 `1`

    ```
     00010000   (16 的补码)            11111000    (-8 的补码)
    +11111000   (-8 的补码)           +11111000    (-8 的补码)
    ---------                        ---------
     00001000   (8 的补码)             11110000    (-16 的补码)  
    ```

    - **[补码本质](http://www.ruanyifeng.com/blog/2009/08/twos_complement.html)**：`负数 = 0 - 该负数的绝对值`，不够减就**借位**

        ```
         00000000   (0)            
        -00001000   (8)           
                    (-8)

        # 发生借位，实际被减数是 100000000；100000000 = 11111111 + 1
        
         100000000
        - 00001000
        ----------
          11111000

        # 上一步等价于一下两步
         11111111   # 取反的步骤
        -00001000
        ---------
         11110111   # 加 1 的步骤
        +00000001
        ---------
         11111000
        ```

- 正数的原码、反码、补码的最高位都是 0，所以最高位为 1 的一定是以某种码表示的负数
- 计算机内部采用负数的与之对应的正数的 2 的补码（Two's Complement）表示负数
- 计算机读到最高位是 1 的数，知道该数是负数，而负数是以补码形式存储的，所以会将该数当补码处理
- 8 bit - 256 个状态
    1. 无符号数
        - 最小数 `00000000`（`0`），最大数 `11111111`（`255`）
    2. 有符号数
        - 右 7 位用于表示数字范围
        - 正数最小数 `00000000`（`+0`），正数最大数 `01111111`（`127`）
        - 负数最“大”数 `10000000`（`-0`），负数最小数 `11111111`（`-127`）；同时存在正零和负零，存在**冗余**
        - `10000000` 是负数，表示形式为补码，减 `1` 得到 `01111111`，再取反，得到与该负数对应的正数的原码 `10000000`，即 128，所以 `10000000` 表示的负数是 -128