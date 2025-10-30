### Sequence

A sequence is a virtual list that wraps an iterator.

Several built-in methods return sequences - [object.Values](<../Reference/Object/object.Values.md>), [object.Members](<../Reference/Object/object.Members.md>), [object.Assocs](<../Reference/Object/object.Assocs.md>), and [Seq](<../Reference/Seq.md>) and the [Sequence](<../Reference/Sequence.md>) function can be used to create a sequence from a user defined iterator.

Additional, user defined methods can be added by defining a class called "Sequences" which should inherit from "Objects". The main use is to provide sequence specific overrides to Objects methods that do not force instantiation.

Sequences have built-in Size and Join methods that avoid instantiating the sequence.

Sequences can be iterated over (e.g. with for-in) without instantiating an actual list.

For example, something like:

``` suneido
list.Filter(...).Map(...).Join()
```

will "stream" the values and will not create any intermediate list objects.

However, if any built-in Object methods are used on the sequence then the list will be instantiated.