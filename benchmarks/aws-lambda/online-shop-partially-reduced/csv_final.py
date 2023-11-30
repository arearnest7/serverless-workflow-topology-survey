import pandas as pd

# List of CSV files to combine
file_list = ['./checkout.csv']

# Read each CSV and concatenate
combined_data = pd.concat([pd.read_csv(file) for file in file_list])

# Adding 'original_' prefix to column names
combined_data.columns = ['partial_checkout_' + col for col in combined_data.columns]

# Write the combined data to a new CSV file
combined_data.to_csv('partial_checkout.csv', index=False)


file_list = ['./checkout.csv', './cart.csv','./ad.csv','./list.csv']

# Read each CSV and concatenate
combined_data = pd.concat([pd.read_csv(file) for file in file_list])

# Adding 'original_' prefix to column names
combined_data.columns = ['partial_' + col for col in combined_data.columns]

# Write the combined data to a new CSV file
combined_data.to_csv('partial_final.csv', index=False)