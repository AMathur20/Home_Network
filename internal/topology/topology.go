package topology

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadTopology(path string) (*Topology, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Topology{Links: []Link{}}, nil
		}
		return nil, err
	}
	defer file.Close()

	var topo Topology
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&topo); err != nil {
		return nil, err
	}

	return &topo, nil
}

func SaveTopology(path string, topo *Topology) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(topo)
}
