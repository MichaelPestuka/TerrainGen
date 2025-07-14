import { useEffect } from 'react'
import './App.css'
import './renderer.js'
// import { startRenderer } from './renderer.js'
import { getNewTerrain } from './requests.js'

function App() {
  useEffect(() => {
      console.log("Starting Renderer");
      getNewTerrain();
      // startRenderer();
  })

  return (
    <>
      <canvas id="c"></canvas>
    </>
  )
}

export default App
