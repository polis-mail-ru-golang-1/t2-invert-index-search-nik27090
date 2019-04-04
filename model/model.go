package model

import (
	"github.com/go-pg/pg"
)

type Model struct {
	db *pg.DB
}

func New(db *pg.DB) Model {
	return Model{
		db: db,
	}
}

type File struct {
	Id   int
	Name string `sql:"name"`
}

func (m Model) ClearModel() {
	_, err := m.db.Exec("TRUNCATE occurences CASCADE")
	if err != nil {
		panic(err)
	}

	_, err = m.db.Exec("TRUNCATE files CASCADE")
	if err != nil {
		panic(err)
	}

	_, err = m.db.Exec("TRUNCATE words CASCADE")
	if err != nil {
		panic(err)
	}
}

type Word struct {
	Id   int
	Word string `sql:"word"`
}

type Occurences struct {
	Id      int
	WordId  int `sql:"word_id"`
	FileId  int `sql:"file_id"`
	Counter int `sql:"count"`
}

func (m Model) AddOccurences(wordid int, fileid int, counter int) {
	Occurences := Occurences{
		WordId:  wordid,
		FileId:  fileid,
		Counter: counter,
	}

	err := m.db.Insert(&Occurences)
	if err != nil {
		panic(err)
	}
}

func (m Model) AddCountersBulk(counters []Occurences) {
	m.db.Insert(&counters)
}

func (m Model) AddCoutersBulk(counters []Occurences) {
	err := m.db.Insert(&counters)
	if err != nil {
		panic(err)
	}
}

func (m Model) GetOrAddWord(wordname string) Word {
	word := Word{
		Word: wordname,
	}

	_, err := m.db.Model(&word).
		Column("id").
		Where("word = ?word").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()

	if err != nil {
		panic(err)
	}

	return word
}

func (m Model) AddWordBulk(words []string) []Word {
	wordsToAdd := make([]Word, len(words))

	for i, val := range words {
		wordsToAdd[i] = Word{
			Word: val,
		}
	}

	err := m.db.Insert(&wordsToAdd)
	if err != nil {
		panic(err)
	}

	return wordsToAdd
}

func (m Model) GetWord(wordName string) *Word {
	word := new(Word)
	m.db.Model(word).Where("word = ?", wordName).Select()
	return word
}

func (m Model) GetWords(wordNames []string) *[]Word {
	result := new([]Word)
	err := m.db.Model(result).WhereIn("word in (?)", pg.In(wordNames)).Select()

	if err != nil {
		panic(err)
	}

	return result
}

func (m Model) GetOrAddFile(filename string) File {
	file := File{
		Name: filename,
	}
	_, err := m.db.Model(&file).
		Column("id").
		Where("name = ?name").
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()

	if err != nil {
		panic(err)
	}

	return file
}

func (m Model) GetFile(fileName string) *File {
	file := new(File)
	m.db.Model(file).Where("name = ?", fileName).Select()
	return file
}

func (m Model) GetFiles(ids []int) []File {
	res := make([]File, 0)
	m.db.Model(&res).Where("id in (?)", pg.In(ids)).Select()

	return res
}

type CounterResult struct {
	File    string
	Counter int
}

func (m Model) getCounters(wordsIds []int) []Occurences {

	var res []Occurences
	err := m.db.Model(&Occurences{}).
		Column("fileid").
		ColumnExpr("SUM(counter) as counter").
		WhereIn("wordid in (?)", pg.In(wordsIds)).
		Group("fileid").
		Select(&res)

	if err != nil {
		panic(err)
	}

	return res
}
