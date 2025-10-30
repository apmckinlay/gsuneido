### Decorator

Add extra features to objects dynamically without modifying their class or using inheritance. Also known as *wrappers*. Decorators are more flexible than inheritance:

-	you can attach multiple decorators to the same object, whereas you can only inherit from one class
-	decorators are attached at run time, not just when you write the code


Decorator is similar to Composite since a decorator has the same interface as the object it wraps. But the purpose of decorators is to add features, whereas Composite is designed to structure groups of objects.

Uses:

-	ListViewModelCached is a decorator that wraps ListViewModel 
	and adds caching using LruCache


See Design Patterns by Gamma, Helm, Johnson, Vlissides.