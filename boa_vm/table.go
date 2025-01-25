package main


type Table struct {
  entries  map[string]Value
  count    int
  capacity int
}

func initMap() *Table {
  return &Table{
    entries:  make(map[string]Value),
    count  :  0,
    capacity: 0,
  }
}

func (t *Table) tableSet(k string, v Value){
  t.entries[k] = v
}

func (t *Table) tableGet(k string) *Value {
  v, ok := t.entries[k]; if !ok { return nil } else { return &v }
}

func (t *Table) tableDelete(k string) {
  delete(t.entries, k)
}

