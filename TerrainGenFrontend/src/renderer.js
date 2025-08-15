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
        this.startRenderer.bind(this);
        this.prepareRenderer();
        this.fetchTerrain(worldData);
        this.startRenderer();
        this.cliffMultiplier = 6.0;
    }

    updateShaderValue(name, newValue) {
        this.terrainShader.SetValue(name, newValue);
    }

    fetchTerrain(worldData)
    {
        let new_req = new XMLHttpRequest();
        new_req.addEventListener("load", () => {
            var parsed = JSON.parse(new_req.responseText)
            console.log(parsed)
            this.loadWorld(parsed.Heights, parsed.Width, parsed.Height, parsed.TextureURL)
        });
        new_req.open("QUERY", "http://localhost:8080", true);
        // new_req.setRequestHeader("Access-Control-Allow-Methods", "*");
        new_req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        new_req.send(JSON.stringify(worldData));
    }

    prepareRenderer() {

        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);


        this.canvas = document.querySelector("#c");
        this.renderer = new THREE.WebGLRenderer({canvas: this.canvas, alpha: true});

        this.terrain_geometry = new THREE.BufferGeometry()

        // Load Texture from server
        this.textureLoader = new THREE.TextureLoader();
        this.textureLoader.crossOrigin = "anonymous";

        // Load perlin noise texture

        var perlinTexture = this.textureLoader.load('public/PerlinNoise.png');
        perlinTexture.wrapS = THREE.RepeatWrapping
        perlinTexture.wrapT = THREE.RepeatWrapping

        var forestTexture = this.textureLoader.load('public/forest.png');
        forestTexture.wrapS = THREE.MirroredRepeatWrapping
        forestTexture.wrapT = THREE.MirroredRepeatWrapping

        var sandTexture = this.textureLoader.load('public/sand.png');
        sandTexture.wrapS = THREE.MirroredRepeatWrapping
        sandTexture.wrapT = THREE.MirroredRepeatWrapping

        var rockTexture = this.textureLoader.load('public/rock.jpg');
        rockTexture.wrapS = THREE.MirroredRepeatWrapping
        rockTexture.wrapT = THREE.MirroredRepeatWrapping

        var snowTexture = this.textureLoader.load('public/snow.jpg');
        snowTexture.wrapS = THREE.MirroredRepeatWrapping
        snowTexture.wrapT = THREE.MirroredRepeatWrapping
        // Create material
        this.terrainShader = new TerrainShader(0.0, 1.0, 200, 200);
        this.terrainShader.SetTexture("perlinTexture", perlinTexture);
        this.terrainShader.SetTerrainTresholds([0.4, 0.5, 0.6, 0.7, 0.9])
        this.terrainShader.SetTerrainTextures([ sandTexture, forestTexture, rockTexture, snowTexture, snowTexture]);
        this.terrainShader.SetValue("cliffMultiplier", 6.0);
    }

    loadWorld(premadeValues, width, height, textureURL) {

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
        this.terrainShader.UpdateValues(min_y, max_y, width, height);

        var terrainTexture = this.textureLoader.load(textureURL);
        terrainTexture.wrapS = THREE.RepeatWrapping;
        terrainTexture.wrapT = THREE.RepeatWrapping;
        terrainTexture.minFilter = THREE.NearestFilter;
        terrainTexture.magFilter = THREE.NearestFilter;
        terrainTexture.flipY = false;
        this.terrain_geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
        this.terrain_geometry.setAttribute('uv', new THREE.BufferAttribute(uv_coords, 2));
        this.terrain_geometry.setIndex(getIndices(width, height));
        this.terrain_geometry.computeVertexNormals();

        this.ratio = width / height
        // Create mesh from geometry and add to scene
        const terrain = new THREE.Mesh(this.terrain_geometry, this.terrainShader.material);
        terrain.geometry.center()
        terrain.geometry.scale(50 * this.ratio, 1, 50 / this.ratio);
        this.scene.add(terrain);
    }

    startRenderer()
    {

        // Move camera to position
        this.camera.position.z = 30;
        this.camera.position.y += 30;
        this.camera.rotateX(degToRad(-45));

        var controls = new FlyControls(this.camera, this.canvas)
        controls.movementSpeed = 50
        controls.dragToLook = true
        var animate = (timestamp) => {
            this.terrainShader.SetValue("time", timestamp)
            controls.update(0.01)
            // console.log(timestamp)
            // terrain.rotation.y += 0.005;
            this.renderer.setSize(window.innerWidth, window.innerHeight);
            this.camera.aspect = window.innerWidth / window.innerHeight;
            this.camera.updateProjectionMatrix();
            this.renderer.render(this.scene, this.camera );
        }
        this.renderer.setAnimationLoop( animate );
    }
}