import * as THREE from 'three';

const _VS = `
varying float v_Pos;
varying vec2 texCoord;
uniform float width;
uniform float height; 

varying vec3 v_Norm;
void main() {
    v_Norm = normal;
    if(position.y < 0.0)
    {
        gl_Position = projectionMatrix * modelViewMatrix * vec4(position.x, 0.0 * 5.0, position.z, 1.0);
    }
    else {
        gl_Position = projectionMatrix * modelViewMatrix * vec4(position.x, position.y * 5.0, position.z, 1.0);
    }
    v_Pos = (position.y + 1.0) / 2.0;
    texCoord = uv;
}
`;

const _FS = `
uniform sampler2D perlinTexture;

uniform sampler2D[5] landTextures;
uniform float[5] landHeights;

uniform float waveHeight;

varying float v_Pos;
varying vec3 v_Norm;
varying vec2 texCoord;

uniform float time;
uniform float cliffMultiplier;
uniform float cliffOffset;

float inverseLerp(float a, float b, float value) {
    return clamp((value - a)/(b-a), 0.0, 1.0);
}


void main() {

    float border = ((sin(time / 1000.0) + 1.0) / waveHeight) + 0.5; //0.46875;
    border -= floor(border);
    if(v_Pos <= border - 0.01) // Below waves
    {

        gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
    }
    else if(v_Pos <= border && v_Pos > border - 0.01) { // Foam layer
        gl_FragColor = vec4(0.0, 0.5, 1.0, v_Pos);
    }
    else { // land

        for(int i = 0; i < 5; i++) {
            float drawStrength = inverseLerp(-0.005, 0.005, v_Pos - landHeights[i]);
            switch (i) { // A loathsome workaround, smapler2D arrays have to be indexed with constant expressions
                case 0:
                    gl_FragColor = gl_FragColor * (1.0-drawStrength) + texture(landTextures[0], texCoord * 20.0).rgba * drawStrength;
                    break;
                case 1:
                    gl_FragColor = gl_FragColor * (1.0-drawStrength) + texture(landTextures[1], texCoord * 20.0).rgba * drawStrength;
                    break;
                case 2:
                    gl_FragColor = gl_FragColor * (1.0-drawStrength) + texture(landTextures[2], texCoord * 20.0).rgba * drawStrength;
                    break;
                case 3:
                    gl_FragColor = gl_FragColor * (1.0-drawStrength) + texture(landTextures[3], texCoord * 20.0).rgba * drawStrength;
                    break;
                case 4:
                    gl_FragColor = gl_FragColor * (1.0-drawStrength) + texture(landTextures[4], texCoord * 20.0).rgba * drawStrength;
                    break;
            }
        }
        float angle = acos(dot(normalize(v_Norm), vec3(0.0, 1.0, 0.0)));
        if(angle > 0.0) {
            // gl_FragColor = texture(landTextures[2], texCoord * 20.0).rgba;
            angle += cliffOffset;
            angle = clamp(angle * cliffMultiplier, 0.0, 1.0);
            gl_FragColor = gl_FragColor * (1.0-angle) + texture(landTextures[3], texCoord * 20.0).rgba * angle;
        }

        if(v_Pos < landHeights[0]) {
            gl_FragColor = vec4(texture(landTextures[0], texCoord * 20.0).rgb, 1.0);
        }

        // if(v_Pos > 0.49 && v_Pos < 0.51) {
        //     gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);        
        // }
    }
}
`;

export default class TerrainShader {
    constructor(width, height) {
        this.material = new THREE.ShaderMaterial({
            uniforms: {width : {value : width}, height : {value : height}},
            vertexShader: _VS,
            fragmentShader: _FS,
        });
    }

    UpdateValues(width, height) {
        this.material.uniforms["width"] = {value : width}
        this.material.uniforms["height"] = {value : height}
    }

    SetTexture(name, texture) {
        this.material.uniforms[name] = {value : texture}
    }

    SetValue(name, value) {
        this.material.uniforms[name] = {value : value}
    }

    SetTerrainTresholds(tresholds) {
        this.material.uniforms["landHeights"] = {value : tresholds}
    }
    SetTerrainTextures(textures) {
        this.material.uniforms["landTextures"] = {value : textures}
    }
}