#version 330

// per-instance
in vec4 tint;
in vec2 position;
in vec2 size;
in float scale;
in int texU;
in int texV;
in float angle;


// copied to fragment shader
/*out vec2 frag_size;
flat out int frag_texU;
flat out int frag_texV;
out vec4 frag_tint;*/

out VOut {
	int fragTexU;
	vec4 fragTint;
} ToFrag;

// per-vertex
in vec2 basePosition;

void main()
{
    gl_Position = vec4(scale * size * basePosition + position, 0.0, 1.0);
	ToFrag.fragTint = tint + vec4(1.0, 0.0, 0.0, 0.0);
}
