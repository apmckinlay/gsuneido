### Use try-catch to Check Type

Instead of, for example, using Object? or String? to check the type of a value, just attempt the operation.

For example:

``` suneido
Indexable?

function (x)
    {
    try
        x[0]
    catch (err)
        return err is "member not found: 0" ? true : false
    return true
    }
```

**Note**: In some cases, as in this one, you may have to use the exception to determine the result.

In Python this idiom is known as *Easier to Ask Forgiveness than Permission* (EAFP), as opposed to *Look Before You Leap* (LBYL).