import axios from "axios";

import './App.css';

function App() {
  (async function() {
    const foo = await axios.get("//localhost:8000/dragons", {	headers: {
      'Access-Control-Allow-Origin': '*',
    },});
    console.log(foo);
  })()

  return (
    <div className="App">
      <header className="header">This is my web app</header>
    </div>
  );
}

export default App;
