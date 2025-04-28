import sys
import os

sys.set_int_max_str_digits(10000)


def collatz(num: int):
    if num % 2 == 0:
        return num // 2
    else:
        return num * 3 + 1


num = 2**10000 - 1
starting_num = num
nums: list[int] = []
while num != 1:
    nums.append(num)
    num = collatz(num)
    print(num)

print("Wait for it to write to a file...")

with open(f"collatz_conjecture{str(starting_num)[:10]}" + "...", "w") as file:
    file.write("\n".join([str(num) for num in nums]))

print("Done writing!")