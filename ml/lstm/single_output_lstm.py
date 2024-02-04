import torch
import torch.nn as nn
from torch.utils.data import TensorDataset, DataLoader
from sklearn.preprocessing import MinMaxScaler

class MyLSTM(nn.Module):
    def __init__(self, input_size=1, output_size=1, hidden_size=100, num_layers=1, save_to='',
                 criterion=nn.MSELoss(), optimizer=torch.optim.Adam, learning_rate=0.01, 
                 num_epochs=100, batch_size=64, sequence_length=10, shuffle=True):
        super(MyLSTM, self).__init__()
        self.hidden_size = hidden_size
        self.num_layers = num_layers
        self.lstm = nn.LSTM(input_size, hidden_size, num_layers, batch_first=True)
        self.fc = nn.Linear(hidden_size, output_size)

        self.save_to = save_to
        self.criterion = criterion
        self.optimizer = optimizer
        self.learning_rate = learning_rate
        self.num_epochs = num_epochs
        self.batch_size = batch_size
        self.sequence_length = sequence_length
        self.shuffle = shuffle
        self.scaler = MinMaxScaler()

    def forward(self, x):
        device = torch.device('cuda' if torch.cuda.is_available() else 'cpu') # Use GPU if available

        h0 = torch.zeros(self.num_layers, x.size(0), self.hidden_size).to(device)
        c0 = torch.zeros(self.num_layers, x.size(0), self.hidden_size).to(device)

        out, _ = self.lstm(x, (h0, c0))
        out = self.fc(out[:, -1, :])
        return out

    def train_with_data(self, train_data):
        # Normalize the data
        train_data = self.scaler.fit_transform(train_data.values.reshape(-1, 1))

        # Initialize the model, loss function, and optimizer
        device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        self.lstm = self.lstm.to(device)
        self.criterion = self.criterion.to(device)
        self.optimizer = self.optimizer(self.lstm.parameters(), lr=self.learning_rate)

        # Create sequences and targets
        sequences, targets = create_training_used_sequences(train_data, self.sequence_length)
        dataset = TensorDataset(sequences, targets)
        dataloader = DataLoader(dataset, batch_size=self.batch_size, shuffle=self.shuffle)

        # Training the model
        # Loop over the dataset multiple times
        for epoch in range(self.num_epochs): 
            # Loop over each mini-batch
            for batch_sequences, batch_targets in dataloader:
                batch_sequences, batch_targets = batch_sequences.float().to(device), batch_targets.float().to(device)
                outputs = self(batch_sequences)
                # calculate the loss
                loss = self.criterion(outputs, batch_targets)
                self.optimizer.zero_grad()
                loss.backward()
                self.optimizer.step()

            # if (epoch+1) % 10 == 0:
            print(f'Epoch [{epoch+1}/{self.num_epochs}], Loss: {loss.item():.4f}')

    def predict(self, data):
        # Feature Scaling
        data = self.scaler.transform(data.values.reshape(-1, 1))

        # Create sequences for LSTM
        input_sequences = create_predict_used_sequences(data, self.sequence_length)

        # Make predictions
        with torch.no_grad():
            device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
            inputs = input_sequences.float().to(device)
            predictions = self(inputs).cpu().numpy()

            print('predictions shape:', predictions.shape)
            print('predictions:', predictions)

            # Inverse the scaling
            # for col in self.output_features:
            predictions = self.scaler.inverse_transform(predictions)
            data = self.scaler.inverse_transform(data)
            return predictions

    def save_to(self):
        # Save the model to a file
        # Used to save the model and load from other nodes
        torch.save(self.lstm.state_dict(), self.save_to + '_model.pth')

        # Save the scalers to a file
        torch.save(self.scaler, self.save_to + '_scaler.pth')
    
    def load_model(self):
        # Load the model from a file
        self.lstm = MultiInputLSTM(len(self.input_features), self.hidden_size, self.num_layers, len(self.output_features))
        self.lstm.load_state_dict(torch.load(self.save_to + '_model.pth'))

        # Load the scalers
        self.scaler = torch.load(self.save_to + '_scaler.pth')

def create_training_used_sequences(data, sequence_length):
        sequences, targets = [], []
        for i in range(len(data) - sequence_length):
            seq = data[i:i+sequence_length-1]
            label = data[i+sequence_length]
            sequences.append(seq)
            targets.append(label)
        return torch.tensor(sequences), torch.tensor(targets)

def create_predict_used_sequences(data, sequence_length):
        sequences = []
        for i in range(len(data) - sequence_length):
            seq = data[i:i+sequence_length-1]
            sequences.append(seq)
        return torch.tensor(sequences)