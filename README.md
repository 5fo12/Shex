# Shex scripting langauge
Shex is a small interpreted language.  
It currently is still being updated. 

To run a script with the compiled version:
```
shex scriptname.sx
```

You can also compile your scripts to binaries!
Just use the compiler binary:
```
compile scriptname.sx
```

You can build it yourself by running:  
```
go build shex.go
```

### Comments
Single line comments start with "!".  
Multiline comments start with "<" and end with ">".
```
! This is a comment.

<
    This is 
    a comment.
>
```
### Values
Text values are surronded in double quotes.  
Number values stay as is.
Math operations are surronded by parenthesis.    
Evaluations are surronded by curly brackets.   
Conditions are surronded by brackets.  

**For Conditons, Evaluations and Math operations each part must be seperated by a space**. 

```
"Hello, World!"
50
(43 + 7)
{function param1 param2}
[5 < 2]
```

### Commands
**Var** declares a variable, variables can be accessed in strings by surrounding the variable name in dollar signs.  
Variable values are bool, text and number.

```
! Usage: var name:type
var var1:text

"Var1: $var1$"
```  

**Set** changes the value of a variable.  

```
! Usage: set name value
set var1 "String"
set var1 (5 + 2)
set var1 {function param1 param2}
! Conditions cannot be used for variables yet.
```

**Print** outputs text.

```
! Usage: print value
print "Var1 value: $var1$"
print (5 + 5)
print {function param1 param2}
! Conditions cannot be used for printing yet.
```

**Read** allows you to take use input.  

```
! Usage: read
var input:text
print "Input text:"
set input {read}

print "Input: $input$"
```

**If** allows the execution of code if a condition is true.

```
! Usage:
if [5 < 10]
print "5 is less than 10."
endif

! Use endif to end if blocks.
```

**Func** allows the creation of functions for reusable code.
Parameters are declared the same way as variables.  

```
! Usage:
func greet name:text
print "Hello $name$!"
endfunc

! Use endfunc to end function blocks
```

## Restrictions
Shex currently is very strict on syntax,  
here are some things to know:

Recursive searching of if and function blocks is not implemented.  
Do not indent code in if or function blocks.  
Do not leave whitespace after commands.  

Breaking any of these rules will lead to errors.