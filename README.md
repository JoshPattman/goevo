# GoEvo Evolutionary Algorithms
<img src="./PACKAGE_ART/icon-background.svg" width=256 align="right"/>

GoEvo is an evolutionary algorithms package for Go that provides both flexible and fast implementations of many genetic algorithms.

Some Key Features:
- **Many Algorithms**: Support for many types of evolutionary algorithms, from basic hill-climbers to a full [NEAT](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf) implementation.
- **Optimise Anything**: NEAT genotypes, slices of floats, or any type that you can perform crossover, mutation, and fitness evaluation on are supported by this package.
- **Flexible for Your Use-Case**: As long as your components (such as mutation functions, selection functions, etc) implement the easy-to-understand interfaces specified, you can implement interesting and unique custom behaviour.

## Documentation
The documentation is stored on the [GoEvo Wiki](https://github.com/JoshPattman/goevo/wiki).

## Overview of Structure
The parent module, `goevo`, defines a set of interfaces that lay out the behavior for which the different components must implement. However, the parent module does not implement any of these interfaces, but instead, some default implementations are stored in the sub-modules. As long as you can provide implementations for the following interfaces, it is possible to use any combination of implementations to create a customized evolutionary simulation:

- `Population[T]`: Represents a population storing a set of genotypes of type `T`. All default implementations of this interface are stored in `pop/`.
- `MutationStrategy[T]`: Represents a mutation strategy that can be applied to a genotype of type `T` to perform a mutation. The default implementations reside in `geno/`. You can have multiple mutation strategies for any given genotype type.
- `CrossoverStrategy[T]`: Represents a crossover strategy that can be applied to two genotypes of type `T` to perform a crossover (this can have any number of parents). The default implementations reside in `geno/`. You can have multiple crossover strategies for any given genotype type.
- `Selection[T]`: Represents a selection strategy that can be applied to a population of genotypes of type `T` to select a subset of the population. The default implementations reside in `selec/`. You can have multiple selection strategies for any given genotype type.
- Genotypes: These do not need to implement any interfaces themselves. However, they are the DNA upon which the crossover and mutation strategies are applied. The default implementations reside in `geno/`.

It is also entirely possible to only use certain components of the package. For example, a continuous evolutionary simulation may not require a `Selection` or `Population`, as these are inherited from interactions with the environment.

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
