import random
import os

def is_even(num : int) -> bool:
    return num % 2 == 0

def colatz(num : int) -> int:
    if is_even(num):
        return num//2
    else:
        return num * 3 + 1

def get_filename(num: int) -> str:
    """Generate the filename for a given number"""
    return f"collatz/collatz_conjecture_{num if len(str(num)) < 235 else str(num)[:187]}.txt"

def calculate(num: int):
    """Calculate the Collatz sequence for a number, using previously saved files when possible"""
    # Check if file already exists for this starting number
    if os.path.exists(get_filename(num)):
        print(f"Skipping {num} - already calculated")
        return
        
    starting_num = num
    nums = []
    print(num)
    
    while num != 1:
        # Check if we have already calculated the sequence for this number
        current_filename = get_filename(num)
        if num != starting_num and os.path.exists(current_filename):
            print(f"Found file for {num}, using saved sequence")
            # Read the sequence from the file
            with open(current_filename, "r") as file:
                rest_of_sequence = file.read().strip().split("\n")
                nums.extend(rest_of_sequence)
                print(f"Added {len(rest_of_sequence)} numbers from file")
            break
        
        # Otherwise, calculate the next number
        num = colatz(num)
        nums.append(str(num))
    
    print(f"done, len={len(nums)}")
    
    # Ensure the collatz directory exists
    os.makedirs("collatz", exist_ok=True)
    
    # Save to file
    with open(get_filename(starting_num), "w") as file:
        file.write("\n".join(nums))

if __name__ == "__main__":
    # Create the collatz directory if it doesn't exist
    os.makedirs("collatz", exist_ok=True)
    
    for i in range(2, 100001):
        calculate(i)