import random

F = open('biglist.txt', 'w+')

for i in range(100000000):
    F.write(str(random.randrange(2147483647)) + '\n')
