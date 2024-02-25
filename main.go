package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

var checked map[string]bool

func init() {
	checked = make(map[string]bool)
}

type SliceType interface {
	~string | ~int | ~float64 // add more *comparable* types as needed
}

func removeDuplicates[T SliceType](s []T) []T {
	if len(s) < 1 {
		return s
	}

	// sort
	sort.SliceStable(s, func(i, j int) bool {
		return s[i] < s[j]
	})

	prev := 1
	for curr := 1; curr < len(s); curr++ {
		if s[curr-1] != s[curr] {
			s[prev] = s[curr]
			prev++
		}
	}

	return s[:prev]
}

func main() {
	rules, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections:     []string{},
		IgnoreInlineComment:     true,
		IgnoreContinuation:      true,
		SkipUnrecognizableLines: true,
	}, "WARZONE/rules.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	sound, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections:     []string{},
		IgnoreInlineComment:     true,
		IgnoreContinuation:      true,
		SkipUnrecognizableLines: true,
	}, "WARZONE/sound01.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	art, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections:     []string{},
		IgnoreInlineComment:     true,
		IgnoreContinuation:      true,
		SkipUnrecognizableLines: true,
	}, "WARZONE/art.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	soundOutputCheck, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections:     []string{},
		IgnoreInlineComment:     true,
		IgnoreContinuation:      true,
		SkipUnrecognizableLines: true,
	}, "INI/sound01.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	soundOutput, err := ini.Load([]byte(""))
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	// var arr map[string]interface{}
	for _, value := range rules.Sections() {
		if value.Name() == "ORCAAC" {
			fmt.Println(value.Name())
			rules2, arts, sounds := FindItems(value, art, sound, rules)
			fmt.Println("Rules", removeDuplicates(rules2))
			fmt.Println("Arts", removeDuplicates(arts))
			fmt.Println("Sounds", removeDuplicates(sounds))

			list, _ := soundOutput.NewSection("SoundList")
			sourceList, _ := soundOutputCheck.GetSection("SoundList")
			var maxKey int
			for _, key := range sourceList.Keys() {
				i, _ := strconv.Atoi(key.Name())
				if i > maxKey && i != 999 {
					maxKey = i
				}
			}
			sourceList.Keys()
			for _, soundItem := range removeDuplicates(sounds) {
				if soundOutputCheck.HasSection(soundItem) {
					continue
				}
				println("Adding sound", soundItem)

				newSection, err := sound.GetSection(soundItem)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				soundOutput.NewRawSection(soundItem, newSection.Body())

				maxKey++
				list.NewKey(fmt.Sprintf("%d", maxKey), soundItem)

				soundOutput.SaveToIndent("INI/sound01_copy.ini", "")
			}

			artOutputCheck, err := ini.LoadSources(ini.LoadOptions{
				UnparseableSections:     []string{},
				IgnoreInlineComment:     true,
				IgnoreContinuation:      true,
				SkipUnrecognizableLines: true,
			}, "INI/art.ini")
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
				os.Exit(1)
			}

			artOutput, err := ini.LoadSources(ini.LoadOptions{
				UnparseableSections:     []string{},
				IgnoreInlineComment:     true,
				IgnoreContinuation:      true,
				SkipUnrecognizableLines: true,
			}, []byte(""))
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
				os.Exit(1)
			}

			for _, item := range removeDuplicates(arts) {
				if artOutputCheck.HasSection(item) {
					continue
				}
				println("Adding art", item)

				newSection, err := art.GetSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				println(newSection.KeyStrings())
				s, err := artOutput.NewSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				for _, key := range newSection.KeyStrings() {
					println("key", key, strings.Join(s.KeyStrings(), ","))
					fmt.Printf("key: %v\n", s)
					s.NewKey(key, newSection.Key(key).Value())

					// _, err := s.g
					// if err != nil {
					// 	fmt.Printf("no section: %v", err)
					// 	os.Exit(1)
					// }
				}

				artOutput.SaveToIndent("INI/art_copy.ini", "")
			}
		}
	}
}

func FindItems(section *ini.Section, art *ini.File, sound *ini.File, rules *ini.File) ([]string, []string, []string) {
	re := regexp.MustCompile(`^[A-Z-_0-9,]*$`)
	re2 := regexp.MustCompile(`[A-Z]+`)

	var arts []string
	var sounds []string
	var rules2 []string

	if checked[section.Name()] {
		return rules2, arts, sounds

	}

	checked[section.Name()] = true

	for _, value := range section.Keys() {
		if re.MatchString(value.Value()) && re2.MatchString(value.Value()) {
			items := value.Strings(",")
			for _, item := range items {
				item := strings.Trim(item, " ")
				if sound.HasSection(item) {
					sounds = append(sounds, item)
					rules3, arts3, sounds3 := FindItems(sound.Section(item), art, sound, rules)
					rules2 = append(rules2, rules3...)
					arts = append(arts, arts3...)
					sounds = append(sounds, sounds3...)
				}
				if art.HasSection(item) {
					arts = append(arts, item)
					rules3, arts3, sounds3 := FindItems(art.Section(item), art, sound, rules)
					rules2 = append(rules2, rules3...)
					arts = append(arts, arts3...)
					sounds = append(sounds, sounds3...)
				}
				if rules.HasSection(item) {
					rules2 = append(rules2, item)
					rules3, arts3, sounds3 := FindItems(rules.Section(item), art, sound, rules)
					rules2 = append(rules2, rules3...)
					arts = append(arts, arts3...)
					sounds = append(sounds, sounds3...)
				}
			}
		}
	}

	return rules2, arts, sounds
}
