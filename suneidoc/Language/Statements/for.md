### for

There are several forms of the **for** statement.

#### for i in from..to

Simple counted loops that increment by one.
<pre>or i in <i>from</i>..<i>to</i></pre>

Is equivalent to:
<pre>or (i = <i>from</i>; i &lt; <i>to</i>; ++i)</pre>

*from* may be omitted if it is 0

``` suneido
for i in ..to
```

To just do a loop a certain number of times, without a variable, use:

``` suneido
for ..to
```

Note: Unlike the classic form, the limit expression (*to*) is only evaluated once, not once per loop.

Note: Unlike the classic form, modifying the loop variable will **not** affect the iteration.

#### for x in y

To loop through the values of any *iterable* object.
<pre>
for <i>variable</i> in <i>expression</i>
    <i>statement</i>
</pre>

This form is equivalent to:
<pre>
iter = <i>expression</i>.Iter()
while (iter isnt <i>variable</i> = iter.Next())
    <i>statement</i>
</pre>

Note: If members are added or removed during iteration, an "object modified during iteration" exception will be thrown. This can usually be "fixed" by making a copy of the members, for example:

``` suneido
for m in ob.Members().Copy()
    if ...
        ob.Delete(m)
```

The expression must evaluate to an object that supports .Iter() which must return an iterator that has a .Next() method that returns the iterator itself when it reaches the end.

#### for m,v in ob

With objects, records, classes, or instances you can use:
<pre>
for <i>var1</i>, <i>var2</i> in x
</pre>

where var1 will be the index or member name and var2 will be the value.

This is "snapshot" iteration. i.e. it will not "see" changes made during iteration. It is not necessary to use .Copy()

#### for (init; cond; inc)

The "classic" form is similar to C/C++/Java.
<pre>
for ( <i>expressions1</i> ; <i>logical-expression</i> ; <i>expressions2</i> )
    <i>statement</i>
</pre>

Executes *expressions1* once, and then repeatedly executes the *statement* followed by *expressions2* as long as the *logical-expression* is true. The *statement* and *expressions2* will not be executed at all if the *logical-expression* is initially false.

This form is equivalent to:
<pre>
<i>expressions1</i>
while ( <i>logical-expression</i> )
    {
    <i>statement</i>
    <i>expressions2</i> 
    }
</pre>

*expressions1* and *expressions2* can consist of zero or more comma separated expressions

For example:

``` suneido
for (i = 0; i < 10; ++i)
    Print(i)
```

If the *logical-expression* is omitted, it is taken as true, presumably with some other provision for ending the loop.

Note: If logical-expression evaluates to something other than true or false, an exception will result: "conditionals require true or false".

See also: [forever](<forever.md>)