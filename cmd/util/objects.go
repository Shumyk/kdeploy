package cmd

type Entry struct {
	Key   string
	Value string
}

func EntryOf(key, value string) Entry {
	return Entry{key, value}
}
