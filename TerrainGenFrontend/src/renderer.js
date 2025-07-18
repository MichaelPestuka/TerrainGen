import * as THREE from 'three';
import { createNoise2D } from 'simplex-noise';
import { degToRad } from 'three/src/math/MathUtils.js';

const terrain_scale = 0.02;


const _VS = `
varying float v_Pos;
varying vec2 texCoord;
uniform float width;
uniform float height; 
void main() {
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position.x, position.y * 5.0, position.z, 1.0);
    v_Pos = (position.y + 1.0) / 2.0;
    texCoord = vec2(position.x / (width/4.0) + 0.5, position.z / (height/4.0) + 0.5);
}
`;

const _FS = `
uniform sampler2D texture1;
uniform float min_y;
uniform float max_y;
uniform float terrain_opacity;
varying float v_Pos;
varying vec2 texCoord;

void main() {
    // float relative_height = (v_Pos.y - min_y) / (max_y);
    float relative_height = v_Pos;
    gl_FragColor = vec4(relative_height , relative_height, relative_height, 1.0);
    if(relative_height < 0.5 && relative_height >= 0.0)
    {
        gl_FragColor = vec4(0.2, 0.2, 0.5, 1.0);
    }
    if(texture(texture1, texCoord).r < 0.1 && texture(texture1, texCoord).b < 0.1 ) {
    
        gl_FragColor = vec4(relative_height, relative_height, relative_height, 1.0);
    }
    else {
        gl_FragColor = vec4(texture(texture1, texCoord).rgb, 1.0);
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

export function startRenderer(usePremadeValues, premadeValues, width, height, textureURL)
{
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);


    const canvas = document.querySelector("#c");
    const renderer = new THREE.WebGLRenderer({canvas, alpha: true});

    const terrain_geometry = new THREE.BufferGeometry()


    var positions, min_y, max_y
    if(usePremadeValues == true)
    {
        var new_positions = [];
        for(let x = 0; x < width; x++)
        {
            for(let y = 0; y < height; y++)
            {
                new_positions.push(x / (width * terrain_scale));
                new_positions.push(premadeValues[x + width*y]);
                new_positions.push(y / (height * terrain_scale));
            }
        }
        positions = new Float32Array(new_positions);
        min_y = Math.min.apply(this, premadeValues);
        max_y = Math.max.apply(this, premadeValues);
    }
    else
    {
        let pos_ret = getPositions(width, width * terrain_scale);
        positions = pos_ret[0];
        min_y = pos_ret[1];
        max_y = pos_ret[2];
    }
    terrain_geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    terrain_geometry.setIndex(getIndices(width));
    terrain_geometry.computeVertexNormals();


    var textureLoader = new THREE.TextureLoader();
    textureLoader.crossOrigin = "anonymous";
    var texture = textureLoader.load(textureURL);
    texture.wrapS = THREE.RepeatWrapping;
    texture.wrapT = THREE.RepeatWrapping;
    texture.minFilter = THREE.NearestFilter;
    texture.magFilter = THREE.NearestFilter;
    texture.flipY = false;
    console.log(texture);

    const material = new THREE.ShaderMaterial({
        uniforms: {max_y : {value: max_y}, min_y : {value: min_y}, width : {value : width}, height : {value : height}, terrain_opacity : {value : 0.5}, texture1 : {value : texture}},
        vertexShader: _VS,
        fragmentShader: _FS,
    });


    const terrain = new THREE.Mesh(terrain_geometry, material);
    terrain.geometry.center()
    scene.add(terrain);
    console.log(terrain.position);

    camera.position.z = 30;
    camera.position.y += 30;
    camera.rotateX(degToRad(-45));


    function animate() {
        material.uniforms.terrain_opacity.value = Math.min(material.uniforms.terrain_opacity.value + 0.01, 1.0);
        terrain.rotation.y += 0.005;
        renderer.setSize(window.innerWidth, window.innerHeight);
        camera.aspect = window.innerWidth / window.innerHeight;
        camera.updateProjectionMatrix();
        renderer.render( scene, camera );
    }
    renderer.setAnimationLoop( animate );
}