import './App.css';
import { Box, Button, Grommet, Header, Menu } from 'grommet'
import { Trigger, Cpu, Memory, Disc, Cluster, Search } from 'grommet-icons'
import GaugeChart from 'react-gauge-chart'
import { useEffect, useState } from 'react';

function App() {
  const [percent, setPercent] = useState(0.86);

  const randomBetween = (min, max) => {
    return min + (max - min) * Math.random();
  }

  useEffect(() => {
    const timer = setInterval(() => {
      setPercent(randomBetween(0.5, 0.7));
    }, 1000)

    return () => clearInterval(timer);
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
      <Box pad="small">
        <GaugeChart id="gauge-chart2"
          nrOfLevels={20}
          percent={percent}
          textColor="black"
          animDelay={100}
          animateDuration={500}
        />
      </Box>
    </Grommet>
  );
}

export default App;
