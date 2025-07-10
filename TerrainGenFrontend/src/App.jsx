import { useEffect } from 'react'
import './App.css'
import './renderer.js'
import { startRenderer } from './renderer.js'

function App() {
  useEffect(() => {
      console.log("Starting Renderer");
      startRenderer();
  })

  return (
    <>
      <canvas id="c"></canvas>
    </>
  )
}

export default App
