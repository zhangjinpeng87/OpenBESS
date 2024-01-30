import pandas as pd
from sklearn.cluster import DBSCAN
from sklearn.preprocessing import StandardScaler
import matplotlib.pyplot as plt

# Load data from a CSV file
file_path = 'normalized-mvotage-mcurrent.csv'
data = pd.read_csv(file_path)

# Assuming your data is in a DataFrame and you want to use all columns for clustering
X = data.values

# Standardize the features
scaler = StandardScaler()
X_std = scaler.fit_transform(X)

# Set the parameters for DBSCAN
eps = 0.5  # epsilon, the maximum distance between two samples for one to be considered as in the neighborhood of the other
min_samples = 10  # the number of samples (or total weight) in a neighborhood for a point to be considered as a core point

# Perform DBSCAN clustering
dbscan = DBSCAN(eps=eps, min_samples=min_samples)
dbscan.fit(X_std)

# Add cluster labels to your original data
data['Cluster'] = dbscan.labels_

# Visualize the clusters (for 2D data)
# Modify this part based on the number of features in your data
plt.scatter(X_std[:, 0], X_std[:, 1], c=dbscan.labels_, cmap='viridis', alpha=0.5)
plt.title('DBSCAN Clustering')
plt.xlabel('Max Voltage')
plt.ylabel('Max Current')
plt.show()
