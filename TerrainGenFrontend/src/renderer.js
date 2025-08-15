import * as THREE from 'three';
import { degToRad } from 'three/src/math/MathUtils.js';
import TerrainShader from './terrainshader';
import { FlyControls } from 'three/addons/controls/FlyControls.js';
// const terrain_scale = 0.02;

function getIndices(width, height)
{
    var indices = [];
    for(let x = 0; x < width - 1 ; x++)
    {
        for(let y = 0; y < height - 1; y++)
        {
            let index = y * width + x;
            indices.push(index);
            indices.push(index + 1);
            indices.push(index + width);

            indices.push(index + width + 1);
            indices.push(index + width);
            indices.push(index + 1);
        }
    }
    console.log(indices)
    return indices;
}

export default class TerrainRenderer {

    constructor(worldData) {
        this.fetchTerrain(worldData);
    }

    fetchTerrain(worldData)
    {
        let new_req = new XMLHttpRequest();
        new_req.addEventListener("load", () => {
            var parsed = JSON.parse(new_req.responseText)
            console.log(parsed)
            this.startRenderer(parsed.Heights, parsed.Width, parsed.Height, parsed.TextureURL)
        });
        new_req.open("QUERY", "http://localhost:8080", true);
        // new_req.setRequestHeader("Access-Control-Allow-Methods", "*");
        new_req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        new_req.send(JSON.stringify(worldData));
    }

    startRenderer(premadeValues, width, height, textureURL)
    {
        const scene = new THREE.Scene();
        const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);


        const canvas = document.querySelector("#c");
        const renderer = new THREE.WebGLRenderer({canvas, alpha: true});

        const terrain_geometry = new THREE.BufferGeometry()


        var positions, min_y, max_y
        var new_positions = [];
        var uv_coords = []
        for(let x = 0; x < height; x++)
        {
            for(let y = 0; y < width; y++)
            {
                new_positions.push(x / (height));
                new_positions.push(premadeValues[y + width * x]);
                new_positions.push(y / (width));
                uv_coords.push(y / width)
                uv_coords.push(x / height)
            }
        }
        positions = new Float32Array(new_positions);
        uv_coords = new Float32Array(uv_coords)
        min_y = Math.min.apply(this, premadeValues);
        max_y = Math.max.apply(this, premadeValues);
        

        terrain_geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
        terrain_geometry.setAttribute('uv', new THREE.BufferAttribute(uv_coords, 2));
        terrain_geometry.setIndex(getIndices(width, height));
        terrain_geometry.computeVertexNormals();


        // Load Texture from server
        var textureLoader = new THREE.TextureLoader();
        textureLoader.crossOrigin = "anonymous";
        var terrainTexture = textureLoader.load(textureURL);
        terrainTexture.wrapS = THREE.RepeatWrapping;
        terrainTexture.wrapT = THREE.RepeatWrapping;
        terrainTexture.minFilter = THREE.NearestFilter;
        terrainTexture.magFilter = THREE.NearestFilter;
        terrainTexture.flipY = false;

        // Load perlin noise texture

        var perlinTexture = textureLoader.load('public/PerlinNoise.png');
        perlinTexture.wrapS = THREE.RepeatWrapping
        perlinTexture.wrapT = THREE.RepeatWrapping

        var forestTexture = textureLoader.load('public/forest.png');
        forestTexture.wrapS = THREE.MirroredRepeatWrapping
        forestTexture.wrapT = THREE.MirroredRepeatWrapping

        var sandTexture = textureLoader.load('public/sand.png');
        sandTexture.wrapS = THREE.MirroredRepeatWrapping
        sandTexture.wrapT = THREE.MirroredRepeatWrapping

        var rockTexture = textureLoader.load('public/rock.jpg');
        rockTexture.wrapS = THREE.MirroredRepeatWrapping
        rockTexture.wrapT = THREE.MirroredRepeatWrapping

        var snowTexture = textureLoader.load('public/snow.jpg');
        snowTexture.wrapS = THREE.MirroredRepeatWrapping
        snowTexture.wrapT = THREE.MirroredRepeatWrapping
        // Create material
        const terrainShader = new TerrainShader(min_y, max_y, width, height);
        terrainShader.SetTexture("perlinTexture", perlinTexture);
        terrainShader.SetTerrainTresholds([0.4, 0.5, 0.6, 0.7, 0.9])
        terrainShader.SetTerrainTextures([ sandTexture, forestTexture, rockTexture, snowTexture, snowTexture])

        // Create mesh from geometry and add to scene
        const terrain = new THREE.Mesh(terrain_geometry, terrainShader.material);
        terrain.geometry.center()
        let ratio = width / height
        terrain.geometry.scale(50 * ratio, 1, 50 / ratio);
        scene.add(terrain);

        // Move camera to position
        camera.position.z = 30;
        camera.position.y += 30;
        camera.rotateX(degToRad(-45));

        var controls = new FlyControls(camera, canvas)
        controls.movementSpeed = 50
        controls.dragToLook = true
        function animate(timestamp) {
            terrainShader.SetValue("time", timestamp)
            controls.update(0.01)
            // console.log(timestamp)
            // terrain.rotation.y += 0.005;
            renderer.setSize(window.innerWidth, window.innerHeight);
            camera.aspect = window.innerWidth / window.innerHeight;
            camera.updateProjectionMatrix();
            renderer.render( scene, camera );
        }
        renderer.setAnimationLoop( animate );
    }
}