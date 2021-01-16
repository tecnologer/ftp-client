import "./App.css";
import FilesContainer from "./components/filesContainer";
import InputCn from "./components/inputConnection";
import bodyParser from "body-parser";

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>FTP Client</h1>
      </header>
      <div>
        <InputCn />
      </div>
      <div>
        <FilesContainer />
      </div>
    </div>
  );
}

export default App;
