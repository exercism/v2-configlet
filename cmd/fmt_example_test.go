package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ExampleFormat() {
	tmp, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	diff, err := formatFile(filepath.FromSlash("../fixtures/format/unformatted/config.json"), tmp.Name(), formatTopics, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(diff)

	formatted, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(formatted))

	// Output:
	// @@ -2 +1,0 @@
	// -  "language": "Numbers",
	// @@ -5 +4,28 @@
	// -  "ignore_pattern": "example(?!.*test)",
	// +  "exercises": [
	// +    {
	// +      "core": false,
	// +      "difficulty": 1,
	// +      "slug": "one",
	// +      "topics": [
	// +        "booleans",
	// +        "control_flow_conditionals",
	// +        "integers",
	// +        "logic"
	// +      ],
	// +      "unlocked_by": null,
	// +      "uuid": "001"
	// +    },
	// +    {
	// +      "core": false,
	// +      "difficulty": 1,
	// +      "slug": "two",
	// +      "topics": [
	// +        "equality",
	// +        "mathematics",
	// +        "text_formatting",
	// +        "time"
	// +      ],
	// +      "unlocked_by": null,
	// +      "uuid": "002"
	// +    }
	// +  ],
	// @@ -10,26 +36,3 @@
	// -  "exercises": [
	// -{
	// -      "uuid": "001",
	// -      "slug": "one",
	// -      "core": false,
	// -      "unlocked_by": null,
	// -      "difficulty": 1,
	// -      "topics": [
	// -            "Control-flow (conditionals)",
	// -            "Logic",
	// -            "Booleans",
	// -            "Integers"
	// -      ]
	// -    },
	// -    {
	// -      "uuid": "002",
	// -      "slug": "two",
	// -      "core": false,
	// -      "unlocked_by": null,
	// -      "difficulty": 1,
	// -      "topics": [
	// -        "Time",
	// -        "Mathematics",
	// -        "Text formatting",
	// -        "Equality"
	// -      ]}]}
	// +  "ignore_pattern": "example(?!.*test)",
	// +  "language": "Numbers"
	// +}
	//
	// {
	//   "active": true,
	//   "blurb": "",
	//   "exercises": [
	//     {
	//       "core": false,
	//       "difficulty": 1,
	//       "slug": "one",
	//       "topics": [
	//         "booleans",
	//         "control_flow_conditionals",
	//         "integers",
	//         "logic"
	//       ],
	//       "unlocked_by": null,
	//       "uuid": "001"
	//     },
	//     {
	//       "core": false,
	//       "difficulty": 1,
	//       "slug": "two",
	//       "topics": [
	//         "equality",
	//         "mathematics",
	//         "text_formatting",
	//         "time"
	//       ],
	//       "unlocked_by": null,
	//       "uuid": "002"
	//     }
	//   ],
	//   "foregone": [
	//     "three",
	//     "four"
	//   ],
	//   "ignore_pattern": "example(?!.*test)",
	//   "language": "Numbers"
	// }
}
