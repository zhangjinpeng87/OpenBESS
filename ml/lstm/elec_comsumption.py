import torch
import pandas as pd
import numpy as np
from sklearn.preprocessing import MinMaxScaler
import matplotlib.pyplot as plt
from my_lstm import MyLSTM
import argparse

from torch.utils.data import DataLoader, Dataset

class DemandDataset(Dataset):
    def __init__(self, data, sequence_length, target_length):
        self.data = data
        self.sequence_length = sequence_length
        self.target_length = target_length

    def __len__(self):
        return len(self.data) - (self.sequence_length + self.target_length) + 1

    def __getitem__(self, idx):
        idx_end = idx + self.sequence_length
        input_sequence = self.data[idx:idx_end]

        idx_target = idx_end + self.target_length - 1
        target_sequence = self.data[idx_end:idx_target+1]

        return torch.tensor(input_sequence).view(-1, 1), torch.tensor(target_sequence).view(-1, 1)

def load_data(file_path) -> pd.DataFrame:
    raw_data = pd.read_csv(file_path, sep=',', quotechar='"')
    raw_data['Local Time'] = pd.to_datetime(raw_data['Local Time'])
    raw_data['demand'] = pd.to_numeric(raw_data['Demand'].str.replace(',', ''), errors='coerce')

    return raw_data

def main() -> int:
    # Get parameters from the command line
    parser = argparse.ArgumentParser(
        prog='electric_demand_predict',
        description='Train a LSTM model to predict electricity demand')
    parser.add_argument('--num_epochs', type=int, default=10, help='Number of epochs')
    parser.add_argument('--learning_rate', type=float, default=0.01, help='Learning rate')
    parser.add_argument('--batch_size', type=int, default=64, help='Batch size')
    parser.add_argument('--hidden_size', type=int, default=100, help='Hidden size')
    parser.add_argument('--num_layers', type=int, default=1, help='Number of layers')
    parser.add_argument('--sequence_length', type=int, default=168, help='Sequence length means how many hours to look back in the past, 168 means look back 7 days.')
    parser.add_argument('--target_length', type=int, default=24, help='Target length means how many hours to predict in the future, 24 means predict next 1 days.')
    parser.add_argument('--train_dataset', type=str, default='./train.csv', help='Path to the training dataset')
    parser.add_argument('--verify_dataset', type=str, default='./test.csv', help='Path to the verify dataset')
    parser.add_argument('--num_threads', type=int, default=4, help='Number of threads to use when training the model')
    parser.add_argument('--save_to', type=str, default='model', help='Path to save the model')
    args = parser.parse_args()

    # Load train dataset
    train_raw_data = load_data(args.train_dataset)
    train_data = train_raw_data['demand'].values

    print(f'train_data.shape: {train_data.shape}')
    print(f'train_data: {train_data}')

    demand_scaler = MinMaxScaler(feature_range=(0, 1))
    train_data = demand_scaler.fit_transform(train_data.reshape(-1, 1))
    train_dataloader = DataLoader(DemandDataset(train_data, args.sequence_length, args.target_length), 
                                  batch_size=args.batch_size, shuffle=True, num_workers=args.num_threads)

    # Train the model
    lstm = MyLSTM(hidden_size=args.hidden_size, sequence_length=args.sequence_length, target_length=args.target_length, 
                          num_layers=args.num_layers, save_to=args.save_to, learning_rate=args.learning_rate, 
                          num_epochs=args.num_epochs)
    lstm.train()
    loss_list = lstm.train_with_data(train_dataloader)
    plt.figure(figsize=(12, 6))
    plt.plot(loss_list)
    plt.xlabel('Epoch')
    plt.ylabel('Loss')
    plt.title('Training Loss')
    plt.show()

    # Predict the demand
    # Use last 30 days of the training data to predict the next 7 days
    input_sequences = train_data[-args.sequence_length:]
    lstm.eval()
    predictions = lstm.predict_next_n_days(input_sequences)
    predictions = demand_scaler.inverse_transform(predictions.reshape(-1, 1))

    # Load real demand data
    real_demand = load_data(args.verify_dataset)

    real_time = train_raw_data['Local Time'].values[-args.sequence_length:]
    real_time = np.append(real_time, real_demand['Local Time'].values[:args.target_length])
    real_data = train_raw_data['demand'].values[-args.sequence_length:]
    real_data = np.append(real_data, real_demand['demand'].values[:args.target_length])


    # Draw the graph
    plt.figure(figsize=(12, 6))
    plt.plot(real_time[-args.target_length:], predictions, label='Predicted Consumption', linestyle='--')
    plt.plot(real_time[:], real_data, label='Real Consumption', linestyle='-')
    plt.xlabel('Time')
    plt.ylabel('Consumption')
    plt.title('Consumption Prediction vs Real Consumption')
    plt.legend()
    plt.show()

    return 0

if __name__ == '__main__':
    main()