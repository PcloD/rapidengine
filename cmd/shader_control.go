package cmd

import "rapidengine/material"

type ShaderControl struct {
	programs map[string]*material.ShaderProgram
}

func NewShaderControl() ShaderControl {
	return ShaderControl{make(map[string]*material.ShaderProgram)}
}

func (shaderControl *ShaderControl) Initialize() {
	shaderControl.programs = map[string]*material.ShaderProgram{
		"texture":       &material.TextureProgram,
		"colorLighting": &material.ColorLightingProgram,
		"color":         &material.ColorProgram,
		"skybox":        &material.SkyBoxProgram,
	}
	for _, prog := range shaderControl.programs {
		prog.Compile()
	}
}

func (shaderControl *ShaderControl) GetShader(name string) uint32 {
	return shaderControl.programs[name].GetID()
}