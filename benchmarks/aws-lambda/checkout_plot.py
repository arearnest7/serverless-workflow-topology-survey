import pandas as pd
import matplotlib.pyplot as plt

# Read the first CSV file
file1 = './online-shop-original/original_checkout.csv'  # Replace with your file name
column1_name = 'original_checkout_response-time'  # Replace with the column name you want to extract

data1 = pd.read_csv(file1)
column1 = data1[column1_name]

# Downsample the data - Example: taking every 10th data point
downsampled_column1 = column1[::10]  # Adjust the sampling rate as needed

# Read the second CSV file
file2 = './online-shop-partially-reduced/partial_checkout.csv'  # Replace with your file name
column2_name = 'partial_checkout_response-time'  # Replace with the column name you want to extract

data2 = pd.read_csv(file2)
column2 = data2[column2_name]

# Downsample the data - Example: taking every 10th data point
downsampled_column2 = column2[::10]  # Adjust the sampling rate as needed

# Plotting the downsampled data
plt.figure(figsize=(8, 6))  # Adjust the figure size if needed
plt.plot(downsampled_column1, label='Original', linestyle='-', linewidth=1)  # Plotting downsampled data from file1.csv
plt.plot(downsampled_column2, label='Partially-reduced', linestyle='-', linewidth=1)  # Plotting downsampled data from file2.csv

plt.xlabel('X-axis Label')  # Replace with your X-axis label
plt.ylabel('Y-axis Label')  # Replace with your Y-axis label
plt.title('Comparison of Columns (Downsampled)')  # Replace with your plot title
plt.legend()  # Show legend

plt.show()
