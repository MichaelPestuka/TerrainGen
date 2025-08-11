import * as THREE from 'three';

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
uniform sampler2D terrainTexture;
uniform sampler2D perlinTexture;
uniform float min_y;
uniform float max_y;
varying float v_Pos;
varying vec2 texCoord;

uniform float time;

void main() {
    // float relative_height = (v_Pos.y - min_y) / (max_y);
    float relative_height = v_Pos;
    gl_FragColor = vec4(relative_height , relative_height, relative_height, 1.0);
    if(relative_height < 0.5 && relative_height >= 0.0)
    {
        gl_FragColor = vec4(0.2, 0.2, 0.5, 1.0);
    }
    if(texture(terrainTexture, texCoord).b < 0.1 ) {
    
        gl_FragColor = vec4(relative_height, relative_height, relative_height, 1.0);
    }
    else {
        gl_FragColor = vec4(texture(terrainTexture, texCoord).rgb, 1.0);
        if(gl_FragColor.b >= 0.9 && gl_FragColor.g >= 0.9) {
            float foamHeight = (sin(time / 500.0) + 1.0)/5.0;
            gl_FragColor = texture(perlinTexture, vec2(texCoord.x * 5.0 + sin(time / 10000.0), texCoord.y * 5.0 + sin(time / 10000.0))).r < foamHeight ? vec4(0, 1.0, 1.0, 1.0) : vec4(0.0, 0.0, 1.0, 1.0);
        }
        // gl_FragColor = vec4(relative_height, relative_height, relative_height, 1.0);
    }
    float border = ((sin(time / 1000.0) + 1.0) / 64.0) + 0.45; //0.46875;
    border -= floor(border);
    
    if(relative_height < border && relative_height > border - 0.02) {
        gl_FragColor = vec4(0.0, 0.5, 1.0, relative_height);
    }
    // gl_FragColor = vec4(relative_height, relative_height, relative_height, 1.0);
}
`;

export default class TerrainShader {
    constructor(max_y, min_y, width, height) {
        this.material = new THREE.ShaderMaterial({
            uniforms: {max_y : {value: max_y}, min_y : {value: min_y}, width : {value : width}, height : {value : height}},
            vertexShader: _VS,
            fragmentShader: _FS,
        });
    }

    SetTexture(name, texture) {
        this.material.uniforms[name] = {value : texture}
    }

    SetValue(name, value) {
        this.material.uniforms[name] = {value : value}
    }

}