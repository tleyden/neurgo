package neurgo

// TODO: any way to make this generic using interface{}?
func pruneEmptyElements(things []*weightedInput) []*weightedInput {

	result := make([]*weightedInput,0)

	for _, element := range things {

		if element != nil {
			result = append(result, element)
		}
	}
	return result

}

