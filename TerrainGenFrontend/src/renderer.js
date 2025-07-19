import * as THREE from 'three';
import { degToRad } from 'three/src/math/MathUtils.js';

// const terrain_scale = 0.02;

const _VS = `
varying float v_Pos;
varying vec2 texCoord;
uniform float width;
uniform float height; 
void main() {
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position.x, position.y * 5.0, position.z, 1.0);
    v_Pos = (position.y + 1.0) / 2.0;
    // texCoord = vec2(position.x / float(width) + 0.5, position.z / float(height) + 0.5);
    texCoord = uv;
}
`;

const _FS = `
uniform sampler2D texture1;
uniform float min_y;
uniform float max_y;
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
    if(texture(texture1, texCoord).r < 0.1 || texture(texture1, texCoord).b > 0.1 ) {
    
        gl_FragColor = vec4(relative_height, relative_height, relative_height, 1.0);
    }
    else {
        gl_FragColor = vec4(texture(texture1, texCoord).rgb, 1.0);
    }
}
`;

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

    constructor(width, height) {
        this.fetchTerrain(width, height);
    }

    fetchTerrain(width, height)
    {
        console.log("getting");
        console.log(width, height);
        let new_req = new XMLHttpRequest();
        new_req.addEventListener("load", () => {
            var parsed = JSON.parse(new_req.responseText)
            console.log(parsed)
            this.startRenderer(parsed.Heights, parsed.Width, parsed.Height, parsed.TextureURL)
        });
        new_req.open("QUERY", "http://localhost:8080", true);
        // new_req.setRequestHeader("Access-Control-Allow-Methods", "*");
        new_req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        new_req.send(JSON.stringify({"width": width, "height": height}));
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
        var texture = textureLoader.load(textureURL);
        texture.wrapS = THREE.RepeatWrapping;
        texture.wrapT = THREE.RepeatWrapping;
        texture.minFilter = THREE.NearestFilter;
        texture.magFilter = THREE.NearestFilter;
        texture.flipY = false;

        // Create material from texture
        const material = new THREE.ShaderMaterial({
            uniforms: {max_y : {value: max_y}, min_y : {value: min_y}, width : {value : width}, height : {value : height}, texture1 : {value : texture}},
            vertexShader: _VS,
            fragmentShader: _FS,
        });

        // Create mesh from geometry and add to scene
        const terrain = new THREE.Mesh(terrain_geometry, material);
        terrain.geometry.center()
        let ratio = width / height
        terrain.geometry.scale(50 * ratio, 1, 50 / ratio);
        scene.add(terrain);

        // Move camera to position
        camera.position.z = 30;
        camera.position.y += 30;
        camera.rotateX(degToRad(-45));


        function animate() {
            terrain.rotation.y += 0.005;
            renderer.setSize(window.innerWidth, window.innerHeight);
            camera.aspect = window.innerWidth / window.innerHeight;
            camera.updateProjectionMatrix();
            renderer.render( scene, camera );
        }
        renderer.setAnimationLoop( animate );
    }
}