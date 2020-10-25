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
