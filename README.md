# GoEvo Rewrite
This branch is an experimental rewrite from the ground up. I am making it as fast, simple, and bug-free as possible.

# GoEvo - Evolutionary Algorithms Based on NEAT, in Golang
GoEvo is a package for Go that performs a variety of evolutionary algorithms to optimise a neural network. The networks do not only evolve their weights, but also the neurons within, and therefore their whole structure. I have tried to base GoEvo networks on the ones described in the original NEAT paper, [Evolving Neural Networks Through Augmenting Topologies](https://nn.cs.utexas.edu/downloads/papers/stanley.ec02.pdf), however, I have taken some liberties during development. Despite this, the networks function in broadly the same way as originally described.

GoEvo also does not just support the NEAT algorithm. NEAT genotypes can be used with many other types of genetic algorithms, with various methods of selection, mutation, and speciation. For this reason, GoEvo strives to support many different algorithms, all focusing on optimising the Genotype type, which equates to DNA. GoEvo is also built to be very extensible, allowing the user to implement new evolutionary processes with ease.

## Future Work
- HyperNEAT: I already implemented this in the original package before I rewrote it, so this should be trivial. HyperNEAT basically allows NEAT to evolve much larger networks.
- Implement NEAT algorithm: This should be coming very soon, as I just need to convert the code from the original into this package. Currently there is a simple population.
- Document EVERYTHING
- Examples directory