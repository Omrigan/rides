# Taxi rides
Task: https://querifylabs.notion.site/Analyzing-taxi-rides-3b4c61decf3f474da9ee90211fd3973d

## Implementation

Project includes two index versions:
 - Dummy, which uses the example algorithm from the task
 - Binsearch, where:
   - An array with all rides ordered by its start time
   - A binary search is performed to limit search space to _candidates_ - rides which only started in given range
   - Then a full scan of the remaining rides is filtered to leave only rides which ended in given range
   - Full scan can be parallelized
 
## Tradeoffs
Taking into account the relative smallness of the dataset, it is reasonable to keep all data in memory. 

In case of a larger dataset, some disk-based B-tree or LSM could be used.

## Possible improvements
The binsearch algorithm can be further improved by having a log-sized array attached to every ride.

The idea here is being able to skip ahead during the scan through _candidates_.

i-th element of this array corresponds to next `2^i` rides, following this particular ride.
In this element is stored `(max_of_finish_times, sum, cnt)`.

So, if in a particular request `(start, finish)` on the first ride we see `max_of_finish_times < finish`, 
this means we can skip ahead `2^i` rides, which are all before `finish`.
## Building instructions
```bash
cd main/
go build
```

 
