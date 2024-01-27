# OpenBMS
OpenBMS is a open source battery management system(BMS), aim to provide BMS for battery energy storage systems. OpenBMS monitor SOC and SOH of each battery cell in real-time, automatically balance the charge and discharge operations to achieve longer lifespan and higher efficiency for batteries. OpenBMS also collect batterier data and upload them to cloud, periodically analyze and learn these latest data and historical data to predict the remain lifetime of these batteries and provide maintenance suggestions in advance. 

![image](./resources/bms-monitor.png)

## Domain Knowleges
- [x] [Characteristics](./docs/characteristic-of-lithium-ion-battery.md) of Lithium-ion batteries
- [x] [Thermal runaway issue](./docs/thermal-runaway-issue.md) of Li-ion batteries
- [x] [Why low temperature decrease batteries' performance?](./docs/decreased-perf-under-low-temperature.md)
- [x] [How to maximize Lithium-ion battery life?](./docs/how-to-maximize-lithium-ion-battery-life.md)

## Data Store
- [ ] Implement robust data management for monitoring and collecting data from individual batteries.
- [ ] Store time series data locally.
- [ ] Upload data to cloud storage like S3.

## Intelligent and Predictive Maintenance
- [ ] Process latest collected data and historical data with ML and BigData algorithms, to predict SOH(State-Of-Health) of each battery pack, and provide maintenance suggestions in advance.
- [ ] Analyze operational data to identify opportunities for improving overall energy efficiency, suggesting adjustments in the system configuration or operation.

## Battery Managment

### State of Charge Calculation
- [ ] Develop algorithms for accurate SoC calculations based on battery voltage, current, and temperature measurements.
- [ ] Consider incorporating advanced techniques like Kalman filters for better estimation.

### Balancing and Equalization:
- [ ] Implement balancing algorithms to ensure that individual cells within a battery pack are charged and discharged uniformly.
- [ ] Consider equalization strategies to prolong the overall battery life.

### Energy Management
- [ ] Optimize energy usage based on grid demand and supply.
- [ ] Implement algorithms for determining when to charge and discharge batteries to maximize efficiency.

### Thermal Management:
- [ ] Monitor and control battery temperatures to optimize performance and prevent overheating.
- [ ] Implement thermal management strategies to ensure safe operation.

### State of Health Monitoring
- [ ] Monitor state of health for each pack and each cell
- [ ] Show SOH on dashboard 
