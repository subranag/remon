import logo from './logo.svg';
import './App.css';
import { Button, Grommet, Header, Menu } from 'grommet'
import { Trigger,  Cpu, Memory, Disc, Cluster, Search} from 'grommet-icons'

function App() {
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
    </Grommet>
  );
}

export default App;
