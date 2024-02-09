import torch
import torch.nn as nn
from torch.utils.data import TensorDataset, DataLoader

class MyLSTM(nn.Module):
    _hidden_size: int
    _num_layers: int
    _sequence_length: int
    _target_length: int
    _criterion: nn.Module
    _optimizer: torch.optim.Optimizer
    _learning_rate = 0.01
    _num_epochs = 100
    _lstm: nn.LSTM
    _fc: nn.Linear

    def __init__(self, input_size=1, output_size=1, sequence_length=168, target_length=24, hidden_size=100, num_layers=1,
                 criterion=nn.MSELoss(), optimizer=torch.optim.Adam, learning_rate=0.01, 
                 num_epochs=100):
        super(MyLSTM, self).__init__()
        self._hidden_size = hidden_size
        self._num_layers = num_layers
        self._sequence_length = sequence_length
        self._target_length = target_length
        self._lstm = nn.LSTM(input_size, hidden_size, num_layers, batch_first=True)
        self._fc = nn.Linear(hidden_size, output_size)

        self._criterion = criterion
        self._optimizer = optimizer
        self._learning_rate = learning_rate
        self._num_epochs = num_epochs

    def forward(self, x):
        device = torch.device('cuda' if torch.cuda.is_available() else 'cpu') # Use GPU if available

        h0 = torch.zeros(self._num_layers, x.size(0), self._hidden_size).to(device)
        c0 = torch.zeros(self._num_layers, x.size(0), self._hidden_size).to(device)

        out, _ = self._lstm(x, (h0, c0))
        # out: tensor of shape (batch_size, seq_length, output_size)
        out = self._fc(out[:, -self._target_length:, :])
        return out

    def train_with_data(self, dataloader: DataLoader):
        # Initialize the model, loss function, and optimizer
        device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
        self._lstm = self._lstm.to(device)
        self._criterion = self._criterion.to(device)
        self._optimizer = self._optimizer(self._lstm.parameters(), lr=self._learning_rate)

        loss_list = []
        # Training the model
        # Loop over the dataset multiple times
        for epoch in range(self._num_epochs): 
            # Loop over each mini-batch
            for batch_sequences, batch_targets in dataloader:
                batch_sequences, batch_targets = batch_sequences.float().to(device), batch_targets.float().to(device)

                # Forward pass
                outputs = self(batch_sequences)

                # Flatten the target tensor to match the shape of predictions
                last_predictions = outputs.reshape(-1)

                # Calculate the loss
                loss = self._criterion(last_predictions.view(-1), batch_targets.view(-1))

                # Backward and optimize
                self._optimizer.zero_grad()
                loss.backward()
                self._optimizer.step()

            loss_list.append(loss.item())
            print(f'Epoch [{epoch+1}/{self._num_epochs}], Loss: {loss.item():.4f}')

        return loss_list

    def predict_next_n_days(self, input_sequences):
        # Predict the next n days
        print(f'input_sequences.shape: {input_sequences.shape}')
        print(f'input_sequences: {input_sequences}')
        if len(input_sequences) != self._sequence_length:
            raise ValueError(f'input_sequences must be of length {self._sequence_length}')

        input_sequences = input_sequences.reshape(1, -1, 1)
        print(f'reshaped input_sequences.shape: {input_sequences.shape}')
        print(f'reshaped input_sequences: {input_sequences}')

        # Make predictions
        with torch.no_grad():
            device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
            inputs = torch.tensor(input_sequences).float().to(device)

            output = self(inputs).cpu().numpy()

            print(f'output.shape: {output.shape}')
            print(f'output: {output}')

            return output


    def save_to(self, path):
        # Save the model to a file
        # Used to save the model and load from other nodes
        torch.save(self._lstm.state_dict(), path + '_model.pth')
    
    def load_model(self, path):
        # Load the model from a file
        self._lstm.load_state_dict(torch.load(path + '_model.pth'))