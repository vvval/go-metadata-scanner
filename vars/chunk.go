package vars

type Chunk []string

func (c Chunk) Split(size int) []Chunk {
	var chunks []Chunk

	for i := 0; i < len(c); i += size {
		end := i + size
		if end > len(c) {
			end = len(c)
		}

		chunk := c[i:end]
		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}
