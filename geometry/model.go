package geometry

import (
	"rapidengine/material"
)

// A model can be imported from a 3D object file
// format such as OBJ or STL, and can contain multiple
// meshes / shapes.

type Model struct {
	Meshes    []Mesh
	Materials map[int]material.Material
}

func (m *Model) Render(viewMtx *float32, modelMtx *float32, projMtx *float32, totalTime float64) {
	for _, ms := range m.Meshes {
		ms.Render(m.Materials[ms.ModelMaterial], viewMtx, modelMtx, projMtx, 0, totalTime, 1)
	}
}

func (m *Model) ComputeTangents() {
	for _, ms := range m.Meshes {
		ms.ComputeTangents()
	}
}

func (m *Model) EnableInstancing(num int) {
	for _, ms := range m.Meshes {
		ms.InstancingEnabled = true
		ms.NumInstances = num
	}
}

func NewModel(m Mesh, mat material.Material) Model {
	return Model{
		Meshes:    []Mesh{m},
		Materials: map[int]material.Material{0: mat},
	}
}
