import { useEffect, useState } from 'react'
import './App.css'
import './renderer.js'
import { Box, Button, Slider, Toolbar } from '@mui/material'
import './Sliders.jsx'
// import { startRenderer } from './renderer.js'
import { getNewTerrain } from './requests.js'
import InputSlider from './Sliders.jsx'

function App() {
  useEffect(() => {
      console.log("Starting Renderer");
      getNewTerrain();
      // startRenderer();
  }, [])

  const [width, setWidth] = useState(200)
  const [height, setHeight] = useState(200)

  const updateWidth = (data) => {
    setWidth(data);
  };

  const updateHeight = (data) => {
    setHeight(data);
  };
  return (
    <>
      <canvas id="c"></canvas>
      <Box sx={{width: '100vw', height: '100vh', zIndex: '1'}}>

      <Box sx={{marginRight: '0', width: '20%', height: '100%'}}>
        <Button variant='contained' onClick={() => {
          console.log("Reloafdinf")
          getNewTerrain();
        }}> Regenerate </Button>
        <InputSlider label={"World width"} updateData={updateWidth} minValue={100} maxValue={300} defaultValue={200} />
        <InputSlider label={"World height"} updateData={updateHeight} minValue={100} maxValue={300} defaultValue={200}/>
      </Box>

      </Box>
    </>
  )
}

export default App
