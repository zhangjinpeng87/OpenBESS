package simulation




// Sensor is a sensor.
type MockSensor struct {
	Station int `json:"station"`
	Container int `json:"container"`
	Pack int `json:"pack"`
	Cell int `json:"cell"`
)



func NewMockSensor(station, container, pack, cell int) *MockSensor {
	return &Sensor{
		Station: station,
		Container: container,
		Pack: pack,
		Cell: cell,
	}
}

func (s *Sensor) Start(ctx context, ) error {
	errg, c := errgroup.WithContext(ctx)
	errg.Go(func() error {
		timer := time.NewTimer(1 * time.Second)
		for {

		}
}