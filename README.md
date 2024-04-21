# Warning
This package is currently under a period of intense development and should be considered extremely unstable until version `v0.5.0`.

## TODO pre `v0.5.0`
- NEAT Population
- Evolutionary tree diagram
- N point crossover on floatarr
- Convert floatarr to arr (provide multiple mutators for floats, chars, etc. Should be able to use the same crossovers for all)
- Add more tests and clean up test dir
- Add an examples dir
- Make wiki


# GoEvo Evolutionary Algorithms
<img src="./PACKAGE_ART/icon-background.svg" width=256 align="right"/>

GoEvo is an evolutionary algorithms package for Go that provides both flexible and fast implementations of many genetic algorithms.

Some Key Features:
- **Many Algorithms**: Support for many types of evolutionary algorithms, from basic hill-climbers to full [NEAT](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf).
- **Optimize Anything**: NEAT genotypes, slices of floats, or any type that you can perform crossover, mutation, and fitness evaluation on are supported by this package.
- **Flexible for Your Use-Case**: As long as your components (such as mutation functions, selection functions, etc) implement the easy-to-understand interfaces specified, you can implement interesting and unique custom behavior.

## Documentation
The documentation is stored on the [GoEvo Wiki](https://github.com/JoshPattman/goevo/wiki).

## Built-In Features List
### Algorithms
- [NEAT](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf) (Neuro Evolution of Augmenting Topologies)
- Simple one-species population
- One-by-one replacement

### Genotypes
- [NEAT](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf)
  	- [HyperNEAT](https://axon.cs.byu.edu/~dan/778/papers/NeuroEvolution/stanley3**.pdf) is also supported
- Float array
