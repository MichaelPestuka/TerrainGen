import { useEffect, useState, useRef } from 'react'
import './App.css'
import { Box, Button, Slider, TextField, Toolbar } from '@mui/material'
import './Sliders.jsx'
import InputSlider from './Sliders.jsx'
import TerrainRenderer from './renderer.js'


function App() {
  useEffect(() => {
      console.log("Starting Renderer");
      renderer.current = new TerrainRenderer({'Width':worldSize, 'Height':worldSize, 'Sealevel':seaLevel})
      // startRenderer();
  }, [])


  var renderer = useRef()
  const [worldSize, setSize] = useState(200)
  const [seaLevel, setSeaLevel] = useState(0.5)
  const [cliffMultiplier, setCliffMutliplier] = useState(6.0)
  const [cliffOffset, setCliffOffset] = useState(0.0)
  // const [seed, setSeed] = useState(Math.floor(Math.random() * 1000000))

  useEffect(() => {
    // console.log(cliffMultiplier);
    renderer.current.updateShaderValue("cliffMultiplier", cliffMultiplier);
    renderer.current.updateShaderValue("cliffOffset", cliffOffset);
  }, [cliffMultiplier, cliffOffset])
  return (
    <>
      <canvas id="c"></canvas>
      <Box sx={{width: '100vw', height: '100vh', zIndex: '1', position: 'absolute', top: '0', left: '0'}}>

      <Box sx={{marginRight: '0', width: '20%', height: '100%', padding: '0 32px'}}>
        <Button variant='contained' onClick={() => {
          console.log("Reloafdinf")
          renderer.current.fetchTerrain({'Width':worldSize, 'Height':worldSize, 'Sealevel':seaLevel});
        }}> Regenerate </Button>
        <TextField label={"Seed"} type='number'/>
        <InputSlider label={"World size"} updateData={setSize} minValue={50} maxValue={300} defaultValue={200} />
        <InputSlider label={"Sea level"} updateData={setSeaLevel} minValue={0.0} maxValue={1.0} defaultValue={0.5} />
        <InputSlider label={"Cliff Multiplier"} updateData={setCliffMutliplier} minValue={0.0} maxValue={10.0} defaultValue={6.0} />
        <InputSlider label={"Cliff offset"} updateData={setCliffOffset} minValue={-0.1} maxValue={0.1} defaultValue={0.0} />
      </Box>

      </Box>
    </>
  )
}

export default App
