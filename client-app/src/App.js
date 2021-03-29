import './App.css';
import { Box, Button, Grid, Grommet, Header, Menu } from 'grommet'
import { Trigger, Cpu, Memory, Disc, Cluster, Search } from 'grommet-icons'
import GaugeChart from 'react-gauge-chart'
import { useEffect, useRef, useState } from 'react';
import { Bar } from 'react-chartjs-2'

const client = new WebSocket(`ws://${window.location.host}/api/ws/cpu`)

const cpuBarChartOptions = {
  scales: {
    yAxes: [
      {
        ticks: {
          beginAtZero: true,
        },
      },
    ],
    xAxes: [],
  },
}

const cpuBarChartData = {
  labels: [],
  datasets: [
    {
      label: 'Per core utilization',
      data: [],
      backgroundColor: 'rgb(255, 99, 132)',
    }
  ],
}

function App() {
  const [percent, setPercent] = useState(0.0);

  // reference to the bar chart
  const barChartRef = useRef(null);
  const chartLables = [];
  const chartValues = [];

  useEffect(() => {
    client.onopen = () => {
      console.log('websocket client connected');
    };

    client.onmessage = (message) => {
      const data = JSON.parse(message.data);
      const cpuAgg = (data['cpu'] / 100.00);
      setPercent(cpuAgg);
      delete data['cpu'];

      var labels = barChartRef.current.props.data.labels;
      var chartData = barChartRef.current.props.data.datasets[0].data;
      
      if (labels.length === 0) {
        for (var key in data) {
          labels.push(key);
        }
        labels.sort();
      }
      
      if (chartData.length === 0) {
        for (var index in labels) {
          chartData.push(data[labels[index]]);
        }
      } else {
        var counter = 0;
        for (var index in labels) {
          chartData[counter++] = data[labels[index]]
        }
      }
      
      barChartRef.current.chartInstance.update();
    };
  }, []);

  return (
    <Grommet>

      <Header background="brand" align="center" justify="start" pad="small">
        <Button icon={<Cpu />} label="CPU" hoverIndicator />
        <Button icon={<Memory />} label="Memory" hoverIndicator />
        <Button icon={<Disc />} label="Disk" hoverIndicator />
        <Button icon={<Cluster />} label="Network" hoverIndicator />
        <Button icon={<Trigger />} label="Processes" hoverIndicator />
        <Button icon={<Search />} label="Logs" hoverIndicator />
      </Header>
      <Grid
        fill
        rows={['flex', 'flex']}
        columns={['flex', 'flex']}
        gap="small"
        areas={[
          { name: 'one', start: [0, 0], end: [1, 0] },
          { name: 'two', start: [1, 0], end: [1, 0] },
          { name: 'three', start: [0, 1], end: [0, 1] },
          { name: 'four', start: [1, 1], end: [1, 1] },
        ]}
      >
        <Box gridArea="one">
          <GaugeChart id="gauge-chart2"
            marginInPercent="0.02"
            nrOfLevels={20}
            percent={percent}
            textColor="black"
            animDelay={100}
            animateDuration={500}
          />
        </Box>
        <Box gridArea="two">
          <Bar data={cpuBarChartData} options={cpuBarChartOptions} ref={barChartRef} />
        </Box>
        <Box gridArea="three">

        </Box>
        <Box gridArea="four">

        </Box>
      </Grid>
    </Grommet>
  );
}

export default App;
