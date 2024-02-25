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
	err := os.RemoveAll("./temp")
	if err != nil {
		fmt.Printf("Fail to remove file: %v", err)
		os.Exit(1)
	}

	err = os.Mkdir("temp", 0755) // 0755 sets permissions for the directory
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

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

	rulesCount, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections:     []string{},
		IgnoreInlineComment:     true,
		IgnoreContinuation:      true,
		SkipUnrecognizableLines: true,
	}, "INI/rules.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	Animations, _ := rulesCount.GetSection("Animations")
	BuildingTypes, _ := rulesCount.GetSection("BuildingTypes")
	AircraftTypes, _ := rulesCount.GetSection("AircraftTypes")
	VehicleTypes, _ := rulesCount.GetSection("VehicleTypes")

	var maxAnimationKey int
	for _, key := range Animations.Keys() {
		i, _ := strconv.Atoi(key.Name())
		if i > maxAnimationKey && i != 999 {
			maxAnimationKey = i
		}
	}

	var maxBuildingKey int
	for _, key := range BuildingTypes.Keys() {
		i, _ := strconv.Atoi(key.Name())
		if i > maxBuildingKey && i != 999 {
			maxBuildingKey = i
		}
	}

	var maxAircraftKey int
	for _, key := range AircraftTypes.Keys() {
		i, _ := strconv.Atoi(key.Name())
		if i > maxAircraftKey && i != 999 {
			maxAircraftKey = i
		}
	}

	var maxVehicleKey int
	for _, key := range VehicleTypes.Keys() {
		i, _ := strconv.Atoi(key.Name())
		if i > maxVehicleKey && i != 999 {
			maxVehicleKey = i
		}
	}

	Animations, _ = rules.GetSection("Animations")
	BuildingTypes, _ = rules.GetSection("BuildingTypes")
	AircraftTypes, _ = rules.GetSection("AircraftTypes")
	VehicleTypes, _ = rules.GetSection("VehicleTypes")

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
		if value.Name() == "EXCITER" {
			fmt.Println(value.Name())
			rules2, arts, sounds := FindItems(value, art, sound, rules)
			fmt.Println("Rules", removeDuplicates(rules2))
			fmt.Println("Arts", removeDuplicates(arts))
			fmt.Println("Sounds", removeDuplicates(sounds))

			rules2 = append(rules2, value.Name())
			arts = append(arts, value.Name())

			list, _ := soundOutput.NewSection("SoundList")
			sourceList, _ := soundOutputCheck.GetSection("SoundList")
			var maxKey int
			for _, key := range sourceList.Keys() {
				i, _ := strconv.Atoi(key.Name())
				if i > maxKey && i != 999 {
					maxKey = i
				}
			}

			for _, item := range removeDuplicates(sounds) {
				if soundOutputCheck.HasSection(item) {
					continue
				}
				println("Adding sound", item)

				newSection, err := sound.GetSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				s, err := soundOutput.NewSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				for _, key := range newSection.KeyStrings() {
					s.NewKey(key, newSection.Key(key).Value())

					// _, err := s.g
					// if err != nil {
					// 	fmt.Printf("no section: %v", err)
					// 	os.Exit(1)
					// }
				}
				maxKey++
				list.NewKey(fmt.Sprintf("%d", maxKey), item)

				soundOutput.SaveToIndent("temp/sound01_copy.ini", "")
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

				newSection, err := art.GetSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				s, err := artOutput.NewSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				for _, key := range newSection.KeyStrings() {
					s.NewKey(key, newSection.Key(key).Value())

					// _, err := s.g
					// if err != nil {
					// 	fmt.Printf("no section: %v", err)
					// 	os.Exit(1)
					// }
				}

				artOutput.SaveToIndent("temp/art_copy.ini", "")
			}

			rulesOutputCheck, err := ini.LoadSources(ini.LoadOptions{
				UnparseableSections:     []string{},
				IgnoreInlineComment:     true,
				IgnoreContinuation:      true,
				SkipUnrecognizableLines: true,
			}, "INI/rules.ini")
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
				os.Exit(1)
			}

			rulesOutput, err := ini.LoadSources(ini.LoadOptions{
				UnparseableSections:     []string{},
				IgnoreInlineComment:     true,
				IgnoreContinuation:      true,
				SkipUnrecognizableLines: true,
			}, []byte(""))
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
				os.Exit(1)
			}

			AnimationsOutput, _ := rulesOutput.NewSection("Animations")
			BuildingTypesOutput, _ := rulesOutput.NewSection("BuildingTypes")
			AircraftTypesOutput, _ := rulesOutput.NewSection("AircraftTypes")
			VehicleTypesOutput, _ := rulesOutput.NewSection("VehicleTypes")

			for _, item := range removeDuplicates(rules2) {
				if rulesOutputCheck.HasSection(item) {
					continue
				}

				var isAnimation bool
				var isBuilding bool
				var isVehicle bool
				var isAircraft bool

				for _, a := range Animations.Keys() {
					v, _ := Animations.GetKey(a.Name())
					if v.Value() == item {
						isAnimation = true
						maxAnimationKey++
						AnimationsOutput.NewKey(fmt.Sprintf("%d", maxAnimationKey), item)
					}
				}

				for _, a := range BuildingTypes.Keys() {
					v, _ := BuildingTypes.GetKey(a.Name())
					if v.Value() == item {
						isBuilding = true
						maxBuildingKey++
						BuildingTypesOutput.NewKey(fmt.Sprintf("%d", maxBuildingKey), item)
					}
				}

				for _, a := range VehicleTypes.Keys() {
					v, _ := VehicleTypes.GetKey(a.Name())
					if v.Value() == item {
						isVehicle = true
						maxVehicleKey++
						VehicleTypesOutput.NewKey(fmt.Sprintf("%d", maxVehicleKey), item)
					}
				}

				for _, a := range AircraftTypes.Keys() {
					v, _ := AircraftTypes.GetKey(a.Name())
					if v.Value() == item {
						isAircraft = true
						maxAircraftKey++
						AircraftTypesOutput.NewKey(fmt.Sprintf("%d", maxAircraftKey), item)
					}
				}

				if !isAnimation && !isBuilding && !isVehicle && !isAircraft {
					fmt.Printf("WARNING: unknown type: %v \n", item)
				}

				newSection, err := rules.GetSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				s, err := rulesOutput.NewSection(item)
				if err != nil {
					fmt.Printf("no section: %v", err)
					os.Exit(1)
				}
				for _, key := range newSection.KeyStrings() {
					s.NewKey(key, newSection.Key(key).Value())

					// _, err := s.g
					// if err != nil {
					// 	fmt.Printf("no section: %v", err)
					// 	os.Exit(1)
					// }
				}

				rulesOutput.SaveToIndent("temp/rules_copy.ini", "")
			}
		}
	}
}

var checkedRules map[string]bool
var checkedArt map[string]bool
var checkedSound map[string]bool

func init() {
	checkedRules = make(map[string]bool)
	checkedArt = make(map[string]bool)
	checkedSound = make(map[string]bool)
}

func FindItems(section *ini.Section, art *ini.File, sound *ini.File, rules *ini.File) ([]string, []string, []string) {
	re := regexp.MustCompile(`^[a-zA-Z-_0-9,]*$`)
	re2 := regexp.MustCompile(`[A-Z]+`)

	var arts []string
	var sounds []string
	var rules2 []string

	for _, value := range section.Keys() {
		if re.MatchString(value.Value()) && re2.MatchString(value.Value()) {
			items := value.Strings(",")
			for _, item := range items {
				item := strings.Trim(item, " ")
				if sound.HasSection(item) && !checkedSound[item] {
					checkedSound[item] = true
					sounds = append(sounds, item)
					rules3, arts3, sounds3 := FindItems(sound.Section(item), art, sound, rules)
					rules2 = append(rules2, rules3...)
					arts = append(arts, arts3...)
					sounds = append(sounds, sounds3...)
				}
				if art.HasSection(item) && !checkedArt[item] {
					checkedArt[item] = true
					arts = append(arts, item)
					rules3, arts3, sounds3 := FindItems(art.Section(item), art, sound, rules)
					rules2 = append(rules2, rules3...)
					arts = append(arts, arts3...)
					sounds = append(sounds, sounds3...)
				}
				if rules.HasSection(item) && !checkedRules[item] {
					checkedRules[item] = true
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
