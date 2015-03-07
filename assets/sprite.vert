#version 330

layout (std140) uniform Props {
	vec2 position;
	vec2 size;

	float scale;
	int texU;
	int texV;
	float angle;

	vec4 tint;
};

//in vec2 position;

void main()
{
    gl_Position = vec4(position, 0.0, 1.0);
}
