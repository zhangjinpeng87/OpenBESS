from sklearn.svm import SVR


svr = SVR(kernel='rbf', C=100, gamma=0.1, epsilon=.1)