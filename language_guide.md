# A Tour of Rose

## Introduction
Rose is a general-purpose, statically typed language that focuses on writeability, readability, safety, and simplicity (in that order).  It is designed to provide the expressiveness and elegance of dynamically typed languages, like Python, but keep the safety and maintainability of statically typed languages, such as Go. Rose has explicit support for concurrency, and aims to make error handling easy to write and understand. Rose is also garbage collected, so you don't have to worry about managing memory.

## Basic Syntax
### Hello World
I guess I should show you the obligatory Hello World program... here it is:
```
print("Hello, World!")
```
That's it. No declaring main, no specifying what package you're in, nothing. Outside of functions and packages are automatically run in "main". Rose also does not need terminating semicolons.

### Comments
Comments are C-style in Rose. 
1.  _Line comments_ start with the character sequence `//` and stop at the end of the line.
2.  _General comments_ start with the character sequence `/*` and stop with the first subsequent character sequence `*/`.

A comment cannot start inside a char or string literal, or inside a comment. A general comment containing no newlines acts like a space. Any other comment acts like a newline.

### Packages
If you want write a reusable library, you have to declare your package:
```
package myMath

fn Add(x int, y int) int {
	return x + y
}
```
Any name that begins with a capital letter is *exported*, or usable outside of the package it's in. 
To use a package, you have to import it:
```
import "myMath"
// idiomatic way to import multiple packages
import (
	"moarMath"
	"newMath"
)
```
You can also create alias names for imports for convenience:
```
import myPkg "this/is/a/very/long/package/name"
```

### Variables
Variable types can be inferred in most cases:
```
hello = "hello there!"
grade = 88
```
Variables can be declared with the var keyword:
```
var i int
var s string
var (
    x int
    n, f float
)
var python, java, c bool
```
If you want to be verbose, you can explicitly specify variable types when instantiating a variable as well:
```
var i int = 5
var s string = "im feeling stringy"
```
Notice that type names appear *after* the variable name. This is to maximize readability, especially when working with functions as we'll see later.

## Logic Control
### If Statements
If statements have standard syntax in Rose. Curly braces `{}` are always required though. Else statements *must* start directly after the closing curly brace of the preceding if statement. Rose uses the keywords `and, or, not` as boolean operators, these can be used to create complex but readable boolean expressions.
```
imblue = true
if imblue {
	print("I'm blue")
} else {
	print("...")
}

myColor = "blue"
if imblue and myColor == "blue" {
	print("Da ba dee da ba da")
}
```
In addition, Rose supports simple statements before the condition. This can be useful to reuse the returned value from a function multiple times:
```
import "math"

if s = math.Sqrt(16); if s > 0 and s == 4 {
	print("math doesn't lie")
}
```

### Looping
Rose only has one kind of loop: the for loop. This may seem limiting, however Rose's for loop can handle all common looping types. A traditional for loop looks about how you'd expect it to. A for loop with one condition is essentially a while loop, and a for loop with no condition is an infinite loop.
```
for i = 0; i < 5; i++ { // standard for loop
	print(i)
}

stillGoing = true
for stillGoing { // essentially a while loop
	stillGoing = false
}

for { // infinite loop
	print("This will print FOREVER")
}
```

## Types
### Basic Types
As you saw above, all the common types you'd expect are present in Rose: `bool, int, float, string`. 
These basic types are immutable, meaning their values can never change. That isn't to say that variables cannot be reassigned. Let me show you an example:
```
truther = true
numbah = 5
percent = 96.7
str = "Hello, World!"
str[1] = '3' // illegal: strings are immutable
str = "Yes, this works" // fine: variable "str" can be reassigned
```
Note that numeric types (`int` and `float`) can be of arbitrary size. The only limit to how big these numbers can be is the machine's memory you're running a Rose program on.

Rose requires explicit conversions when working with items of different types. This is to minimize confusion and make intent clear.
```
// Some basic numeric conversions:
i = 42
f = float(i)

// Basic string conversions:
x = 78
xStr = string(x)

// Working with floats and ints
sum = 48.9 + float(x) // sum is a float of 126.9
```
Rose also has some less commonly used, but still useful types: `char, byte, bytes`
The `char` and `byte` types are immutable, while the `byte` type is mutable. A `string` holds `char`s, and a `bytes` holds `byte`s.
```
b = bytes("Good old-fashioned oatmeal")
b[9] = '_' // fine, bytes are mutable

s = "New modern  oatmeal"
s[11] = 'g' // illegal, string is immutable
```
As you undoubtedly noticed, char literals use single quotes instead of double quotes. 
String literals *must* use double quotes, unless you want a raw string. Raw string literals are created by enclosing a string in back quotes, and escapes don't apply and the string may contain newlines.
```
raw = `Hello!
This is a long, pointless message.

Bye!
`
```
#### The 'any' Type
The `any` type in Rose is special in that it can represent anything. Shocking I know. Declared variables of the `any` type can be assigned anything at any point.
```
var foo any
foo = 5
foo = "five"
foo = 'f'
```

### Zero Values
Every type in Rose is has a zero value, and declared variables are zero-initialized. In other words, there is no such thing as an uninitialized variable in Rose. If you don't initialize your variables, Rose will for you. Zero values for the previously discussed basic types are as follows:
```
bool: false
int: 0
float: 0.0
char: unicode null character (U+0000)
string: "" (empty string)
any: nil
```
`any` is unique in that it's zero type is `nil`, which is a constant that represents the absence of a value.
The zero value of a type is falsey, and any non-zero value of a type is truthy.

### Container Types
What good is a programming language without support for lists and maps? Fret not, Rose has those and more. Rose natively supports these container types: `list, map, set, tuple`. All container types in Rose are mutable, except for tuples, which are always immutable.

The zero value of all containers is a container with a length and capacity of 0, or an empty container.
All mutable container types (lists, maps, and sets) have both a length and capacity. The length is the number of elements in the container, and the capacity is the number of elements allocated in memory. The following relationship always holds:
```
0 <= len(s) <= cap(s)
``` 
As shown above, the `len` and `cap` builtin functions are used to check the length and capacity of containers.
Note that because tuples are immutable, their capacity will always equal their length.

#### Lists
Lists are dynamically sized arrays, their usage should be pretty familiar:
```
myList = [1, 2, 3] // list of int
var sameList list[int] = myList // explicitly specifying the type of list
anyList = ["blue", 5, 'e', false] // list of any

myList.append(7)
print(myList[6]) // 7
```
If a list literal contains multiple types, the list will be a `list` of `any`.

Lists also have some operators for making working with them easy:
```
list4 = [1, 2, 3]
list4 += [4, 5, 6] // the '+' and '+=' operators combine lists
print(list4) // [1, 2, 3, 4, 5, 6]
list5 = ["this", "is", "a", "list"]
list5 << "son" // the "<<" operator can append to lists
"hey" >> list5 // the ">>" operator can prepend to lists
```

#### Maps
Maps are unordered key-value pairs, also referred to in other languages as 'dictionaries' or 'hashes'. Maps have a constant time lookup time for a value, provided that a key that exists in the map is given. This also means that maps are guaranteed to contain unique keys. Only immutable types can be keys though.
```
myMap = {"key": "value", "peanut": "butter", 1: true} // an any any map
var dankMap map[int]float = {1: 6.7, 6: 8.9, 3: 99.45} // explicit specification of map type
print(myMap["peanut"]) // butter
maMap["newKey"] = 'f' // map assignment
```
Trying to access a non-existent key results in the zero value of the map's value type getting returned. For example:
```
var emptyMap map[string]int
print(emptyMap["key"]) // 0
```
There are ways to check if a key exists in a map, but we'll touch on that later.

#### Sets
Sets are very similar to lists, except they are unordered, and contain unique values. Like maps they can only contain immutable types. Sets are useful when you want to preform math operations on sets:
```
mySet = {1, 2, 3} // a set of int
var yourSet set[int] = {2, 4, 3} // explicit specification of set type
print(mySet | yourSet) // set union; prints {1, 2, 3, 4}
print(mySet - yourSet) // set difference; prints {1}
```
Sets can be used to quickly remove duplicates from lists as well:
```
nums = [1, 4, 6, 3, 8, 4, 9, 6, 6]
uniqueNums = set(nums)
print(uniqueNums) // {1, 4, 6, 3, 8, 9}
```

#### Tuples
Tuples are immutable collections of values. Tuples can contain any combination of values, whether they are immutable or not. Tuple literals are values separated by commas (and optionally parentheses):
```
t = 123, "hello!", 'f'
var t2 tuple[int, string, char] = t // explicit specification of tuple type
t3 = (t, t2)
```
When creating empty tuples or tuples with one element, the syntax is a bit different. Empty parentheses `()` specifies an empty tuple, and a tuple with one element can be created by appending a comma to the value:
```
emptyTuple = ()
babyTuple = "one element",
```
Like lists and sets, tuple elements can be accessed by indexing.

#### Iterating Over Containers
Rose has for in loops that make it easy to iterate over containers. An inline if statement is optional and can be used to filter elements.
```
nums = [1, 2, 3]
for n in nums {
	print(n + 1)
}

grades = [50, 96, 7, 88, 73, 10]
for grade in grades if grade > 70 {
	print("Passing grade: " + string(grade))
}
```

#### Container Unpacking
Containers can easily be unpacked by assigning multiple variables to a container. Multiple assignment is really just unpacking a tuple inline. This can even be used to great effect in for in loops:
```
t = ['h', 'i'], 8
x, y = t

people = [
	("Bob", 42, "Mechanic"),
	("James", 24, "Artist"),
	("Harry", 32, "Lecturer")
]

for name, age, profession in people {
	printf("Name: %s, Age: %d, Profession: %s", name, age, profession)
}
```
This requires that there are the same number of variables getting assigned as there are elements in the container.
An underscore `_` in place of a variable can be used to ignore elements from the container:
```
l = [1, 2, 3]
a, _, b = l // a = 1, b = 3
people = ["james", "charity", "john"]
```
In some situations, you may want to isolate a few elements, then keep the rest of the elements together. An asterisk `*` can be used to tell Rose to collect the unassigned elements from the unpacking, or to unpack containers in-place:
```
nums = [1, 2, 3, 4, 5]
head, *tail = nums  // head = 1, tail = [2, 3, 4, 5]
*most, last = nums  // most = [1, 2, 3, 4], last = 5
x, *middle, y = nums  // x = 1, middle = [2, 3, 4], y = 5

moarNums = [head, *middle, *tail] // [1, 2, 3, 4, 2, 3, 4, 5]
```

#### Testing For Membership
The `in` keyword can be used to easily test if a value is in a container:
```
l = [56, 78, 43]
if 56 in l {
	print("It's here")
}

if "i" in "team" {
	print("'i' IS in team!!")
}
```

#### Slicing
Rose supports slicing on lists, sets, strings, bytes, and tuples. Slicing is a convenient way to access multiple elements of a container at once:
```
msg = "Light the fire at midnight"
print(msg[10:4]) // fire
altMsg = msg[10:]
print(altMsg) // fire at midnight
```

#### Using make to Create Containers
Sometimes you know in advance roughly how large a container needs to be before you need to use it. In this case, the builtin function `make` can be used to create lists, sets and bytes with length and capacity hints, and maps with capacity hints. `make` takes at most 3 arguments: type of container to create, then the length and capacity of the new container.

`make` uses these rules to create containers:
1. When `make` is called with a non-zero length argument `n`, a container is created with `n` elements initialized to their zero value. The length and capacity of this container is `n`.
2. When `make` is called with a capacity argument `c` that is greater than the length argument `n`, a container is created with `n` initialized elements and allocated but uninitialized space for `c - n` elements.
3. When `make` is called to create a `map`, only one argument is accepted for the `map`'s capacity.
```
theWonderful101 = make(list[string], 101) // len = 101, cap = 101
badSet = make(set[float], 10, 0) // error: length is greater than capacity
preAlocd = make(list[int], 0, 1000) // len = 0, cap = 101
bigMap = make(map[string]string, 10000) // len = 0, cap = 10000
```

### How Are Types Passed Around?
Immutable types are copied when they are passed in an expression (assigned to another variable, passed into a function, etc). Mutable types' reference is copied when passed in an expression. This is a bit confusing, so here's some code to demonstrate:
```
str = "bleu"
str2 = str // str's value is copied into str2
str2 += " cheese"
print(str) // bleu

l = [1, 2, 3]
l2 = l // l2 and l now point to the same list
l2[0] = 0
print(l) // [0, 2, 3]
```

### Constants
Rose has two ways to declare constants: with `const` and with `let`. The difference is that variables declared with `const` must have their value known at *compile time*, while variables declared with `let` must be know at *runtime*.
```
import "math/rand"

const pi = 3.14159 // ok
let randBoi = rand.Intn(10) // ok
const badBoi = rand.Intn(15) // illegal: result of function not known at compile time
let unrecommendedBoi = 42 // ok, but not recommended, 42 is known at compile time
```

## Functions
Functions are declared with the `fn` keyword. Additionally, functions are first-class citizens in Rose. This means that functions can be assigned to variables, passed into functions as arguments, and returned from functions. Here's some simple code to demonstrate:
```
fn each(seq list[int], f fn(i int)) {
    for e in seq {
        f(e)
    }
}

fn sum(init int, seq list[int]) int {
    each(seq, fn(i int) { init += i })
    return init
}

print(sum(0, [1, 2, 3])) // 6
```
As you can see, function return types are *after* the function name and parameters. This makes lambda and closure declarations very easy to read. 

## TODO
- Switch statements
- Enums
