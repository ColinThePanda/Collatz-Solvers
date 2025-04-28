import random

def is_even(num : int) -> bool:
    return num % 2 == 0

def colatz(num : int) -> int:
    if is_even(num):
        return num//2
    else:
        return num * 3 + 1

nums = []
starting_num : int = 2**10000 - 1
num : int = starting_num
print(num)
while num != 1:
    num = colatz(num)
    nums.append(str(num))
    print(num)
print(f"done, len={len(nums)}")
filename = f"collatz_conjecture_{starting_num if len(str(starting_num)) < 235 else str(starting_num)[:187]}"
with open(f"{filename}.txt", "w") as file:
    file.write("\n".join(nums))