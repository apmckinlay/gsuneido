### Number

Suneido has a single numeric type that is used to support both integers and floating point numbers. Numbers use a decimal representation in order to handle decimal fractions exactly.

Support arithmetic and bitwise operations.

Additional, user defined methods can be added by defining a class called "Numbers".

Numeric literals can be written in the following formats:

decimal:

``` suneido
an optional leading sign (+ or -)

one or more digits with an optional decimal place (.)

an optional exponent consisting of:

    the letter 'e' or 'E'

    an optional leading sign (+ or -)

    one or two digits

For example: 123, +1, -1, 1.23, .001, 1e6, 1e-2
```

hex:

``` suneido
0x followed by one or more hex digits (0-9, a-f, A-F)

For example: 0xc4
```

**Note:** Currently hex values are limited to 32 bits.

Numbers have a minimum of 13 digits of precisions and a maximum of 16 digits. Integers have 16 digits of precision.