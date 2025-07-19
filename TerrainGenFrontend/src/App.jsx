import { useEffect, useState, useRef } from 'react'
import './App.css'
import { Box, Button, Slider, Toolbar } from '@mui/material'
import './Sliders.jsx'
import InputSlider from './Sliders.jsx'
import TerrainRenderer from './renderer.js'


function App() {
  useEffect(() => {
      console.log("Starting Renderer");
      renderer.current = new TerrainRenderer(worldSize, worldSize)
      // startRenderer();
  }, [])

  var renderer = useRef()
  const [worldSize, setSize] = useState(200)

  const updateSize = (data) => {
    setSize(data);
  };

  return (
    <>
      <canvas id="c"></canvas>
      <Box sx={{width: '100vw', height: '100vh', zIndex: '1', position: 'absolute', top: '0', left: '0'}}>

      <Box sx={{marginRight: '0', width: '20%', height: '100%', padding: '0 32px'}}>
        <Button variant='contained' onClick={() => {
          console.log("Reloafdinf")
          renderer.current.fetchTerrain(worldSize, worldSize);
        }}> Regenerate </Button>
        <InputSlider label={"World size"} updateData={updateSize} minValue={50} maxValue={300} defaultValue={200} />
      </Box>

      </Box>
    </>
  )
}

export default App
