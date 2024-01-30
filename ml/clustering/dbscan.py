'''
The DBSCAN algorithm views clusters as areas of high density separated by areas of low density. 
Due to this rather generic view, clusters found by DBSCAN can be any shape, as opposed to k-means
which assumes that clusters are convex shaped. The central component to the DBSCAN is the concept
of core samples, which are samples that are in areas of high density. A cluster is therefore a set
of core samples, each close to each other (measured by some distance measure) and a set of non-core
samples that are close to a core sample (but are not themselves core samples). There are two parameters
to the algorithm, min_samples and eps, which define formally what we mean when we say dense. Higher
min_samples or lower eps indicate higher density necessary to form a cluster.

More formally, we define a core sample as being a sample in the dataset such that there exist min_samples
other samples within a distance of eps, which are defined as neighbors of the core sample. This tells us
that the core sample is in a dense area of the vector space. A cluster is a set of core samples that can 
be built by recursively taking a core sample, finding all of its neighbors that are core samples, finding 
all of their neighbors that are core samples, and so on. A cluster also has a set of non-core samples, 
which are samples that are neighbors of a core sample in the cluster but are not themselves core samples. 
Intuitively, these samples are on the fringes of a cluster.

Any core sample is part of a cluster, by definition. Any sample that is not a core sample, and is at least
eps in distance from any core sample, is considered an outlier by the algorithm.

While the parameter min_samples primarily controls how tolerant the algorithm is towards noise (on noisy
and large data sets it may be desirable to increase this parameter), the parameter eps is crucial to choose
appropriately for the data set and distance function and usually cannot be left at the default value. It 
controls the local neighborhood of the points. When chosen too small, most data will not be clustered at all
(and labeled as -1 for “noise”). When chosen too large, it causes close clusters to be merged into one cluster,
and eventually the entire data set to be returned as a single cluster. Some heuristics for choosing this 
parameter have been discussed in the literature, for example based on a knee in the nearest neighbor distances
plot.
'''

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
