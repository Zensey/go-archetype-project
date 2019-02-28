# How to run

* In docker: `make docker-build`
* Locally: `make test`

# Rem.

The time complexity of getMedian() is O(log n).
It can be reduced to O(1) by using two AVL trees: left -- for smaller half and right for greater half.
We must also keep length of trees different at most by 1. 
If length are equal median can be calculated like this : 
`median := ( left.max() + right.min() ) / 2`
To make run-time of left.max() / right.min() constant we can cache them.   