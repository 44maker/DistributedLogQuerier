import numpy as np
import matplotlib.pyplot as plt

# 提取每次运行的耗时数据
elapsed_times = [
    35.334544,
    42.762791,
    41.1956,
    45.642704,
    50.298383
]

# 计算平均值和标准偏差
mean_elapsed_time = np.mean(elapsed_times)
std_elapsed_time = np.std(elapsed_times)

# 创建两个子图
fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(10, 5))

# 绘制平均查询时间柱状图
ax1.bar(range(1, len(elapsed_times) + 1), elapsed_times)
ax1.set_xticks(range(1, len(elapsed_times) + 1))
ax1.set_xticklabels([f"Run {i}" for i in range(1, len(elapsed_times) + 1)])
ax1.set_ylabel("Time (ms)")
ax1.set_title("Average Query Time for Each Run")

# 绘制标准偏差柱状图
ax2.bar(1, mean_elapsed_time, yerr=std_elapsed_time, capsize=10)
ax2.set_xticks([1])
ax2.set_xticklabels(["Average Elapsed Time"])
ax2.set_ylabel("Time (ms)")
ax2.set_title("Average Elapsed Time with Standard Deviation")

# 调整子图间距
plt.subplots_adjust(wspace=0.5)

plt.show()