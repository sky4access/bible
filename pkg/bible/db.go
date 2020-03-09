package db

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/doug-martin/goqu.v4"
)

type Entry struct {
	Verse   int32  `db:"verse"`
	Content string `db:"content"`
}

type BibleDB struct {
	*goqu.Database
}

func (m BibleDB) QueryBible(book string, chapter int, verses ...int) (string, []Entry, error) {
	ds := m.From("bible")
	var q *goqu.Dataset
	bookName := BooksByNumber[BooksByName[book]]
	var entryName string
	if len(verses) == 1 {
		q = ds.Select(
			goqu.I("verse").As("verse"),
			goqu.I("content").As("content")).
			Where(goqu.I("book").Eq(BooksByName[book]),
				goqu.I("chapter").Eq(chapter),
				goqu.I("verse").Eq(verses[0]))
		entryName = fmt.Sprintf("%v %v:%v", bookName, chapter, verses[0])

	} else if len(verses) == 2 {
		q = ds.Select(
			goqu.I("verse").As("verse"),
			goqu.I("content").As("content")).
			Where(goqu.I("book").Eq(BooksByName[book]),
				goqu.I("chapter").Eq(chapter),
				goqu.I("verse").Gte(verses[0]),
				goqu.I("verse").Lte(verses[1]))
		entryName = fmt.Sprintf("%v %v:%v-%v\n", bookName, chapter, verses[0], verses[1])
	} else {
		q = nil
	}

	if q != nil {
		var entries []Entry
		if err := q.ScanStructs(&entries); err != nil {
			return entryName, nil, err
		}
		return strings.TrimSpace(entryName), entries, nil
	} else {
		return entryName, nil, errors.New("verses must be 1 or 2 arguments")
	}

}

func (m *BibleDB) ParseVerses(s string) string {
	arr := strings.Split(s, "+")
	b, cv := arr[0], arr[1]
	arr = strings.Split(cv, ":")
	c, v := arr[0], arr[1]
	vs := strings.Split(v, "-")

	chapter, _ := strconv.Atoi(c)
	if len(vs) == 1 {
		v1, _ := strconv.Atoi(vs[0])
		entryName, entries, err := m.QueryBible(b, chapter, v1)
		check(err)
		return printVerse(entryName, entries)
	} else if len(vs) == 2 {
		v1, _ := strconv.Atoi(vs[0])
		v2, _ := strconv.Atoi(vs[1])
		entryName, entries, err := m.QueryBible(b, chapter, v1, v2)
		check(err)
		return printVerse(entryName, entries)
	}

	return ""
}

func printVerse(entryName string, entries []Entry) string {

	buffer := bytes.NewBufferString("")
	buffer.WriteString(entryName + " ")
	for _, entry := range entries {
		buffer.WriteString(fmt.Sprintf("(%v)%v", strconv.Itoa(int(entry.Verse)), entry.Content))
	}
	return buffer.String()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	BooksByNumber = map[int]string{
		1:  "창세기",
		2:  "출애굽기",
		3:  "레위기",
		4:  "민수기",
		5:  "신명기",
		6:  "여호수아",
		7:  "사사기",
		8:  "룻기",
		9:  "사무엘상",
		10: "사무엘하",
		11: "열왕기상",
		12: "열왕기하",
		13: "역대상",
		14: "역대하",
		15: "에스라",
		16: "느헤미야",
		17: "에스더",
		18: "욥기",
		19: "시편",
		20: "잠언",
		21: "전도서",
		22: "아가",
		23: "이사야",
		24: "예레미야",
		25: "예레미야애가",
		26: "에스겔",
		27: "다니엘",
		28: "호세아",
		29: "요엘",
		30: "아모스",
		31: "오바댜",
		32: "요나",
		33: "미가",
		34: "나훔",
		35: "하박국",
		36: "스바냐",
		37: "학개",
		38: "스가랴",
		39: "말라기",
		40: "마태복음",
		41: "마가복음",
		42: "누가복음",
		43: "요한복음",
		44: "사도행전",
		45: "로마서",
		46: "고린도전서",
		47: "고린도후서",
		48: "갈라디아서",
		49: "에베소서",
		50: "빌립보서",
		51: "골로새서",
		52: "데살로니가전서",
		53: "데살로니가후서",
		54: "디모데전서",
		55: "디모데후서",
		56: "디도서",
		57: "빌레몬서",
		58: "히브리서",
		59: "야고보서",
		60: "베드로전서",
		61: "베드로후서",
		62: "요한1서",
		63: "요한2서",
		64: "요한3서",
		65: "유다서",
		66: "요한계시록",
	}

	BooksByName = map[string]int{
		"genesis":           1,
		"exodus":            2,
		"leviticus":         3,
		"numbers":           4,
		"deuteronomy":       5,
		"joshua":            6,
		"judges":            7,
		"ruth":              8,
		"1st samuel":        9,
		"2nd samuel":        10,
		"1st kings":         11,
		"2nd kings":         12,
		"1st chronicles":    13,
		"2nd chronicles":    14,
		"ezra":              15,
		"nehemiah":          16,
		"esther":            17,
		"job":               18,
		"psalms":            19,
		"proverbs":          20,
		"ecclesiastes":      21,
		"song of solomon":   22,
		"isaiah":            23,
		"jeremiah":          24,
		"lamentations":      25,
		"ezekiel":           26,
		"daniel":            27,
		"hosea":             28,
		"joel":              29,
		"amos":              30,
		"obadiah":           31,
		"jonah":             32,
		"micah":             33,
		"nahum":             34,
		"habakkuk":          35,
		"zephaniah":         36,
		"haggai":            37,
		"zechariah":         38,
		"malachi":           39,
		"matthew":           40,
		"mark":              41,
		"luke":              42,
		"john":              43,
		"acts":              44,
		"romans":            45,
		"1st corinthians":   46,
		"2nd corinthians":   47,
		"galatians":         48,
		"ephesians":         49,
		"philippians":       50,
		"colossians":        51,
		"1st thessalonians": 52,
		"2nd thessalonians": 53,
		"1st timothy":       54,
		"2nd timothy":       55,
		"titus":             56,
		"philemon":          57,
		"hebrews":           58,
		"james":             59,
		"1st peter":         60,
		"2nd peter":         61,
		"1st john":          62,
		"2nd john":          63,
		"3rd john":          64,
		"jude":              65,
		"revelation":        66,
	}
)
