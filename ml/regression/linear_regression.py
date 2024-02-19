import numpy as np
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error, r2_score

class MyLinearRegression:
    def __init__(self, data: pd.DataFrame, features: list, target: float):
        self.model = LinearRegression()
        self.data = data
        self.features = features
        self.target = target

    def train(self, X_train, y_train):
        self.model.fit(X_train, y_train)

    def predict(self, X_test):
        return self.model.predict(X_test)