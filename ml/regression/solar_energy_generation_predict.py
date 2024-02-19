import numpy as np
import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error, r2_score
import matplotlib.pyplot as plt
import argparse



# We are going to predict the solar energy generation based on the weather data.
# Typically, the solar energy generation is influenced by the following factors:
# - Temperature
# - UV index
# - Sunshine duration
# - Weather description
# 
# But these data is not easy to get, we will use the following features to predict
# the solar energy generation:
# - Month
# - Season



def main() -> int:
    parser = argparse.ArgumentParser(
        prog='solar_energy_generation_predict',
        description='Train a linear regression model to predict electricity demand')
    
    
    # Load data from a CSV file


    return 0





if __name__ == '__main__':
    main()