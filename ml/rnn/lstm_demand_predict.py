import torch
import torch.nn as nn
import numpy as np
import pandas as pd
from sklearn.preprocessing import MinMaxScaler
from torch.utils.data import DataLoader, TensorDataset

# Sample Data Loading: format is "date, weather, demand"
data = pd.read_csv('dataset.csv')
data['date'] = pd.to_datetime(data['date'])
data = data.sort_values('date')
data = data.set_index('date')

# Feature Scaling
scaler = MinMaxScaler()
data_scaled = scaler.fit_transform(data['demand'].values.reshape(-1, 1))

# Create sequences for LSTM
def create_sequences(data, sequence_length):
    sequences = []
    targets = []
    for i in range(len(data) - sequence_length):
        seq = data[i:i+sequence_length]
        label = data[i+sequence_length:i+sequence_length+1]
        sequences.append(seq)
        targets.append(label)
    return torch.tensor(sequences), torch.tensor(targets)

# Hyperparameters
input_size = 1
hidden_size = 100
num_layers = 1
output_size = 1
num_epochs = 100
learning_rate = 0.01
sequence_length = 10  # Adjust as needed

# Create sequences and targets
sequences, targets = create_sequences(data_scaled, sequence_length)

# Create DataLoader
dataset = TensorDataset(sequences, targets)
dataloader = DataLoader(dataset, batch_size=64, shuffle=True)

# LSTM Model
class LSTM(nn.Module):
    def __init__(self, input_size, hidden_size, num_layers, output_size):
        super(LSTM, self).__init__()
        self.hidden_size = hidden_size
        self.num_layers = num_layers
        self.lstm = nn.LSTM(input_size, hidden_size, num_layers, batch_first=True)
        self.fc = nn.Linear(hidden_size, output_size)

    def forward(self, x):
        h0 = torch.zeros(self.num_layers, x.size(0), self.hidden_size).to(device)
        c0 = torch.zeros(self.num_layers, x.size(0), self.hidden_size).to(device)

        out, _ = self.lstm(x, (h0, c0))
        out = self.fc(out[:, -1, :])
        return out

# Initialize the model, loss function, and optimizer
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
model = LSTM(input_size, hidden_size, num_layers, output_size).to(device)
criterion = nn.MSELoss()
optimizer = torch.optim.Adam(model.parameters(), lr=learning_rate)

# Training the model
for epoch in range(num_epochs):
    for batch_sequences, batch_targets in dataloader:
        batch_sequences, batch_targets = batch_sequences.float().to(device), batch_targets.float().to(device)

        # Forward pass
        outputs = model(batch_sequences)
        loss = criterion(outputs, batch_targets)

        # Backward and optimize
        optimizer.zero_grad()
        loss.backward()
        optimizer.step()

    if (epoch+1) % 10 == 0:
        print(f'Epoch [{epoch+1}/{num_epochs}], Loss: {loss.item():.4f}')

## Save only the model parameters (state dictionary)
# torch.save(model.state_dict(), 'lstm_demand_model.pth')

## Load the model from other nodes
# model = LSTM(input_size, hidden_size, num_layers, output_size).to(device)
# model.load_state_dict(torch.load('lstm_demand_model.pth'))

# Prediction
model.eval()
with torch.no_grad():
    test_sequence = data_scaled[-sequence_length:].reshape(1, sequence_length, input_size)
    test_sequence = torch.tensor(test_sequence).float().to(device)
    predicted_demand = model(test_sequence).item()

# Inverse transform the predicted value
predicted_demand = scaler.inverse_transform(np.array([[predicted_demand]]))

print(f'Predicted Electricity Demand for the next day: {predicted_demand[0][0]:.2f}')
