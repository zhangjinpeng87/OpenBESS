import pandas as pd
from sklearn.cluster import KMeans
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

# Choose the number of clusters
num_clusters = 5

# Perform KMeans clustering
kmeans = KMeans(n_clusters=num_clusters, random_state=42)
kmeans.fit(X_std)

# Add cluster labels to your original data
data['Cluster'] = kmeans.labels_

# Visualize the clusters (for 2D data)
# Modify this part based on the number of features in your data
plt.scatter(X_std[:, 0], X_std[:, 1], c=kmeans.labels_, cmap='viridis', alpha=0.5)
plt.scatter(kmeans.cluster_centers_[:, 0], kmeans.cluster_centers_[:, 1], marker='x', s=300, c='red')
plt.title('KMeans Clustering')
plt.xlabel('Feature 1')
plt.ylabel('Feature 2')
plt.show()
