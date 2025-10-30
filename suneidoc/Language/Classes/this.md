### this

Within class member functions the keyword `this` refers to the object (class instance) that the member function was called from.

Note: `this` may not be assigned to.

Within class methods, `this.member` may be abbreviated to just `.member`

Note: Because private member names are prefixed with the name of their class

``` suneido
this["private_member"]
```

will not work. (The actual member name will be of the form "Class_private_member".)