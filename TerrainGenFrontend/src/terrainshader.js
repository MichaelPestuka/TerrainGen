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
uniform sampler2D seafoamTexture;
uniform sampler2D distortionTexture;

uniform sampler2D[5] landTextures;
uniform float[5] landHeights;

uniform float waveHeight;

varying float v_Pos;
varying vec3 v_Norm;
varying vec2 texCoord;

uniform float time;
uniform float cliffMultiplier;
uniform float cliffOffset;

uniform float wavesOffset;
uniform float wavesBlending;

float inverseLerp(float a, float b, float value) {
    return clamp((value - a)/(b-a), 0.0, 1.0);
}

float wave(vec2 dir, float amplitude, float moveSpeed, float wavelength, float phase) {
    wavelength *= 0.01;
    moveSpeed *= 0.01;
    float k = 2.0 * 3.141 / wavelength;
    float w = moveSpeed * k;
    float dotProduct = dot(dir, texCoord);
    float periodicAmplitude  = amplitude * 0.5 * (1.0 + sin(6.282 * time * 0.001 / wavelength + phase));
    return periodicAmplitude * sin(k * dotProduct - w * time);
}

float allWaves() {
    float height =      wave(vec2(1.0, 1.0), 0.3, 0.5, 5.0, 0.0);
    height +=      wave(vec2(0.5, 1.0), 0.15, 0.632356, 2.5, 0.5);
    height +=      wave(vec2(1.0, 0.8), 0.10, 0.98976868, 2.5, 1.0);
    height +=      wave(vec2(0.4, 0.8), 0.20, 1.231251, 2.5, 3.0);

    return clamp(height / 4.0, 0.0, 1.0);
}

void main() {

    float border = ((sin(time) + 1.0) / waveHeight) + 0.5 - (sin(1.57) + 1.0) / waveHeight; //0.46875;
    border -= floor(border);
    if(v_Pos <= border) // Below waves
    {
        float wave_x = allWaves();
        vec4 seaColor = vec4(0.05, 0.3, 0.8, 1.0);
        seaColor.b += texture(perlinTexture, texCoord * 200.0).r * wave_x;

        float wavesMinimum = ((sin(time - 0.8) + 1.0) / waveHeight) + 0.5 - (sin(1.57) + 1.0) / waveHeight - 0.0002;
        float wavyness = inverseLerp(-wavesBlending, wavesBlending, v_Pos - wavesMinimum + wavesOffset);
        vec2 foamOffset = texture(distortionTexture, texCoord + time / 100.0).xy;
        vec4 foamColor = vec4(texture(perlinTexture, texCoord * 100.0 + foamOffset).x, texture(perlinTexture, texCoord * 100.0 + foamOffset).y + 0.5, 1.0, 1.0);
        gl_FragColor = seaColor * (1.0 - wavyness) + foamColor * wavyness;
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

        // if(v_Pos > 0.5 && v_Pos < 0.51) {
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