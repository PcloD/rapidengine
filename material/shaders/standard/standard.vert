#version 410

out vec3 FragPos;
out vec3 TexCoords;
out mat3 TBN;
out vec3 Normal;
out vec3 ReflectedVector;
out vec3 RefractedVector;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 tex;
layout (location = 2) in vec3 normal;
layout (location = 3) in vec3 tangent;
layout (location = 4) in vec3 bitTangent;

uniform float scale;
uniform float displacement;

uniform sampler2D heightMap;

uniform mat4 modelMtx;
uniform mat4 viewMtx;
uniform mat4 projectionMtx;

uniform vec3 viewPos;

uniform float refractivity;

float getDisplacement();

void main() {
    vec3 finalPosition = position + (normal * getDisplacement());

    //	Vertex position 
    gl_Position = projectionMtx * viewMtx * modelMtx * vec4(finalPosition, 1.0);

    // Normal vector
    Normal = mat3(transpose(inverse(modelMtx))) * normal;

    // Fragment position
    FragPos =  vec3(modelMtx * vec4(finalPosition, 1.0));

    // Texture coordinates
    TexCoords = tex / scale;

    // Normal mapping
    vec3 T = normalize(vec3(modelMtx * vec4(tangent,   0.0)));
    vec3 B = normalize(vec3(modelMtx * vec4(bitTangent, 0.0)));
    vec3 N = normalize(vec3(modelMtx * vec4(normal,    0.0)));
    TBN = mat3(T, B, N);

    // Reflection/refraction
    vec3 viewVector = normalize(finalPosition - viewPos);
    ReflectedVector = reflect(viewVector, normal);
    RefractedVector = refract(viewVector, normal, refractivity);
}

float getDisplacement() {
    return texture(heightMap, tex.xy / scale).x * displacement;
}
