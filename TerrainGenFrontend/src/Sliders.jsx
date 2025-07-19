import { Slider, Box, Typography, Grid, Input } from '@mui/material'
import { useState } from 'react'

export default function InputSlider({ label, updateData, minValue, maxValue, defaultValue }) {
  const [value, setValue] = useState(defaultValue);

  const handleSliderChange = (event, newValue) => {
    setValue(newValue);
    updateData(newValue);
  };

  const handleInputChange = (event) => {
    setValue(event.target.value === '' ? minValue : Number(event.target.value));
    updateData(event.target.value);
  };

  const handleBlur = () => {
    if (value < minValue) {
      setValue(minValue);
    } else if (value > maxValue) {
      setValue(maxValue);
    }
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Grid container spacing={2} sx={{ alignItems: 'center' }}>
        <Grid>
            <Typography id="input-slider" gutterBottom>
                {label}
            </Typography>
        </Grid>
        <Grid size="grow">
          <Slider
            value={typeof value === 'number' ? value : minValue}
            onChange={handleSliderChange}
            aria-labelledby="input-slider"
            min={minValue}
            max={maxValue}
          />
        </Grid>
        <Grid>
          <Input
            value={value}
            size="small"
            onChange={handleInputChange}
            onBlur={handleBlur}
            inputProps={{
              step: 10,
              min: minValue,
              max: maxValue,
              type: 'number',
              'aria-labelledby': 'input-slider',
            }}
          />
        </Grid>
      </Grid>
    </Box>
  );
}