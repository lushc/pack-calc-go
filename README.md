# Pack Calc (Go)

A solution for calculating the number of packs required to satisfy a requested quantity of items while abiding by the following rules:

1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of rule 1, send out no more items than necessary to fulfil the order.
3. Within the constraints of rule 1 & 2, send out as few packs as possible to fulfil each order.

Derived from my prior implementation [pack-calc-php](https://github.com/lushc/pack-calc-php).

## Getting started

Build the image:

```
docker build -t pack-calc-go .
```

Run the dev container with hot reload testing:

```
docker run --name pack-calc-go-dev -v $(pwd):/app pack-calc-go
```

Reference the container name with `docker exec` to run arbitary commands, e.g. format all files:

```
docker exec -it pack-calc-go-dev gofmt -s -w .
```

## Benchmarks

```
BenchmarkGraphPackCalculator_Calculate/default_packs_with_1-12                             40418             29614 ns/op           19297 B/op        305 allocs/op
BenchmarkGraphPackCalculator_Calculate/default_packs_with_250-12                           40042             29918 ns/op           19554 B/op        309 allocs/op
BenchmarkGraphPackCalculator_Calculate/default_packs_with_251-12                           30428             39265 ns/op           24827 B/op        367 allocs/op
BenchmarkGraphPackCalculator_Calculate/default_packs_with_501-12                           23762             50416 ns/op           30924 B/op        476 allocs/op
BenchmarkGraphPackCalculator_Calculate/default_packs_with_12001-12                          1866            643788 ns/op          324358 B/op       7207 allocs/op
BenchmarkGraphPackCalculator_Calculate/default_packs_with_47501043056-12                      44          25187338 ns/op        11692241 B/op     278826 allocs/op
BenchmarkGraphPackCalculator_Calculate/prime_packs_with_32-12                              23742             50508 ns/op           32424 B/op        443 allocs/op
BenchmarkGraphPackCalculator_Calculate/prime_packs_with_500-12                               232           5101087 ns/op         2689115 B/op      43820 allocs/op
BenchmarkGraphPackCalculator_Calculate/prime_packs_with_758-12                               100          10396431 ns/op         4916666 B/op      94367 allocs/op
BenchmarkGraphPackCalculator_Calculate/off_by_one_pack_with_500-12                           358           3344326 ns/op         1950087 B/op      27037 allocs/op
BenchmarkGraphPackCalculator_Calculate/edge_case_pack_permutation-12                        8546            139276 ns/op           82748 B/op       1346 allocs/op
BenchmarkGraphPackCalculator_Calculate/choose_smallest_pack_count-12                         120          10056510 ns/op         4630608 B/op      90139 allocs/op
BenchmarkGraphPackCalculator_Calculate/prime_stress_test-12                                   88          13814504 ns/op         6747548 B/op     127597 allocs/op
BenchmarkGraphPackCalculator_Calculate/prime_stress_test_with_3_sizes-12                     667           1793968 ns/op         1106331 B/op      16703 allocs/op
BenchmarkGraphPackCalculator_Calculate/prime_stress_test_with_2_sizes-12                     361           3298535 ns/op         1814290 B/op      32965 allocs/op
BenchmarkSimplePackCalculator_Calculate/zero_quantity-12                                   77960             15231 ns/op           10822 B/op        139 allocs/op
BenchmarkSimplePackCalculator_Calculate/negative_quantity-12                               78458             15038 ns/op           10821 B/op        139 allocs/op
BenchmarkSimplePackCalculator_Calculate/single_pack-12                                     77690             15082 ns/op           10822 B/op        139 allocs/op
BenchmarkSimplePackCalculator_Calculate/divisible-12                                       77738             15116 ns/op           10822 B/op        139 allocs/op
BenchmarkSimplePackCalculator_Calculate/indivisible-12                                     77396             15165 ns/op           10823 B/op        139 allocs/op
```

## Microservice deployment

The [Serverless framework](https://serverless.com/) is used for a microservice deployment:

```
npm install -g serverless
serverless config credentials --provider aws --key <key> --secret <secret>
```

To build and deploy the function:

```
docker exec -it pack-calc-go-dev make build
make deploy
```

The API Gateway endpoint and API key will then be printed to console. An example request would be:

```
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-Api-Key: example_key" \
  -d '{"quantity": 12001, "packSizes": [250,500,1000,2000,5000]}' \
  https://example.execute-api.eu-west-1.amazonaws.com/dev/calculate
```

The service responds with a JSON object describing the number of required pack sizes:

```
{"250":1,"2000":1,"5000":2}
```

## Implementation

### Algorithm

1. Initialise a directed multigraph where nodes are quantities and edges are pack sizes
   - When the quantity exceeds an arbitary threshold (sum of pack sizes * 50) we first reduce the problem space by subtracting as many of the largest packs as possible while still leaving enough headroom to permutate a best fit
   - In the case of a single pack size, a graph isn't necessary to calculate the required packs and so a separate implementation is used instead
2. Recursively build the graph by subtracting pack sizes from ancestors, starting from the root node (initial quantity)
   - Packs are subtracted from the current node's quantity in descending order
   - Either a new node is created for the calculated quantity or an existing node is located
   - A weighted edge between the current node and new node is created to track the pack size used (i.e. a new permutation)
   - Nodes with a quantity <= 0 are treated as a candidate and no further subtraction occurs
   - Nodes with a quantity > 0 continue to recurse
   - Permutation generation is halted when a number of paths to 0 are found to prevent an exhaustive and expensive search
   - The available pack sizes are reduced on each iteration over the root as this helps produce different permutations
3. Candidate nodes are sorted (by quantity) descending, with the first being chosen as it's either 0 or closest to 0
4. The graph is pruned to remove nodes that are of a lower quantity than the chosen candidate node
5. The graph is pruned further to remove other nodes which don't have any outgoing edges (i.e. they're a dead end)
6. Graph traversal is performed to find the shortest path between the root node and the candidate node
7. Each edge in the path is iterated over, using their weight to tally the number of packs used at each size
8. A map is returned where keys are pack sizes and values are the counts

### Remarks

These are my thoughts on how the implementation can be improved or a different approach could be taken.

#### Graph traversal

Although the search algorithm used is A*, the graph is currently reporting a uniform cost for edge weights and no heuristic function is in place for performing a more informed search. Taking the pack size weight at face value meant larger pack edges were treated as being a higher cost towards the goal (greater distance, more moves requried etc.) and this was reflected in the results where smaller packs were favoured, meaning less efficient pack size configurations.

I did try to implement some heuristics so that larger weights were more favoured but this did not affect the results as I anticipated, so more work could be done to develop a better best-fit search.

#### Non-deterministic results

The graph calculator does not always produce the same result. This can be seen when running tests as sometimes a test can fail, for example:

```
graph_test.go:48: {PackSizes:[23 31 53 151 757]}.Calculate(758) == map[31:3 53:4 151:3], want map[23:4 31:2 151:4]
```

Only a subset of tests are affected and the total number of packs are equal in terms of count and quantity. I suspect this is due to the order of iteration over nodes not being stable.

#### Memory usage

Allocation of memory is quite high when using very large pack sizes and quantities, for example:

``
{PackSizes:[250 500 1000 2000 5000 10000 20000 50000]}.Calculate(5000001) Alloc = 122 MiB
``
