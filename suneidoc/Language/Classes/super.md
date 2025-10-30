### super

Within a class method `super.fn(...)` can be used to call a method in a base class that has been overridden in the current class.

Within a `New` method `super(...)` can be used to pass arguments to the base class `New`

**Note: An explicit `super(...)` call must be the first statement of the `New` method.**

If there is no explicit `super(...)` call, then an implicit `super()` call will be generated.