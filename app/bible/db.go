package bible

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
	book = strings.ToLower(book)
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
		"genesis":        1,
		"gen":            1,
		"exodus":         2,
		"exod":           2,
		"leviticus":      3,
		"lev":            3,
		"numbers":        4,
		"num":            4,
		"deuteronomy":    5,
		"deut":           5,
		"joshua":         6,
		"josh":           6,
		"judges":         7,
		"judg":           7,
		"ruth":           8,
		"1sam":           9,
		"2sam":           10,
		"1kings":         11,
		"1kgs":           11,
		"2kings":         12,
		"2kgs":           12,
		"1chronicles":    13,
		"1chr":           13,
		"2chronicles":    14,
		"2chr":           14,
		"ezra":           15,
		"nehemiah":       16,
		"neh":            16,
		"esther":         17,
		"esth":           16,
		"job":            18,
		"psalms":         19,
		"ps":             19,
		"proverbs":       20,
		"prov":           20,
		"ecclesiastes":   21,
		"eccl":           21,
		"song":           22,
		"isaiah":         23,
		"isa":            23,
		"jeremiah":       24,
		"jer":            24,
		"lamentations":   25,
		"lam":            25,
		"ezekiel":        26,
		"ezek":           26,
		"daniel":         27,
		"dan":            27,
		"hosea":          28,
		"hos":            28,
		"joel":           29,
		"amos":           30,
		"obadiah":        31,
		"obad":           31,
		"jonah":          32,
		"micah":          33,
		"mic":            33,
		"nahum":          34,
		"nah":            34,
		"habakkuk":       35,
		"hab":            35,
		"zephaniah":      36,
		"zep":            36,
		"haggai":         37,
		"hag":            37,
		"zechariah":      38,
		"zech":           38,
		"zec":            38,
		"malachi":        39,
		"mal":            39,
		"matthew":        40,
		"matt":           40,
		"mark":           41,
		"luke":           42,
		"john":           43,
		"acts":           44,
		"romans":         45,
		"rom":            45,
		"1corinthians":   46,
		"1cor":           46,
		"2corinthians":   47,
		"2cor":           47,
		"galatians":      48,
		"gal":            48,
		"ephesians":      49,
		"eph":            49,
		"philippians":    50,
		"phil":           50,
		"colossians":     51,
		"col":            51,
		"1thessalonians": 52,
		"1thess":         52,
		"2thessalonians": 53,
		"2thess":         53,
		"1timothy":       54,
		"1tim":           54,
		"2timothy":       55,
		"2tim":           55,
		"titus":          56,
		"philemon":       57,
		"phlm":           57,
		"hebrews":        58,
		"heb":            58,
		"james":          59,
		"jas":            59,
		"1peter":         60,
		"1pet":           60,
		"2peter":         61,
		"2pet":           61,
		"1john":          62,
		"2john":          63,
		"3john":          64,
		"jude":           65,
		"revelation":     66,
	}
)
