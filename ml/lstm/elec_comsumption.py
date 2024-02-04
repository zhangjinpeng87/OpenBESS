import pandas as pd
import matplotlib.pyplot as plt
from single_output_lstm import MyLSTM
import argparse

def load_data(file_path) -> pd.DataFrame:
    raw_data = pd.read_csv(file_path, sep=',', quotechar='"')
    raw_data['Local Time'] = pd.to_datetime(raw_data['Local Time'])
    raw_data['demand'] = pd.to_numeric(raw_data['Demand'].str.replace(',', ''), errors='coerce')

    return raw_data

def main() -> int:
    # Get parameters from the command line
    parser = argparse.ArgumentParser(
        prog='electric_demand_predict',
        description='Train a multi-input LSTM model to predict demand')
    parser.add_argument('--num_epochs', type=int, default=10, help='Number of epochs')
    parser.add_argument('--learning_rate', type=float, default=0.001, help='Learning rate')
    parser.add_argument('--batch_size', type=int, default=64, help='Batch size')
    parser.add_argument('--hidden_size', type=int, default=100, help='Hidden size')
    parser.add_argument('--num_layers', type=int, default=1, help='Number of layers')
    parser.add_argument('--sequence_length', type=int, default=24, help='Sequence length')
    parser.add_argument('--train_dataset', type=str, default='./train.csv', help='Path to the training dataset')
    parser.add_argument('--test_dataset', type=str, default='./test.csv', help='Path to the test dataset')
    parser.add_argument('--save_to', type=str, default='model', help='Path to save the model')
    args = parser.parse_args()

    # Load the data
    data = load_data(args.train_dataset)

    print('data loaded, start to train the model')
    print(data.head())

    # Train the model
    lstm = MyLSTM(hidden_size=args.hidden_size, 
                          num_layers=args.num_layers, save_to=args.save_to, learning_rate=args.learning_rate, 
                          num_epochs=args.num_epochs, batch_size=args.batch_size, sequence_length=args.sequence_length)
    lstm.train()
    lstm.train_with_data(data['demand'])

    # Predict the demand
    data = load_data(args.test_dataset)

    lstm.eval()
    predictions = lstm.predict(data['demand'])

    # Draw the graph
    plt.figure(figsize=(12, 6))
    plt.plot(data['Local Time'][args.sequence_length:], predictions, label='Predicted Consumption', linestyle='--')
    plt.plot(data['Local Time'][args.sequence_length:], data['demand'][args.sequence_length:], label='Real Consumption', linestyle='-')
    plt.xlabel('Time')
    plt.ylabel('Consumption')
    plt.title('Consumption Prediction vs Real Consumption')
    plt.legend()
    plt.show()

    return 0

if __name__ == '__main__':
    main()