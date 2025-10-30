## Test If a Value is One of a List

**Category:** Coding

**Problem**

You need to check if a value is one of several different choices.

**Ingredients**

[if](<../Language/Statements/if.md>), 
[Regular Expressions](<../Language/Regular Expressions.md>), 
[object.Has?](<../Language/Reference/Object/object.Has?.md>), 
[object.Member?](<../Language/Reference/Object/object.Member?.md>), 
[switch](<../Language/Statements/switch.md>)

**Recipe**

The most obvious solution is to simply use **if** and **or**, for example:

``` suneido
if (value is 'red' or value is 'green' or value is 'blue')
```

However, if you have more than a few choices, or if you need to do this test several times, it soon becomes awkward.

Another option is to use regular expressions:

``` suneido
if (value =~ "^(red|green|blue)$")
```

The ^(...)$ ensure that you are matching the entire string, not just a portion.

This option will only work for strings. For other types of values, you can use an object to contain the list:

``` suneido
if (#('red', 'green', 'blue').Has?(value))
```

For simple values (e.g. strings, numbers, dates), it's faster to make the values the members (keys) of the object:

``` suneido
if (#('red':, 'green':, 'blue':).Member?(value))
```

This is faster because it takes advantage of the hash lookup that object use to find members. Notice that we didn't supply any values for the members - the default value is **true**.

Another option is to use a **switch**:

``` suneido
switch (value)
    {
case 'red', 'green', 'blue':
    ...
    }
```

This is usually not the best choice unless you need to test for several different cases.

**Discussion**

For a small number of string values I would use a regular expression. My next choice for simple values would be .Member?. Otherwise, the most general solution is to use .Has?.