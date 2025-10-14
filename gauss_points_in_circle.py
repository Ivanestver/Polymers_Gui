import math
import sys

def gauss_circle_optimized(R):
    """Оптимизированная версия метода Гаусса"""
    if R == 0:
        return 1
    
    count = 0
    R_sq = R * R
    limit = int(R)
    
    # Считаем только первый квадрант (включая оси)
    for x in range(0, limit + 1):
        y_max = int(math.sqrt(R_sq - x*x))
        count += y_max + 1  # +1 для точки на оси y=0
    
    # Умножаем на 4 и вычитаем лишние точки на осях
    return 4 * count - 4 * limit - 3

if len(sys.argv) != 2:
    print(f"Usage: {sys.argv[0]} R")
    exit(0)

print(gauss_circle_optimized(float(sys.argv[1])))