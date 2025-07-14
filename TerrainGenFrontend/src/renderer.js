import * as THREE from 'three';
import { createNoise2D } from 'simplex-noise';
import { degToRad } from 'three/src/math/MathUtils.js';

const terrain_div = 81;
const terrain_scale = 0.02;


const _VS = `
varying vec4 v_Pos;

void main() {
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position.x, position.y * 3.0, position.z, 1.0);
    v_Pos = vec4(position, 1.0);
}
`;

const _FS = `

uniform float min_y;
uniform float max_y;
uniform float terrain_opacity;
// uniform vec3 cameraPosition;
varying vec4 v_Pos;

void main() {
    float relative_height = (v_Pos.y - min_y) / (max_y - min_y);
    // float relative_height = v_Pos.y;
    gl_FragColor = vec4(relative_height , relative_height , relative_height , 1.0);
    if(relative_height < 0.01)
    {
        gl_FragColor = vec4(0.3, 0.3, 1.0, 1.0);
    }
}
`;

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

function getIndices(div)
{
    var indices = [];
    for(let x = 0; x < div - 1; x++)
    {
        for(let y = 0; y < div - 1; y++)
        {
            let index = x + div * y;
            indices.push(index);
            indices.push(index + 1);
            indices.push(index + div);

            indices.push(index + div + 1);
            indices.push(index + div);
            indices.push(index + 1);
        }
    }
    return indices;
}

function getPositions(div, scale)
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

export function startRenderer(usePremadeValues, premadeValues)
{
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);


    const canvas = document.querySelector("#c");
    const renderer = new THREE.WebGLRenderer({canvas, alpha: true});

    const terrain_geometry = new THREE.BufferGeometry()

    let pos_ret = getPositions(terrain_div, terrain_div * terrain_scale);

    var positions, min_y, max_y
    if(usePremadeValues == true)
    {
        var new_positions = [];
        for(let x = 0; x < terrain_div; x++)
        {
            for(let y = 0; y < terrain_div; y++)
            {
                new_positions.push(x / (terrain_div * terrain_scale));
                new_positions.push(premadeValues[x + terrain_div*y]);
                new_positions.push(y / (terrain_div * terrain_scale));
            }
        }
        positions = new Float32Array(new_positions);
        min_y = Math.min.apply(this, premadeValues);
        max_y = Math.max.apply(this, premadeValues);  
    }
    else
    {
        positions = pos_ret[0];
        min_y = pos_ret[1];
        max_y = pos_ret[2];
    }
    terrain_geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    terrain_geometry.setIndex(getIndices(terrain_div));
    terrain_geometry.computeVertexNormals();

    const geometry = new THREE.BoxGeometry( 1, 1, 1 );
    const material = new THREE.ShaderMaterial({
        uniforms: {max_y : {value: max_y}, min_y : {value: min_y}, terrain_opacity : {value : 0.5}},
        vertexShader: _VS,
        fragmentShader: _FS
    });
    const cube = new THREE.Mesh( geometry, material );


    const terrain = new THREE.Mesh(terrain_geometry, material);
    terrain.geometry.center()
    scene.add(terrain);
    console.log(terrain.position);

    camera.position.z = 30;
    camera.position.y += 30;
    camera.rotateX(degToRad(-45));
    cube.position.setY(1.0);


    function animate() {
        cube.rotation.y += 0.01;
        material.uniforms.terrain_opacity.value = Math.min(material.uniforms.terrain_opacity.value + 0.01, 1.0);
    // terrain.rotation.x += 0.05;
    terrain.rotation.y += 0.005;
    // terrain.rotation.z += 0.05;
    // camera.rotation.y += 0.05;
    renderer.setSize(window.innerWidth, window.innerHeight);
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.render( scene, camera );
    }
    renderer.setAnimationLoop( animate );
}