/**
 * @author Michael Pestuka
 * Helper rendering functions
 */

import { createNoise2D } from 'simplex-noise';

function fractalNoise(x, y, noise, octaves, frequency, amplitude)
{
    var result = 0.0;
    for (let i = 0; i < octaves; i++)
    {
        result += amplitude * noise(x * frequency, y * frequency);
        frequency *= 2.0;
        amplitude *= 0.5;
    }
    return result
}

 
export function getPositions(div, scale)
{
    const noise = createNoise2D();
    var positions = [];
    let min_y = 1000;
    let max_y = -1000;
    for(let x = 0; x < div; x++)
    {
        for(let y = 0; y < div; y++)
        {
            positions.push(x / scale);
            let new_y = fractalNoise(x, y, noise, 4, 0.015, 1.0);
            min_y = Math.min(min_y, new_y);
            max_y = Math.max(max_y, new_y);
            positions.push(new_y);
            positions.push(y / scale);
        }
    }
    return [new Float32Array(positions), min_y, max_y];
}