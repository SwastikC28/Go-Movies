package web

// Get common associations between passed associations and options.
func GetMatchedAssociations(includes, associations []string) []string {
	var matchedAssociations = []string{}

	if len(includes) == 0 || len(associations) == 0 {
		return []string{}
	}

	for _, association := range associations {
		for _, include := range includes {
			if include == association {
				matchedAssociations = append(matchedAssociations, include)
			}
		}
	}

	return matchedAssociations
}
