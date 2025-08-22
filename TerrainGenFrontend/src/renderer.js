import * as THREE from 'three';
import { degToRad, randFloat } from 'three/src/math/MathUtils.js';
import TerrainShader from './terrainshader';
import { FlyControls } from 'three/addons/controls/FlyControls.js';

// import Stats from 'three/examples/jsm/libs/stats.module.js'; for FPS counter

// Returns vertex indices for a plane of size width * height
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
    return indices;
}

// Class holding map rendering and loading logic
export default class TerrainRenderer {

    constructor(worldData) {
        this.startRenderer.bind(this);
        this.prepareRenderer();
        this.fetchTerrain(worldData);
        this.startRenderer();
    }

    // Helper function for changing shader unforms
    updateShaderValue(name, newValue) {
        this.terrainShader.SetValue(name, newValue);
    }

    // Loads map data from server
    fetchTerrain(worldData)
    {
        let new_req = new XMLHttpRequest();
        new_req.addEventListener("load", () => {
            var parsed = JSON.parse(new_req.responseText)
            this.loadWorld(parsed.Heights, parsed.Width, parsed.Height, parsed.TextureURL)
        });
        new_req.open("QUERY", "http://localhost:8080", true);
        new_req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        new_req.send(JSON.stringify(worldData));
    }

    // Prepares renderer and loads textures
    prepareRenderer() {

        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);


        this.canvas = document.querySelector("#c");
        this.renderer = new THREE.WebGLRenderer({canvas: this.canvas, alpha: true});


        // Load terrain Texture from server (Currently unused)
        this.textureLoader = new THREE.TextureLoader();
        this.textureLoader.crossOrigin = "anonymous";

        // Load noise and terrain textures 

        var perlinTexture = this.textureLoader.load('worley.png');
        perlinTexture.wrapS = THREE.MirroredRepeatWrapping
        perlinTexture.wrapT = THREE.MirroredRepeatWrapping

        var distortionTexture = this.textureLoader.load('WaterDistortion.png');
        distortionTexture.wrapS = THREE.MirroredRepeatWrapping
        distortionTexture.wrapT = THREE.MirroredRepeatWrapping

        var forestTexture = this.textureLoader.load('forest.png');
        forestTexture.wrapS = THREE.MirroredRepeatWrapping
        forestTexture.wrapT = THREE.MirroredRepeatWrapping

        var sandTexture = this.textureLoader.load('sand.png');
        sandTexture.wrapS = THREE.MirroredRepeatWrapping
        sandTexture.wrapT = THREE.MirroredRepeatWrapping

        var rockTexture = this.textureLoader.load('rock.jpg');
        rockTexture.wrapS = THREE.MirroredRepeatWrapping
        rockTexture.wrapT = THREE.MirroredRepeatWrapping

        var snowTexture = this.textureLoader.load('snow.jpg');
        snowTexture.wrapS = THREE.MirroredRepeatWrapping
        snowTexture.wrapT = THREE.MirroredRepeatWrapping

        var seafoamTexture = this.textureLoader.load('seafoam.jpg');
        seafoamTexture.wrapS = THREE.MirroredRepeatWrapping
        seafoamTexture.wrapT = THREE.MirroredRepeatWrapping

        // Create shader instance
        this.terrainShader = new TerrainShader(0.0, 1.0, 200, 200);
        this.terrainShader.SetTexture("perlinTexture", perlinTexture);
        this.terrainShader.SetTexture("distortionTexture", distortionTexture);
        this.terrainShader.SetTexture("seafoamTexture", seafoamTexture);
        this.terrainShader.SetTerrainTresholds([0.0, 0.5, 0.52, 0.7, 0.8])
        this.terrainShader.SetTerrainTextures([sandTexture, sandTexture, forestTexture, rockTexture, snowTexture]);
        this.terrainShader.SetValue("cliffMultiplier", 6.0);
    }

    loadWorld(premadeValues, width, height, textureURL) {
        this.scene.clear()

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
        this.positions = new Float32Array(new_positions);
        uv_coords = new Float32Array(uv_coords)
        this.terrainShader.UpdateValues(width, height);

        var terrainTexture = this.textureLoader.load(textureURL);
        terrainTexture.wrapS = THREE.RepeatWrapping;
        terrainTexture.wrapT = THREE.RepeatWrapping;
        terrainTexture.minFilter = THREE.NearestFilter;
        terrainTexture.magFilter = THREE.NearestFilter;
        terrainTexture.flipY = false;
        this.terrain_geometry = new THREE.BufferGeometry()
        this.terrain_geometry.setAttribute('position', new THREE.BufferAttribute(this.positions, 3));
        this.terrain_geometry.setAttribute('uv', new THREE.BufferAttribute(uv_coords, 2));
        this.terrain_geometry.setIndex(getIndices(width, height));
        this.terrain_geometry.computeVertexNormals();

        this.ratio = width / height
        // Create mesh from geometry and add to scene
        this.terrain = new THREE.Mesh(this.terrain_geometry, this.terrainShader.material);
        this.terrain.geometry.center()
        this.terrain.geometry.scale(50 * this.ratio, 1, 50 / this.ratio);
        this.scene.add(this.terrain);
        this.scatterTrees();
    }

    // Creates an instanced mesh of trees scattered on the map
    scatterTrees() {
        // Load tree mesh and texture
        const treeMesh = new THREE.SphereGeometry(0.15, 8, 6) // Trees are just low poly spheres :)
        var treeTexture = this.textureLoader.load('forest.png');
        treeTexture.colorSpace = THREE.SRGBColorSpace // fix for MeshBasicMaterial color rendering
        const mat = new THREE.MeshBasicMaterial({map: treeTexture})

        // Create instanced mesh of trees
        var instanced = new THREE.InstancedMesh(treeMesh, mat, this.positions.length)
        instanced.instanceMatrix.setUsage(THREE.DynamicDrawUsage)

        // Get vertex normals from map, used for determining slopes
        const normalsArray = this.terrain_geometry.getAttribute('normal').array

        var renderedTrees = 0;
        for(let i = 0; i < this.positions.length; i += 3) {
            // Don't render trees on sand and under water
            if(this.positions[i + 1] < 0.05) {           
                continue;
            }

            // Dont render trees on steep slopes
            var angle = Math.acos(new THREE.Vector3(normalsArray[i], normalsArray[i + 1], normalsArray[i + 2]).normalize().dot(new THREE.Vector3(0.0, 1.0, 0.0)))
            if(angle > 0.1 ) {
                continue
            }

            // Set tree position
            var matrix = new THREE.Matrix4()
            instanced.getMatrixAt(renderedTrees, matrix)
            matrix.elements[12] = this.positions[i] + randFloat(-0.1, 0.1)
            matrix.elements[13] = this.positions[i + 1] * 5.0
            matrix.elements[14] = this.positions[i + 2] + randFloat(-0.1, 0.1)
    
            instanced.setMatrixAt(renderedTrees, matrix)
            renderedTrees += 1;
        }
        instanced.count = renderedTrees + 1
        instanced.instanceMatrix.needsUpdate = true
        this.scene.add(instanced)
    }

    // Starts rendering map
    startRenderer()
    {
        // FPS counter for debugging
        // const stats = new Stats()
        // document.body.appendChild(stats.dom)

        // Move camera to position
        this.camera.position.z = 30;
        this.camera.position.y += 30;
        this.camera.rotateX(degToRad(-45));

        // Basic controls
        var controls = new FlyControls(this.camera, this.canvas)
        controls.movementSpeed = 50
        controls.dragToLook = true

        // Animation loop
        var animate = (timestamp) => {
            this.terrainShader.SetValue("time", timestamp / 1000.0)
            controls.update(0.01)
            this.renderer.setSize(window.innerWidth, window.innerHeight);
            this.camera.aspect = window.innerWidth / window.innerHeight;
            this.camera.updateProjectionMatrix();
            this.renderer.render(this.scene, this.camera );

            // stats.update()
        }
        this.renderer.setAnimationLoop( animate );
    }
}