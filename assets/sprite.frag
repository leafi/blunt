#version 330

// per-instance, copied from vertex shader
/*in vec2 frag_size;
flat in int frag_texU;
flat in int frag_texV;
in vec4 frag_tint;*/
//in vec2 size;
//flat in int texU;
//flat in int texV;
//varying vec4 tint;

//layout (location = 0) in vec4 tint2;

out vec4 outColor;

in VOut {
	int fragTexU;
	vec4 fragTint;
} ToFrag;


// gl_FragCoord? probably not actually; scale gets in the way...

void main()
{
    outColor = ToFrag.fragTint + vec4(0.5, 0.5, 1.0, 1.0);
}
