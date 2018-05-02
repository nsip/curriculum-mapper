package main

import (
	"fmt"
	"github.com/jbrukh/bayesian"
	"github.com/jdkato/prose/tokenize"
	"github.com/recursionpharma/go-csv-map"
	"log"
	"os"
	"strings"
)

// assumes tab-delimited file with header
func read_curriculum(filename string) (records []map[string]string, err error) {
	buf, err := os.Open(filename)
	if err != nil {
		log.Printf("%s: ", filename)
		log.Fatalln(err)
	}
	defer buf.Close()
	reader := csvmap.NewReader(buf)
	reader.Reader.Comma = '\t'
	columns, err := reader.ReadHeader()
	if err != nil {
		log.Printf("%s: ", filename)
		log.Fatalln(err)
	}
	reader.Columns = columns
	records, err = reader.ReadAll()
	if err != nil {
		log.Printf("%s: ", filename)
		log.Fatalln(err)
	}
	return records, nil
}

func main() {
	curriculum, err := read_curriculum("./testdata/curriculum.txt")
	if err != nil {
		log.Fatalln(err)
	}
	// Item    Stage   LearningArea    Strand  Substrand       Text    Elaborations
	classes := make([]bayesian.Class, 0)
	curriculum_map := make(map[string]string)
	for _, record := range curriculum {
		classes = append(classes, bayesian.Class(record["Item"]))
		curriculum_map[record["Item"]] = record["Text"]
	}

	classifier := bayesian.NewClassifierTfIdf(classes...)
	for _, record := range curriculum {
		classifier.Learn(tokenize.TextToWords(record["Text"]+". "+record["Elaborations"]), bayesian.Class(record["Item"]))
	}
	classifier.ConvertTermsFreqToTfIdf()

	for _, x := range classes {
		fmt.Printf("\t%s", x)
	}
	fmt.Println()

	syllabus, err := read_curriculum("./testdata/syllabus.txt")
	// Item    Stage   LearningArea    Strand  Substrand       Outcome Content AC content
	syllabus_map := make(map[string]string)
	for _, record := range syllabus {
		syllabus_map[record["Item"]] = record["Outcome"]
	}

	syllabus_outcome_alignment := make(map[string]string)
	syllabus_content_alignment := make(map[string]string)
	for _, record := range syllabus {
		scores1, max1, _ := classifier.LogScores(tokenize.TextToWords(record["Outcome"]))
		scores2, max2, _ := classifier.LogScores(tokenize.TextToWords(record["Content"]))
		ac := strings.Split(strings.Replace(record["AC content"], "\"", "", -1), "; ")
		ac_match := make(map[bayesian.Class]bool)
		for _, a := range ac {
			ac_match[bayesian.Class(a)] = true
		}
		//log.Printf("%+v\n", ac_match)
		fmt.Printf("%s\t", record["Item"])
		for i := 0; i < len(scores1); i++ {
			match := ""
			_, ok := ac_match[classes[i]]
			//log.Printf(" ---%s: %s--- ", classes[i], ok)
			if ok {
				match = "%"
			}
			is_max1 := ""
			if i == max1 {
				is_max1 = "#"
				syllabus_outcome_alignment[record["Item"]] = string(classes[i])
			}
			is_max2 := ""
			if i == max2 {
				is_max2 = "#"
				syllabus_content_alignment[record["Item"]] = string(classes[i])
			}
			fmt.Printf("%0.3f%s:%0.3f%s%s\t", scores1[i], is_max1, scores2[i], is_max2, match)
		}
		fmt.Println()
	}

	fmt.Printf("\n\n")
	for k, v := range syllabus_outcome_alignment {
		fmt.Printf("The best match for NSW Outcome %s: \"%s\" is AC CD %s: \"%s\" based on outcome text, and AC CD %s: \"%s\" based on content text\n",
			k, syllabus_map[k], v, curriculum_map[v], syllabus_content_alignment[k], curriculum_map[syllabus_content_alignment[k]])
	}
}
