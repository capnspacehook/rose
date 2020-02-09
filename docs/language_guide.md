# A Tour of Rose
## Introduction
Rose is a general-purpose, statically typed language that focuses on writeability, readability, safety, and simplicity (in that order).  It is designed to provide the expressiveness and elegance of dynamically typed languages, like Python, but keep the safety and maintainability of statically typed languages, such as Go. Rose has explicit support for concurrency, and aims to make error handling easy to write and understand. Rose is also garbage collected, so you don't have to worry about managing memory.

## Table of Contents
- [Basic Syntax](#basic-syntax)
  * [Hello World](#hello-world)
  * [Comments](#comments)
  * [Packages](#packages)
  * [Variables](#variables)
- [Logic Control](#logic-control)
  * [If Statements](#if-statements)
  * [Looping](#looping)
- [Types](#types)
  * [Basic Types](#basic-types)
    + [String Interpolation](#string-interpolation)
    + [The 'any' Type](#the--any--type)
  * [Zero Values](#zero-values)
  * [Container Types](#container-types)
    + [Lists](#lists)
    + [Maps](#maps)
    + [Sets](#sets)
    + [Tuples](#tuples)
    + [Iterating Over Containers](#iterating-over-containers)
    + [Container Unpacking](#container-unpacking)
    + [Testing For Membership](#testing-for-membership)
    + [Slicing](#slicing)
    + [Using make to Create Containers](#using-make-to-create-containers)
  * [How Are Types Passed Around?](#how-are-types-passed-around-)
  * [Constants](#constants)
- [Functions](#functions)
  * [Declaring Functions](#declaring-functions)
  * [Multiple Returns](#multiple-returns)
  * [Optional Parameters](#optional-parameters)
  * [Variadic Parameters](#variadic-parameters)
  * [Lambdas and Closures](#lambdas-and-closures)
  * [Defer Statements](#defer-statements)
  * [Goroutines](#goroutines)
- [TODO](#todo)

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
If statements have standard syntax in Rose. Curly braces `{}` are always required though. Else and else if clauses *must* start directly after the closing curly brace of the preceding if statement. Rose uses the keywords `and, or, not` as boolean operators, these can be used to create complex but readable boolean expressions.
```
imblue = true
if not imblue {
	print("I'm blue")
} else {
	print("...")
}

myColor = "blue"
if imblue and myColor == "blue" {
	print("Da ba dee da ba da")
} else if imblue and myColor == "green" {
	print("I'm turquoise")
```
In addition, Rose supports simple statements before the condition. This can be useful to reuse the returned value from a function multiple times. The variable is only available inside the scope of the if statement however.
```
import "math"

if s = math.Sqrt(16); if s > 0 and s == 4 {
	print("math doesn't lie")
}

print(s) // compile error: s doesn't exist in this scope
```


### Looping
Rose only has one kind of loop: the for loop. This may seem limiting, however Rose's for loop can handle all common looping types. A traditional for loop looks about how you'd expect it to, a for loop with one condition is essentially a while loop, and a for loop with no condition is an infinite loop.
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
`break` and `continue` can be used as normal to break out of a loop or skip an iteration:
```
// "1\n2" is printed
for i in range(5) {
	if i == 0 {
		continue
	} else if i == 3 {
		break
	}
	print(i)
}
```
For loops can also have an optional else clause that fires when the loop condition evaluates to false on the first iteration:
```
for 2 + 2 == 5 {
	print("math lied to me")
} else {
	print("I guess my professor was right...")
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
Rose also has some less commonly used, but still useful types: `char, byte, bytes`.
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

#### String Interpolation
Strings can be built and formatted very easily and concisely using string interpolation. When curly braces `{}` are present in string literals, their contents will be evaluated as an expression called a *string expression* and replaced with their values. Here's an example:
```
import "math"

name, age, location = "Andrew", 22, "the interwebs"
print("Hi, my name is {name}. I'm {age} years old currently at {location}.")

n = 56
print("The square root of {n} is {math.Sqrt(65)}")
```
If you want a literal curly brace in a string literal, you need to escape the curly brace with another curly brace, or just use a raw string literal, which doesn't support string expressions.
```
print("{{Braces}} {{{{galore}}}}") // {Braces} {{galore}}
print(`Here's some code: if isHungry { buyPizza() }`)
```
Strings are legal inside string expressions, but comments are not:
```
ppl = {"George": "Foreman", "Andy": "Sandberg", "Ronald": "McDonald"}
print("His surname is {ppl["Ronald]}")
print("This will not compile: {ppl // this is illegal}")
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
t = ['h', 'i', 8]
x, y = t

people = [
	("Bob", 42, "Mechanic"),
	("James", 24, "Artist"),
	("Harry", 32, "Lecturer")
]

for name, age, profession in people {
	print("Name: {name}; Age: {age}; Profession: {profession}")
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
Rose has two ways to declare constants: with `const` and with `let`. The difference is that variables declared with `const` must have their value known at *compile time*, while variables declared with `let` must be know at *runtime*. Declaring a variable as constant, with either `const` or `let` ensures that the value of the variable will not change during the program's lifetime.
```
import "math/rand"

const pi = 3.14159 // ok
let randBoi = rand.Intn(10) // ok
const badBoi = rand.Intn(15) // illegal: result of function not known at compile time
let unrecommendedBoi = "don't do this" // ok but not recommended, this string is known at compile time
```
Constants are very powerful in Rose, as *any type* can be a constant. Normally mutable types become immutable when declared as constants, and their methods and operators that normally would mutate their state are unavailable.
```
const philosophy = {"meaning of life": 42}
philosophy["the matrix"] = "real" // error: can't assign to constant map
print(philosophy["meaning of life"]) // 42
```
Additionally, it is important to note that when initializing a constant variable with a value `v`, and when assigning a constant value `v` to a mutable variable, the value is copied. Normally, if `v` is a mutable type, `v`'s reference would be copied. But to make sure constant values are separate and not modified, Rose copies `v` in both of these cases.

## Functions
### Declaring Functions
Functions are declared with the `fn` keyword. Additionally, functions are first-class citizens in Rose. This means that functions can be assigned to variables, passed into functions as arguments, and returned from functions.
```
fn fib(n int) int {
	if n == 1 {
		return 0
	} else if n == 2 {
		return 1
	}
	return fib(n+1) + fib(n+2)
}

print(fib(35)) // 9227465
```
As you can see, function return types are *after* the function name and parameters. This makes lambda and closure declarations very easy to read. 

In Rose, function parameter type and return type annotations are actually optional. If a function parameter is omitted, the compiler sets the parameter type to be `any`. This can be useful if you want to create concise lambdas or functions, or just want may different types to be able to be passed into a function.
```
fn printMe(me) {
	print(me)
}

printMe(5) // 5
printMe("print me!") // print me!
printMe([1, 6, 8]) // [1, 6, 8]
```
When return types are omitted however, Rose follows a few rules to determine what the return type of the function is:
1. If the function does not have any return statements, it is a void function.
2. If the function returns type `T` at least one time, the return type is `T`.
3. If a function returns multiple different types, the return type is `any`.

Here's some examples:
```
fn willHappen(fortune string) { // returns bool
	if len(fortune) % 2 == 0 {
		return true
	}
	return false
}

print(willHappen("I will be given a Cybertruck next Thursday")) // true

fn sayHi() { // void function
	print("...hi")
}

fn returnEVERYTHING(x) { // returns any
	if x > 4 {
		return 9
	} else if x < 4 {
		return "blue"
	} else {
		return {"red", "green", "yogurt"}
	}
}
```

### Multiple Returns
Functions can return multiple values by returning a tuple. Then the returned tuple can be unpacked for convenient access to the multiple return values.
```
fn split(s string) (string, string) {
	middleLen = len(s) / 2
	return s[:middleLen], s[middleLen:]
}

first, last = split("hello there charming")
```
I know what you're probably thinking , and yes, Rose can still infer what the function return type is if you omit a return type annotation and return multiple values. 

### Optional Parameters
Optional parameters are parameters that are, well, optional. They have a default value, and will use the default value if one is not explicitly passed. Parameters that are not optional are called *positional* parameters.
```
fn printChars(s string, reverse bool=false) {
	if reverse {
		s.reverse()
	}
	
	for c in s {
		print(c)
	}
}

alpha = "abc"
printChars(alpha) // prints "a\nb\nc"
printChars(alpha, true) // prints "c\nb\na"
printChars(alpha, reverse=true) // also prints "c\nb\na"
```
Optional parameters *must* come after positional parameters, both when declaring a function and calling it. If a type for a optional parameter is not given, it will be inferred from the default value.

### Variadic Parameters
Variadic parameters are parameters that accept any number of arguments. They are created by prepending ellipses `...` to the parameter name. There can only be *one* variadic parameter per function, and it must be the *last* parameter. Inside the function, all the arguments passed to the variadic parameter is accessible as a list.
```
fn fullFunc(x, y=2, ...z) { // ok, variadic parameter is last
	print("I compile :)")
}

fn badFunc(...stuff, things) { // illegal: variadic parameter is not last
	print("I won't compile")
}

fn variadicFunc(...args) { // ok, variadic parameter is the only parameter
	print(len(args))
}

variadicFunc(1, 2, 3, 'f', "john") // prints "5"
variadicFunc() // prints "0"
```

### Lambdas and Closures
Rose has great support for lambdas, or anonymous functions, and closures, which are special lambdas which "close" around variable(s) in the scope that it is declared in.
```
fn each(seq list[int], f fn(i int)) {
    for e in seq {
        f(e)
    }
}

fn sum(seq list[int], init=0) int {
    each(seq, fn(i) { init += i })
    return init
}

print(sum([1, 2, 3])) // 6
```
If a lambda or closure doesn't have any parameters, the parentheses can be omitted entirely for conciseness: 
```
f = fn { print("I'm a super small lambda") }
```

### Defer Statements
A defer statement defers the execution of a function until the surrounding function returns. Defer statements are evaluated immediately, but the deferred function is not run until the surrounding function returns. Additionally, deferred functions are guaranteed to execute so long as the surrounding function returns and the program doesn't completely crash.
```
fn aTwist() {
	defer fn { print("world") }
	print("hello")
}

aTwist() // hello world
```
Deferred functions are pushed onto a stack, and therefore run in last-in-first-out order.
```
fn testDefers() {
	defer fn { print("I ran 3rd") }
	defer fn { print("I ran 2nd") }
	defer fn { print("I ran 1st") }
}

testDefers() // I ran 1st\nI ran 2nd\nI ran 3rd
```

### Goroutines
All you have to do to run a function asynchronously in Rose is to prepend the `go` keyword to a function call; that's it. It's that simple. Rose has no way of managing or tracking goroutines, so be sure to make sure they eventually finish.
Imagine we have a function `isPrime` that takes an integer and returns true if it's prime or not. Given large inputs, this can take awhile. Luckily we can easily check multiple numbers at once easily with goroutines:
```
bigBois = [100092764521462167, 83756379996131673, 473829456141315, 8376493623621153311]
for boi in bigBois {
	go print(isPrime(boi))
}
```

## Advanced Types

## TODO
- Advanced Types
	- errors
	- structs
	- optionals
	- results
	- channels
	- interfaces
- Switch statements
- Enums
- Iterators
- Generators
	- yield statements
