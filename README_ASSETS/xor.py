from math import *

def sigmoid(x):
    return 1 / (1 + exp(-x))

def xor(N0, N1):
    N508 = cos(-3*N1 + 3*N0)
    N301 = tanh(-3*N508)
    N247 = tanh(-3*N508 + 2.672*N301)
    N4 = sigmoid(-3*N508 + 3*N247 + 3*N301)
    return N4

for (N0, N1) in [(0, 0), (0, 1), (1, 0), (1, 1)]:
    print(f"XOR({N0}, {N1}) = {xor(N0, N1)}")

"""
OUTPUT:
XOR(0, 0) = 0.00012524777618418145
XOR(0, 1) = 0.9998708166272475
XOR(1, 0) = 0.9998708166272475
XOR(1, 1) = 0.00012524777618418145
"""