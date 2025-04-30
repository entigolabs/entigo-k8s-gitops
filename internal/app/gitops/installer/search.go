package installer

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
)

func (i *Installer) search(node *yaml.Node, keys []string) (string, error) {
	returnValue := ""
	var returnError error
	identifier := keys[0]
	if node.Kind == yaml.DocumentNode {
		returnValue, returnError = i.search(node.Content[0], keys)
	}
	if seqPos, err := strconv.Atoi(identifier); err == nil {
		if len(node.Content)-1 < seqPos {
			return "", fmt.Errorf("skiping '%s' copy in %s - key doesn't exist", editInfo.workingKey, editInfo.workingFile) // todo check this error
		}
		seqPosNode := node.Content[seqPos]
		if seqPosNode.Kind == yaml.ScalarNode {
			return seqPosNode.Value, nil
		} else {
			returnValue, returnError = i.search(node.Content[seqPos], keys[1:])
		}
	}
	for j, n := range node.Content {
		if n.Value == identifier {
			if len(keys) <= 1 && node.Content[j+1].Content == nil {
				editInfo.keyExist = true
				return node.Content[j+1].Value, nil
			} else {
				returnValue, returnError = i.search(node.Content[j+1], keys[1:])
			}
		} else {
			if len(keys) > 1 && j == len(node.Content)-1 && !editInfo.keyExist {
				return "", fmt.Errorf("%s don't have following key '%s'", editInfo.workingFile, editInfo.workingKey)
			}
		}
	}
	return returnValue, returnError
}
