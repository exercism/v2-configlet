package cmd

import (
	"fmt"
	"log"
	"path/filepath"
)

func ExampleFormat() {
	diff, formatted, err := formatFile(filepath.FromSlash("../fixtures/format/unformatted/config.json"), formatTopics)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(diff)
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
	// +        "control-flow-(conditionals)",
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
	// +        "text-formatting",
	// +        "time"
	// +      ],
	// +      "unlocked_by": null,
	// +      "uuid": "002"
	// +    }
	// +  ],
	// @@ -10,27 +36,3 @@
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
	// -
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
	//         "control-flow-(conditionals)",
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
	//         "text-formatting",
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
