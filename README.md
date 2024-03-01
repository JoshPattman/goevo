# GoEvo Evolutionary Algorithms
<img src="./PACKAGE_ART/icon-background.svg" width=256 align="right"/>

GoEvo is an evolutionary algorithms package for Go that provides both flexible and fast implementations of many genetic algorithms.

Some Key Features:
- **Many Algorithms**: Support for many types of evolutionary algorithms, from basic hill-climbers to a full [NEAT](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf) implementation.
- **Optimise Anything**: NEAT genotypes, slices of floats, or any type that you can perform crossover, mutation, and fitness evaluation on are supported by this package.
- **Flexible for Your Use-Case**: As long as your components (such as mutation functions, selection functions, etc) implement the easy-to-understand interfaces specified, you can implement interesting and unique custom behaviour.

## Outline
In GoEvo, there is a clear separation between the optimisation algorithms, and the actual data to be optimised (known as a genotype). Any optimisation algorithm can be run on any genotype as long as the correction reproduction, selection, and fitness functions are available.

TODO: Write the rest of this section

## Documentation
TODO: Write a wiki

## Built-In Features List
### Algorithms
- [NEAT](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf) (Neuro Evolution of Augmenting Topologies)
- Simple one-species population
- One-by-one replacement

### Genotypes
- [NEAT](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf)
  	- [HyperNEAT](https://axon.cs.byu.edu/~dan/778/papers/NeuroEvolution/stanley3**.pdf) is also supported
- Float array
- Boolean array
