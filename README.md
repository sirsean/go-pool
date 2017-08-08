# go-pool

This library provides a very simple work pool based on the concept of
"work units".

The pool has a "queue size" and a "number of workers", which effectively
describe how many work units can be buffered and waiting to execute before
adding new work units will block, and how many goroutines will be created
to perform work.

![go pool](https://frinkiac.com/meme/S06E01/195028.jpg?b64lines=SSdNIEdPTk5BCiBTVE9XQVdBWSBVTkRFUldBVEVSLAogQU5EIEdPIFdIRVJFIFRIRSBQT09MIEdPRVM=)

# Usage

```
import (
    "github.com/sirsean/go-pool"
)
```

Create a new pool, with a queue size and number of workers:

```
p := pool.NewPool(10, 20)
```

Start it (this is non-blocking):

```
p.Start()
```

Add work units:

```
p.Add(workUnit1)
p.Add(workUnit2)
// ... as many times as you want
```

Wait until all the work is done:

```
p.Close()
```

This will close the channel and then wait for all the existing work units to be
completed. At this point the pool will no longer accept new work.

# Work Units

But how do you specify what work should happen?

There is an interface called `WorkUnit` with a `Perform()` function defined. You
will want to implement this interface like so:

```
type MyWorkUnit struct {
    Something string
}

func (u MyWorkUnit) Perform() {
    log.Println("I am doing time-consuming work on", u.Something)
}

// ...

p.Add(MyWorkUnit{
    Something: "something",
})
```

A `Pool` can take different kinds of `WorkUnit`s and will execute them.

If you need to get data _out_ of the work unit, you can build one that has a
channel and read from that channel, but that's an implementation detail of
your program.

# License

This project uses the Simplified BSD License.

# Code of Conduct

Please note that this project is released with a Contributor Code of Conduct.
By participating in this project you agree to abide by its terms.
